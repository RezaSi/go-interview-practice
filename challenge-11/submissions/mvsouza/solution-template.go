// Package challenge11 contains the solution for Challenge 11.
package challenge11

import (
	"bytes"
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
	Fetcher           ContentFetcher
	Processor         ContentProcessor
	WorkerCount       int
	RequestsPerSecond int
	Limiter           *rate.Limiter
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	if fetcher == nil || processor == nil || requestsPerSecond <= 0 || workerCount <= 0 {
		return nil
	}
	return &ContentAggregator{
		Fetcher:           fetcher,
		Processor:         processor,
		WorkerCount:       workerCount,
		RequestsPerSecond: requestsPerSecond,
		Limiter:           rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	results, errs := ca.fanOut(ctx, urls)
	return results, joinErrors(errs)
}

type multiError struct {
	errors []error
}

// Join all errors into one
func joinErrors(errs []error) error {
	if len(errs) > 0 {
		return &multiError{errors: errs}
	}
	return nil
}

// Error implements the error interface
func (m *multiError) Error() string {
	var s []string
	for _, err := range m.errors {
		s = append(s, err.Error())
	}
	return strings.Join(s, "\n")
}

// Unwrap implements the interface to allow error unwrapping
func (m *multiError) Unwrap() []error {
	return m.errors
}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	// TODO: Implement proper shutdown logic
	return nil
}

// workerPool implements a worker pool pattern for processing content
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errors chan<- error,
) {
	for i := 0; i < ca.WorkerCount; i++ {
		go func() {
			for url := range jobs {
				if err := ca.Limiter.Wait(ctx); err != nil {
					errors <- err
					continue
				}
				content, err := ca.Fetcher.Fetch(ctx, url)
				if err != nil {
					errors <- err
					continue
				}
				result, err := ca.Processor.Process(ctx, content)
				if err != nil {
					errors <- err
					continue
				}
				results <- result
			}
		}()
	}
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	var wg sync.WaitGroup
	fetchesChannel := make(chan string, ca.WorkerCount)
	resultsChannel := make(chan ProcessedData, len(urls))
	errorsChannel := make(chan error, len(urls))
	go ca.workerPool(ctx, fetchesChannel, resultsChannel, errorsChannel)
	var errs []error
	var results []ProcessedData
	collector := make(chan bool)
	go func() {
		defer close(collector)
		for resultsChannel != nil || errorsChannel != nil {
			select {
			case result, ok := <-resultsChannel:
				if !ok {
					resultsChannel = nil
					continue
				}
				results = append(results, result)
				wg.Done()
			case err, ok := <-errorsChannel:
				if !ok {
					errorsChannel = nil
					continue
				}
				errs = append(errs, err)
				wg.Done()
			}
		}
	}()

	for _, url := range urls {
		wg.Add(1)
		fetchesChannel <- url
	}

	go func() {
		wg.Wait()
		close(resultsChannel)
		close(errorsChannel)
	}()

	close(fetchesChannel)
	<-collector

	return results, errs
}

var limiter *rate.Limiter

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) (body []byte, err error) {
	client := hf.Client
	if client == nil {
		client = http.DefaultClient
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}
	var resp *http.Response
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	for i := 0; i < 3; i++ {
		resp, err = client.Do(request)
		if err != nil {
			continue
		} else if !successResp(resp) {
			resp.Body.Close()
			continue
		} else {
			body, err = io.ReadAll(resp.Body)
			return
		}
	}
	if err != nil {
		return nil, err
	} else if !successResp(resp) {
		return nil, fmt.Errorf("Request not successful: stats code %d", resp.StatusCode)
	}
	return
}

func successResp(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct {
	// TODO: Add any fields needed for HTML processing
}

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	// TODO: Implement HTML processing logic
	result := ProcessedData{}
	if len(content) <= 0 {
		return result, errors.New("No content")
	}
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return result, err
	}
	hp.processHtml(&result, doc)
	if result.Title == "" {
		return result, errors.New("Invalid HTML")
	}
	return result, nil
}

func (hp *HTMLProcessor) processHtml(result *ProcessedData, doc *html.Node) {

	nodes := []*html.Node{doc}
	var n *html.Node
	for len(nodes) > 0 {
		n, nodes = nodes[0], nodes[1:]
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if title := n.FirstChild; title != nil {
					result.Title = title.Data
				}
			case "meta":
				if len(n.Attr) > 0 {
					var name string
					var content string
					for _, attr := range n.Attr {
						switch attr.Key {
						case "content":
							content = attr.Val
						case "name":
							name = attr.Val
						}
					}
					switch name {
					case "description":
						result.Description = strings.TrimSpace(content)
					case "keywords":
						for _, k := range strings.Split(content, ",") {
							result.Keywords = append(result.Keywords, strings.TrimSpace(k))
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			nodes = append(nodes, c)
		}
	}
}
