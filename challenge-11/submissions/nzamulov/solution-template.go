// Package challenge11 contains the solution for Challenge 11.
package challenge11

import (
	"context"
	"net/http"
	"time"
	"sync"
	"io"
	"errors"
	"bytes"
	"strings"
	"math"
	"fmt"
	
	"golang.org/x/net/html"
)

// ContentFetcher defines an interface for fetching content from URLs
type ContentFetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// ContentProcessor defines an interface for processing raw content
type ContentProcessor interface {
	Process(ctx context.Context, content []byte) (ProcessedData, error)
}

// ProcessedData represents structured data extracted from raw content
type ProcessedData struct {
	Title       string
	Description string
	Keywords    []string
	Timestamp   time.Time
	Source      string
}

type RateLimitter struct {
    mu sync.Mutex
    rate int // tokens per second
    burst int // maximum burst capacity
    tokens float64 // current token account
    lastRefill time.Time
}

func NewRateLimitter(rate, burst int) *RateLimitter {
    return &RateLimitter{
        rate: rate,
        burst: burst,
        tokens: float64(burst),
        lastRefill: time.Now(),
    }
}

func (rl *RateLimitter) Allow() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    additional := float64(rl.rate) * time.Since(rl.lastRefill).Seconds()
    rl.tokens = math.Min(rl.tokens + additional, float64(rl.burst))
    rl.lastRefill = time.Now()

    if rl.tokens >= 1.0 {
        rl.tokens -= 1.0
        return true
    }

    return false
}

func (rl *RateLimitter) Wait(ctx context.Context) error {
    for {
        if (rl.Allow()) {
            return nil
        }
    
        rl.mu.Lock()
        waitSec := (1.0 - rl.tokens) / float64(rl.rate)
        if waitSec < 0 {
            waitSec = 0
        }
        rl.mu.Unlock()
        
        duration := time.Duration(waitSec * float64(time.Second))
        cwt, cancel := context.WithTimeout(ctx, duration)
        defer cancel()
        
        timer := time.NewTimer(duration)
        for {
            select {
                case <-cwt.Done():
                    if !timer.Stop() {
                        <-timer.C
                    }
                    return cwt.Err()
                case <-timer.C:
            }
        }
    }
}

// ContentAggregator manages the concurrent fetching and processing of content
type ContentAggregator struct {
	fetcher ContentFetcher
	processor ContentProcessor
	workerCount int
	requestsPerSecond int
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
    if fetcher == nil || processor == nil || workerCount <= 0 || requestsPerSecond <= 0 {
        return nil
    }
	return &ContentAggregator{
	    fetcher: fetcher,
	    processor: processor,
	    workerCount: workerCount,
	    requestsPerSecond: requestsPerSecond,
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	results, errs := ca.fanOut(ctx, urls)
	if len(errs) > 0 {
	    return nil, errs[0]
	}
	return results, nil
}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	return nil
}

// workerPool implements a worker pool pattern for processing content
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errors chan<- error,
) {
    var wg sync.WaitGroup
    wg.Add(ca.workerCount)
    
    rl := NewRateLimitter(ca.requestsPerSecond, ca.requestsPerSecond)

	for i := 0; i < ca.workerCount; i++ {
	    go func() {
	        defer wg.Done()

            for {
                select {
                    case url, ok := <-jobs: {
                        if !ok {
                            return
                        }
    
                        if err := rl.Wait(ctx); err != nil {
                            if ctx.Err() != nil {
                                return
                            }
                            errors <- err
                            continue
                        }
    
                        body, err := ca.fetcher.Fetch(ctx, url)
                        if err != nil {
                            errors <- err
                            continue
                        }
                        
                        data, err := ca.processor.Process(ctx, body)
                        if err != nil {
                            errors <- err
                            continue
                        }
                        
                        data.Source = url
                        results <- data
                    }
                    case <-ctx.Done():
                        return
                }
            }
	    }()
	}

	wg.Wait()
	close(results)
	close(errors)
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
    resultsData := make([]ProcessedData, 0, len(urls))
    resultsError := make([]error, 0, len(urls))
    
	jobs := make(chan string)
	results := make(chan ProcessedData, len(urls))
	errors := make(chan error, len(urls))

	go ca.workerPool(ctx, jobs, results, errors)
	
	go func() {
        for _, url := range urls {
	        select {
                case <-ctx.Done():
                    return
                case jobs <- url:
	        }
	    }
	    close(jobs)
    }()
	
	for results != nil || errors != nil {
	    select {
	        case result, ok := <-results: {
	            if !ok {
	                results = nil
	                continue
	            }
	            resultsData = append(resultsData, result)
	        }
	        case err, ok := <-errors: {
	            if !ok {
	                errors = nil
	                continue
	            }
	            resultsError = append(resultsError, err)
	        }
	    }
	} 

	return resultsData, resultsError
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
}

const timeout = 10 * time.Second
const maxBodyBytesRead = 2 << 20

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
	    return nil, err
	}
	
    req.Header.Set("User-Agent", "challenge11-aggregator/1.0")
	
	if hf.Client == nil {
	    hf.Client = &http.Client{
            Timeout: timeout,
        }
	}
	
	resp, err := hf.Client.Do(req)
	if err != nil {
	    return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
	    return nil, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
 	}
	
    return io.ReadAll(io.LimitReader(resp.Body, maxBodyBytesRead))
}

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct {}

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
    if len(content) == 0 {
        return ProcessedData{}, errors.New("empty HTML page")
    }
    
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
	    return ProcessedData{}, err
	}
	
	title := findValue(doc, "title")
	description := findValue(doc, "description")

	var keywords []string
	rawKW := findValue(doc, "keywords")
	parts := strings.Split(rawKW, ",")
	for _, part := range parts {
	    part = strings.TrimSpace(part)
	    if part != "" {
	        keywords = append(keywords, part)
	    }
	}
	
	if title == "" || description == "" || len(keywords) == 0 {
	    return ProcessedData{}, errors.New("invalid HTML page")
	}
	
	return ProcessedData{
	    Title: title,
	    Description: description,
	    Keywords: keywords,
	    Timestamp: time.Now(),
	}, nil
}

func findValue(n *html.Node, nodeName string) string {
	if n.Type == html.ElementNode && n.Data == nodeName {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			return n.FirstChild.Data
		}
	}
	if n.Type == html.ElementNode && n.Data == "meta" {
        var name, content string
        for _, a := range n.Attr {
            if a.Key == "name" {
                name = a.Val
            }
            if a.Key == "content" {
                content = a.Val
            }
        }
        if name == nodeName && content != "" {
            return content
        }
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if val := findValue(c, nodeName); val != "" {
			return val
		}
	}
	return ""
}