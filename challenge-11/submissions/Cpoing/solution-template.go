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

	"golang.org/x/net/html"
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

// ContentAggregator manages the concurrent fetching and processing of content
type ContentAggregator struct {
	// TODO: Add fields for fetcher, processor, worker count, rate limiter, etc.
	fetcher        ContentFetcher
	processor      ContentProcessor
	workerCount    int
	requestLimiter *rate.Limiter
	wg             sync.WaitGroup
	shutdown       chan struct{}
	shutdownOnce   sync.Once
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	// TODO: Initialize the ContentAggregator with the provided components
	if workerCount <= 0 || requestsPerSecond <= 0 || fetcher == nil || processor == nil {
		return nil
	}

	return &ContentAggregator{
		fetcher:        fetcher,
		processor:      processor,
		workerCount:    workerCount,
		requestLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
		shutdown:       make(chan struct{}),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	// TODO: Implement concurrent fetching and processing with proper error handling

	if len(urls) == 0 {
		return nil, nil
	}

	result, errs := ca.fanOut(ctx, urls)
	if len(errs) > 0 {
		return result, errors.Join(errs...)
	}

	return result, nil
}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	// TODO: Implement proper shutdown logic
	ca.shutdownOnce.Do(func() {
		close(ca.shutdown)
		ca.wg.Wait()
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
	// TODO: Implement worker pool logic
	for i := 0; i < ca.workerCount; i++ {
		ca.wg.Add(1)
		go func() {
			defer ca.wg.Done()
			for {
				select {
				case <-ca.shutdown:
					return
				case <-ctx.Done():
					return
				case url, ok := <-jobs:
					if !ok {
						return
					}
					content, fetchErr := ca.fetcher.Fetch(ctx, url)
					if fetchErr != nil {
						select {
						case errors <- fmt.Errorf("fetch error for %s: %w", url, fetchErr):
						case <-ctx.Done():
						case <-ca.shutdown:
						}
						continue
					}

					processedData, err := ca.processor.Process(ctx, content)
					if err != nil {
						select {
						case errors <- fmt.Errorf("process error for %s: %w", url, err):
						case <-ctx.Done():
						case <-ca.shutdown:
						}
						continue
					}

					processedData.Source = url
					processedData.Timestamp = time.Now()

					select {
					case results <- processedData:
					case <-ctx.Done():
						return
					case <-ca.shutdown:
						return
					}
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
	var results []ProcessedData
	var errs []error

	jobs := make(chan string, len(urls))
	out := make(chan ProcessedData, len(urls))
	errors := make(chan error, len(urls))

	ca.workerPool(ctx, jobs, out, errors)

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

	done := make(chan struct{})
	go func() {
		ca.wg.Wait()
		close(out)
		close(errors)
		close(done)
	}()

	collecting := true
	for collecting {
		select {
		case r, ok := <-out:
			if ok {
				results = append(results, r)
			}
		case e, ok := <-errors:
			if ok && e != nil {
				errs = append(errs, e)
			}
		case <-done:
			collecting = false
		case <-ctx.Done():
			<-done
			collecting = false
		case <-ca.shutdown:
			<-done
			collecting = false
		}
	}

	return results, errs
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
	// TODO: Add fields for rate limiting, etc.
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	// TODO: Implement HTTP-based content fetching with context support

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	resp, err := hf.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
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
	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return ProcessedData{}, err
	}

	var data ProcessedData
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				data.Title = strings.TrimSpace(n.FirstChild.Data)
			}

			if n.Data == "meta" {
				var name, content string
				for _, a := range n.Attr {
					if strings.ToLower(a.Key) == "name" {
						name = strings.ToLower(a.Val)
					}
					if strings.ToLower(a.Key) == "content" {
						content = a.Val
					}
				}
				if name == "description" {
					data.Description = strings.TrimSpace(content)
				}
				if name == "keywords" {
					parts := strings.Split(strings.TrimSpace(content), ",")
					for i := range parts {
						parts[i] = strings.TrimSpace(parts[i])
					}
					data.Keywords = parts
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if strings.TrimSpace(data.Title) == "" {
		return ProcessedData{}, errors.New("invalid title")
	}

	data.Timestamp = time.Now()
	return data, nil
}
