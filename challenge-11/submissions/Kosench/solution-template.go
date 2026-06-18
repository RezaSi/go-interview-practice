// Package challenge11 contains the solution for Challenge 11.
package challenge11

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	// Add any necessary imports here
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

type fetchResult struct {
	url     string
	content []byte
	err     error
}

// ContentAggregator manages the concurrent fetching and processing of content
type ContentAggregator struct {
	fetcher        ContentFetcher
	processor      ContentProcessor
	workerCount    int
	limiter        *rate.Limiter
	shutdown       chan struct{}
	wg             sync.WaitGroup
	mu             sync.RWMutex
	isShuttingDown bool
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	// Validate parameters
	if fetcher == nil || processor == nil {
		return nil
	}

	if workerCount <= 0 {
		return nil
	}

	if requestsPerSecond <= 0 {
		return nil
	}

	return &ContentAggregator{
		fetcher:     fetcher,
		processor:   processor,
		workerCount: workerCount,
		limiter:     rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
		shutdown:    make(chan struct{}),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(ctx context.Context, urls []string) ([]ProcessedData, error) {
	ca.mu.RLock()
	if ca.isShuttingDown {
		ca.mu.RUnlock()
		return nil, errors.New("aggregator is shutting down")
	}
	ca.mu.RUnlock()

	// Track this operation
	ca.wg.Add(1)
	defer ca.wg.Done()

	result, errs := ca.fanOut(ctx, urls)

	if len(errs) > 0 {
		return result, errs[0]
	}

	return result, nil

}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	ca.mu.Lock()
	if ca.isShuttingDown {
		ca.mu.Unlock()
		return nil
	}
	ca.isShuttingDown = true
	ca.mu.Unlock()

	// Signal shutdown to workers
	close(ca.shutdown)

	// Wait for in-flight operations to complete
	ca.wg.Wait()

	return nil
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(ctx context.Context, urls []string) ([]ProcessedData, []error) {
	jobs := make(chan string, len(urls))

	results := make(chan ProcessedData, len(urls))
	errors := make(chan error, len(urls))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ca.workerPool(ctx, jobs, results, errors)
	}()

	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case jobs <- url:
			case <-ctx.Done():
				return
			case <-ca.shutdown:
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	//Fan-In
	var allResult []ProcessedData
	var allErrors []error

	for results != nil || errors != nil {
		select {
		case result, ok := <-results:
			if !ok {
				results = nil
			} else {
				allResult = append(allResult, result)
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
			} else {
				allErrors = append(allErrors, err)
			}
		}
	}

	return allResult, allErrors
}

// workerPool implements a worker pool pattern for processing content
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errors chan<- error,
) {
	var wg sync.WaitGroup

	for i := 0; i < ca.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ca.shutdown:
					return
				case url, ok := <-jobs:
					if !ok {
						return
					}
					ca.processURL(ctx, url, results, errors)
				}
			}

		}(i)
	}

	wg.Wait()
}

// ===== Helper method for processing a single URL =====
func (ca *ContentAggregator) processURL(ctx context.Context, url string, result chan<- ProcessedData, errors chan<- error) {
	if err := ca.limiter.Wait(ctx); err != nil {
		select {
		case errors <- fmt.Errorf("rate limit error for %s: %w", url, err):
		case <-ctx.Done():
		case <-ca.shutdown:
		}
		return
	}

	content, err := ca.fetcher.Fetch(ctx, url)
	if err != nil {
		select {
		case errors <- fmt.Errorf("fetch error for %s: %w", url, err):
		case <-ctx.Done():
		case <-ca.shutdown:
		}
		return
	}

	processed, err := ca.processor.Process(ctx, content)
	if err != nil {
		select {
		case errors <- fmt.Errorf("processing error for %s: %w", url, err):
		case <-ctx.Done():
		case <-ca.shutdown:
		}
		return
	}

	processed.Source = url
	processed.Timestamp = time.Now()

	select {
	case result <- processed:
	case <-ctx.Done():
	case <-ca.shutdown:
	}

}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := hf.Client
	if client == nil {
		client = http.DefaultClient // ← используем стандартный клиент Go
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct{}

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	if len(content) == 0 {
		return ProcessedData{}, errors.New("empty content")
	}

	text := string(content)

	if !strings.Contains(text, "<") || !strings.Contains(text, ">") {
		return ProcessedData{}, errors.New("invalid HTML content")
	}

	title := hp.extractTitle(text)
	description := hp.extractDescription(text)
	keywords := hp.extractKeywords(text)

	if title == "No title" && description == "" && len(keywords) == 0 {
		return ProcessedData{}, errors.New("could not extract any data from HTML")
	}

	return ProcessedData{
		Title:       title,
		Description: description,
		Keywords:    keywords,
	}, nil
}

func (hp *HTMLProcessor) extractTitle(html string) string {
	start := strings.Index(html, "<title>")
	if start == -1 {
		return "No title"
	}
	start += 7
	end := strings.Index(html[start:], "</title>")
	if end == -1 {
		return "No title"
	}
	return html[start : start+end]
}

func (hp *HTMLProcessor) extractDescription(html string) string {
	// Search for <meta name="description" content="...">
	metaTag := `<meta name="description" content="`
	quoteChar := `"`
	start := strings.Index(html, metaTag)
	if start == -1 {
		metaTag = `<meta name='description' content='`
		quoteChar = `'`
		start = strings.Index(html, metaTag)
		if start == -1 {
			return ""
		}
	}

	start += len(metaTag)
	end := strings.Index(html[start:], quoteChar)
	if end == -1 {
		return ""
	}

	return html[start : start+end]
}

func (hp *HTMLProcessor) extractKeywords(html string) []string {
	// Search for <meta name="keywords" content="...">
	metaTag := `<meta name="keywords" content="`
	quoteChar := `"`
	start := strings.Index(html, metaTag)
	if start == -1 {
		metaTag = `<meta name='keywords' content='`
		quoteChar = `'`
		start = strings.Index(html, metaTag)
		if start == -1 {
			return []string{}
		}
	}

	start += len(metaTag)
	end := strings.Index(html[start:], quoteChar)
	if end == -1 {
		return []string{}
	}

	keywordsStr := html[start : start+end]
	keywords := strings.Split(keywordsStr, ",")

	result := make([]string, 0, len(keywords))
	for _, k := range keywords {
		k = strings.TrimSpace(k)
		if k != "" {
			result = append(result, k)
		}
	}

	return result
}
