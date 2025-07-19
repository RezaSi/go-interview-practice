// Package challenge11 contains the solution for Challenge 11.
package challenge11

import (
	"context"
	"net/http"
	"time"
	"sync"
	"errors"
	"fmt"
	"strings"
	"io"

	"golang.org/x/time/rate"
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

// ContentAggregator manages the concurrent fetching and processing of content
type ContentAggregator struct {
	fetcher          ContentFetcher
	processor        ContentProcessor
	workerCount      int
	rateLimiter      *rate.Limiter
	wg               sync.WaitGroup
	shutdown         chan struct{}
	mu               sync.RWMutex
	isShuttingDown   bool
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	if workerCount <= 0 || requestsPerSecond <= 0 {
		return nil
	}

	if fetcher == nil || processor == nil {
		return nil
	}

	return &ContentAggregator{
		fetcher:          fetcher,
		processor:        processor,
		workerCount:      workerCount,
		rateLimiter:      rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
		shutdown:         make(chan struct{}),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	if ca.isShuttingDown {
		return nil, errors.New("shutting down")
	}

	results, errs := ca.fanOut(ctx, urls)
	if len(errs) > 0 {
		return results, fmt.Errorf("got %d errors", len(errs))
	}
	return results, nil
}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	if ca.isShuttingDown {
		return nil
	}
	ca.isShuttingDown = true
	close(ca.shutdown)
	ca.wg.Wait()
	return nil
}

// workerPool implements a worker pool pattern for processing content
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errors chan<- error,
) {
	defer ca.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ca.shutdown:
			return
		case url, ok := <-jobs:
			if ! ok {
				return
			}
			if err := ca.rateLimiter.Wait(ctx); err != nil {
				select {
				case errors <- fmt.Errorf("rate limiter error for %s: %v", url, err):
				case <-ctx.Done():
				case <-ca.shutdown:
				}
				continue
			}

			content, err := ca.fetcher.Fetch(ctx, url)
			if err != nil {
				select {
				case errors <- fmt.Errorf("fetch error for %s: %v", url, err):
				case <-ctx.Done():
				case <-ca.shutdown:
				}
				continue
			}

			data, err := ca.processor.Process(ctx, content)
			if err != nil {
				select {
				case errors <- fmt.Errorf("processing error for %s: %v", url, err):
				case <-ctx.Done():
				case <-ca.shutdown:
				}
				continue
			}

			data.Source = url
			data.Timestamp = time.Now()

			select {
			case results <- data:
			case <-ctx.Done():
			case <-ca.shutdown:
			}
		}
	}
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	jobs := make(chan string, len(urls))
	results := make(chan ProcessedData, len(urls))
	errs := make(chan error, len(urls))

	ca.wg.Add(ca.workerCount)
	for range(ca.workerCount) {
		go ca.workerPool(ctx, jobs, results, errs)
	}

	// Send jobs
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

	var data []ProcessedData
	var errors []error
	done := make(chan struct{})

	go func() {
		for {
			select {
			case result := <-results:
				data = append(data, result)
			case err := <-errs:
				errors = append(errors, err)
			case <-ctx.Done():
				close(done)
				return
			case <-ca.shutdown:
				close(done)
				return
			default:
				if len(data)+len(errors) == len(urls) {
					close(done)
					return
				}
			}
		}
	}()

	<-done
	return data, errors
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %v", err)
	}

	resp, err := hf.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %v", err)
	}

	return body, nil
}

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct {}

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return ProcessedData{}, err
	}

	var (
		title       string
		description string
		keywords    []string
	)

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil {
					title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "meta":
				var name, content string
				for _, attr := range n.Attr {
					switch strings.ToLower(attr.Key) {
					case "name":
						name = strings.ToLower(attr.Val)
					case "content":
						content = attr.Val
					}
				}
				switch name {
				case "description":
					description = strings.TrimSpace(content)
				case "keywords":
					keywords = strings.Split(strings.TrimSpace(content), ",")
					for i := range keywords {
						keywords[i] = strings.TrimSpace(keywords[i])
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if title == "" {
		return ProcessedData{}, errors.New("title not found")
	}

	return ProcessedData{
		Title:       title,
		Description: description,
		Keywords:    keywords,
		Timestamp:   time.Now().UTC(),
	}, nil
} 
