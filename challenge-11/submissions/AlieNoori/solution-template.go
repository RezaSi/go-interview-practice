// Package challenge11 contains the solution for Challenge 11.
package challenge11

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

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
	fetcher      ContentFetcher
	processor    ContentProcessor
	workerCount  int
	shutdown     chan struct{}
	shutdownOnce sync.Once
	sync.WaitGroup
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	if fetcher == nil {
		return nil
	}
	if processor == nil {
		return nil
	}

	if workerCount <= 0 {
		return nil
	}

	if requestsPerSecond <= 0 {
		return nil
	}

	return &ContentAggregator{
		processor:   processor,
		fetcher:     fetcher,
		workerCount: workerCount,
		shutdown:    make(chan struct{}),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	// TODO: Implement concurrent fetching and processing with proper error handling

	processedData, errs := ca.fanOut(ctx, urls)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return processedData, nil
}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	// TODO: Implement proper shutdown logic
	ca.shutdownOnce.Do(func() {
		close(ca.shutdown)
		ca.WaitGroup.Wait()
	})
	return nil
}

// workerPool implements a worker pool pattern for processing content
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errors chan<- error,
) {
	for i := 0; i < ca.workerCount; i++ {
		ca.WaitGroup.Add(1)
		go func() {
			defer ca.WaitGroup.Done()
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}
					data, err := ca.fetcher.Fetch(ctx, job)
					if err != nil {
						select {
						case errors <- fmt.Errorf("fetcher error for %s: %w", job, err):
						case <-ctx.Done():
						}
						return
					}

					processedData, err := ca.processor.Process(ctx, data)
					if err != nil {
						select {
						case errors <- fmt.Errorf("fetcher error for %s: %w", job, err):
						case <-ctx.Done():
						}
						return
					}

					select {
					case results <- processedData:
					case <-ctx.Done():
					}

				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	// TODO: Implement fan-out, fan-in pattern

	results := struct {
		errors        []error
		processedData []ProcessedData
		sync.Mutex
	}{
		errors:        []error{},
		processedData: []ProcessedData{},
	}

	if len(urls) <= 0 {
		return results.processedData, nil
	}

	errChan := make(chan error, len(urls))
	resultChan := make(chan ProcessedData, len(urls))
	jobChan := make(chan string, len(urls))

	go func() {
		for _, url := range urls {
			select {
			case jobChan <- url:
			case <-ctx.Done():
				return
			}
		}
	}()

	ca.workerPool(ctx, jobChan, resultChan, errChan)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for range urls {
			select {
			case err := <-errChan:
				results.Lock()
				results.errors = append(results.errors, err)
				results.Unlock()
			case result := <-resultChan:
				results.Lock()
				results.processedData = append(results.processedData, result)
				results.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	select {
	case <-done:
		close(jobChan)
		close(resultChan)
	case <-ctx.Done():
		return nil, results.errors
	}

	return results.processedData, results.errors
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
	// TODO: Add fields for rate limiting, etc.
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := hf.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}

	return body, nil
}

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct {
	// TODO: Add any fields needed for HTML processing
}

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	// TODO: Implement HTML processing logic

	if len(content) <= 0 {
		return ProcessedData{}, fmt.Errorf("empty content")
	}

	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return ProcessedData{}, fmt.Errorf("error parsing html: %w", err)
	}

	pd := ProcessedData{
		Timestamp: time.Now(),
	}

	for n := range doc.Descendants() {
		switch n.Data {
		case "title":
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				pd.Title = strings.TrimSpace(n.FirstChild.Data)
			}

		case "meta":
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					if attr.Val == "description" {
						for _, attr := range n.Attr {
							if attr.Key == "content" {
								pd.Description = strings.TrimSpace(attr.Val)
							}
						}
					} else if attr.Val == "keywords" {
						for _, attr := range n.Attr {
							if attr.Key == "content" {
								pd.Keywords = strings.Split(attr.Val, ",")
							}
						}
					}
				}
			}
		}
	}

	if pd.Description == "" || pd.Keywords == nil {
		return ProcessedData{}, fmt.Errorf("invalid html")
	}

	return pd, nil
}
