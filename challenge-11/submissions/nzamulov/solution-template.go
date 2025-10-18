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

    var additional = float64(rl.rate) * time.Since(rl.lastRefill).Seconds()
    rl.tokens = min(rl.tokens + additional, float64(rl.burst))
    rl.lastRefill = time.Now()

    if rl.tokens > 0 {
        rl.tokens--
        return true
    }

    return false
}

func (rl *RateLimitter) Wait(ctx context.Context) error {
    if (rl.Allow()) {
        return nil
    }

    var coef = time.Second.Seconds() / float64(rl.rate)
    var refill_diff = time.Since(rl.lastRefill).Seconds()
    var timeout = time.Duration(coef - refill_diff)
    
    cwt, cancel := context.WithTimeout(ctx, timeout * time.Second)
    defer cancel()
    
    for {
        select{
            case <-time.After(timeout * time.Second):
                return nil
            case <-cwt.Done():
                return cwt.Err()
        }
    }
    
    return nil
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
    
    rl := NewRateLimitter(1, 1)

	for range ca.workerCount {
	    go func() {
	        defer wg.Done()

            select {
                case url := <-jobs: {
                    rl.Wait(ctx)

                    body, err := ca.fetcher.Fetch(ctx, url)
                    if err != nil {
                        errors <- err
                        break
                    }
                    
                    data, err := ca.processor.Process(ctx, body)
                    if err != nil {
                        errors <- err
                        break
                    }
                    
                    results <- data
                }
                case <-ctx.Done():
                    errors <- ctx.Err()
                    return
            }
	    }()
	}

	wg.Wait()
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
    resultsData := make([]ProcessedData, 0, len(urls))
    resultsError := make([]error, 0, len(urls))
    
	jobs := make(chan string)
	results := make(chan ProcessedData)
	errors := make(chan error)

	go ca.workerPool(ctx, jobs, results, errors)

	var wg sync.WaitGroup
	wg.Add(len(urls))
	
	done := make(chan struct{})
	
	go func() {
	   for {
    	    select {
    	        case result := <-results: {
    	            resultsData = append(resultsData, result)
    	            wg.Done()
    	        }
    	        case err := <-errors: {
    	            resultsError = append(resultsError, err)
    	            wg.Done()
    	        }
    	        case <-done:
    	            return
    	    }
    	} 
	}()

	for _, url := range urls {
	    go func() {
	        jobs <- url
	    }()
	}
	
	wg.Wait()
	close(done)

	return resultsData, resultsError
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
	    return nil, err
	}
	
	resp, err := hf.Client.Do(req)
	if err != nil {
	    return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
	    return nil, errors.New("HTTP status is not OK")
	}
	
	return io.ReadAll(resp.Body)
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
	keywords := strings.Split(findValue(doc, "keywords"), ",")
	
	if title == "" || description == "" || len(keywords) == 0 {
	    return ProcessedData{}, errors.New("invalid HTML page")
	}
	
	return ProcessedData{
	    Title: title,
	    Description: description,
	    Keywords: keywords,
	}, nil
}

func findValue(n *html.Node, nodeName string) string {
	if n.Type == html.ElementNode && n.Data == nodeName {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			return n.FirstChild.Data
		}
	}
	if n.Type == html.ElementNode && n.Data == "meta" {
	    var extract = true

		for _, a := range n.Attr {
			if a.Key == "name" && a.Val != nodeName {
				extract = false
			}
			if a.Key == "content" && extract {
                return a.Val
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if val := findValue(c, nodeName); val != "" {
			return val
		}
	}
	return ""
}