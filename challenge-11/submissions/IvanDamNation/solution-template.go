// Package challenge11 contains the solution for Challenge 11.
package challenge11

import (
    "bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
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

// ContentAggregator manages the concurrent fetching and processing of content
type ContentAggregator struct {
	Fetcher     ContentFetcher
	Processor   ContentProcessor
	WorkerCount int
	RateLimiter *time.Ticker
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	if fetcher == nil || processor == nil || workerCount < 1 || requestsPerSecond < 1 {
	    return nil
	}
	
	return &ContentAggregator{
	    Fetcher: fetcher,
	    Processor: processor,
	    WorkerCount: workerCount,
	    RateLimiter: time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	res, errs := ca.fanOut(ctx, urls)
	if len(errs) != 0 {
	    return res, errors.Join(errs...)
	}
	return res, nil
}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	ca.RateLimiter.Stop()
	return nil
}

// workerPool implements a worker pool pattern for processing content
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errors chan<- error,
) {
	for url := range jobs {
	    select {
	    case <-ctx.Done():
	        return
	    case <-ca.RateLimiter.C:
	    }
	    
	    data, err := ca.Fetcher.Fetch(ctx, url)
        if err != nil {
            select {
            case <-ctx.Done():
                return
            case errors <- err:
            }
            continue
        }
        
        res, err := ca.Processor.Process(ctx, data)
        if err != nil {
            select {
            case <-ctx.Done():
                return
            case errors <- err:
            }
            continue
        }
        
        select {
        case <-ctx.Done():
            return
        case results <- res:
        }
	}
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	jobCh := make(chan string, len(urls))
	resCh := make(chan ProcessedData, len(urls))
	errCh := make(chan error, len(urls))
	
	var wg sync.WaitGroup
	wg.Add(ca.WorkerCount)
	for range ca.WorkerCount {
	    go func() {
	        defer wg.Done()
	        ca.workerPool(ctx, jobCh, resCh, errCh)
	    }()
	}
	
	go func() {
	    defer close(jobCh)
	    for _, url := range urls {
    	    select {
    	    case <-ctx.Done():
    	        return
    	    default:
    	        jobCh <- url
    	    }
    	}
	}()
	
	go func() {
	    wg.Wait()
	    close(resCh)
	    close(errCh)
	}()
	
	results := make([]ProcessedData, 0, len(urls))
	errs := make([]error, 0, len(urls))
	
	for range len(urls) {
	    select {
	    case <-ctx.Done():
	        errs = append(errs, ctx.Err())
	        return results, errs
	    case res := <-resCh:
	        results = append(results, res)
	    case err := <-errCh:
	        errs = append(errs, err)
	    }
	}
	
	return results, errs
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
	    return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	client := hf.Client
	if client == nil {
	    client = http.DefaultClient
	}
	
	resp, err := client.Do(req)
	if err != nil {
	    return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
	    return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
	    return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	return buf.Bytes(), nil
}

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct {}

var (
	titleRegex       = regexp.MustCompile(`(?i)<title[^>]*>([^<]+)</title>`)
	descriptionRegex = regexp.MustCompile(`(?i)<meta\s+[^>]*name=["']description["'][^>]*content=["']([^"']+)["']`)
	keywordsRegex    = regexp.MustCompile(`(?i)<meta\s+[^>]*name=["']keywords["'][^>]*content=["']([^"']+)["']`)
)

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	if err := ctx.Err(); err != nil {
	    return ProcessedData{}, err
	}
	
	if len(content) == 0 {
	    return ProcessedData{}, fmt.Errorf("empty HTML content")
	}
	
	htmlStr := string(content)
	titleMatches := titleRegex.FindStringSubmatch(htmlStr)
	
	if len(titleMatches) <= 1 {
	    return ProcessedData{}, fmt.Errorf("invalid or malformed HTML: missing title")
	}
	
	var data ProcessedData
	data.Timestamp = time.Now()
	data.Source = ""
	data.Title = strings.TrimSpace(titleMatches[1])
	
	if descMatches := descriptionRegex.FindStringSubmatch(htmlStr); len(descMatches) > 1 {
	    data.Description = strings.TrimSpace(descMatches[1])
	}
	
	if kwMatches := keywordsRegex.FindStringSubmatch(htmlStr); len(kwMatches) > 1 {
	    rawKeywords := strings.Split(kwMatches[1], ",")
	    for _, kw := range rawKeywords {
	        if trimmed := strings.TrimSpace(kw); trimmed != "" {
	            data.Keywords = append(data.Keywords, trimmed)
	        }
	    }
	}
	
	return data, nil
}
