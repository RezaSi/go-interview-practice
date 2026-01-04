// package main

package challenge11

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/time/rate"
)

type ContentFetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}
type ContentProcessor interface {
	Process(ctx context.Context, content []byte) (ProcessedData, error)
}
type ProcessedData struct {
	Title       string
	Description string
	Keywords    []string
	Timestamp   time.Time
	Source      string
}
type ContentAggregator struct {
	fetcher      ContentFetcher
	processor    ContentProcessor
	workerCount  int
	limiter      *rate.Limiter
	wg           sync.WaitGroup
	shutdown     chan struct{}
	shutdownOnce sync.Once
}

// function signature is part of the assignment
func NewContentAggregator(fetcher ContentFetcher, processor ContentProcessor, workerCount int, requestsPerSecond int) *ContentAggregator {
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
		fetcher:     fetcher,
		processor:   processor,
		workerCount: workerCount,
		limiter:     rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
		shutdown:    make(chan struct{}),
	}
}

func (ca *ContentAggregator) FetchAndProcess(ctx context.Context, urls []string) ([]ProcessedData, error) {
	jobs := make(chan string, len(urls))
	results := make(chan ProcessedData, len(urls))
	errors := make(chan error, len(urls))
	var allResults []ProcessedData
	var allErrors []error
	ca.workerPool(ctx, jobs, results, errors)
	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case jobs <- url:
			case <-ctx.Done():
				return
			}
		}
	}()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < len(urls); i++ {
			select {
			case result := <-results:
				allResults = append(allResults, result)
			case err := <-errors:
				allErrors = append(allErrors, err)
			case <-ctx.Done():
				return
			}
		}
	}()
	// Wait for completion or context cancellation
	select {
	case <-done:
		if len(allErrors) > 0 {
			return allResults, fmt.Errorf("encountered %d errors: %v", len(allErrors), allErrors)
		}
		return allResults, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (ca *ContentAggregator) Shutdown() error {
	ca.shutdownOnce.Do(func() {
		close(ca.shutdown)
		ca.wg.Wait() // Wait for all workers to finish
	})
	return nil
}
func (ca *ContentAggregator) workerPool(ctx context.Context, jobs <-chan string, results chan<- ProcessedData, errors chan<- error) {
	for i := 0; i < ca.workerCount; i++ {
		ca.wg.Add(1)
		go func() {
			defer ca.wg.Done()
			for {
				select {
				case url, ok := <-jobs:
					if !ok {
						return // channel closed
					}
					if err := ca.limiter.Wait(ctx); err != nil {
						errors <- fmt.Errorf("rate limit error for %s: %w", url, err)
						continue
					}
					content, err := ca.fetcher.Fetch(ctx, url)
					if err != nil {
						errors <- fmt.Errorf("fetch error for %s: %w", url, err)
						continue
					}
					processedData, err := ca.processor.Process(ctx, content)
					if err != nil {
						errors <- fmt.Errorf("process error for %s: %w", url, err)
						continue
					}
					processedData.Source = url
					processedData.Timestamp = time.Now()
					results <- processedData
				case <-ctx.Done():
					return
				case <-ca.shutdown:
					return
				}
			}
		}()
	}
}

type HTTPFetcher struct {
	Client *http.Client
}

func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("empty url")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ContentAggregator/1.0")
	resp, err := hf.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}
	// Limit response size to 10MB to prevent memory exhaustion
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}
	return body, nil
}

type HTMLProcessor struct{}

func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	select {
	case <-ctx.Done():
		return ProcessedData{}, ctx.Err()
	default:
		if len(content) == 0 {
			return ProcessedData{}, fmt.Errorf("HTML parsing failed. Empty HTML.")
		}
		doc, err := html.Parse(bytes.NewReader(content))
		if err != nil {
			return ProcessedData{}, fmt.Errorf("HTML parsing failed: %w", err)
		}
		var title, description string
		var keywords []string
		var extractData func(*html.Node)
		extractData = func(n *html.Node) {
			if n.Type == html.ElementNode {
				switch n.Data {
				case "title":
					if n.FirstChild != nil {
						title = n.FirstChild.Data
					}
				case "meta":
					var name, content string
					for _, attr := range n.Attr {
						switch attr.Key {
						case "name":
							name = attr.Val
						case "content":
							content = attr.Val
						}
					}
					switch name {
					case "description":
						description = content
					case "keywords":
						keywords = splitKeywords(content)
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				extractData(c)
			}
		}
		extractData(doc)
		if title == "" && description == "" {
			return ProcessedData{}, fmt.Errorf("HTML parsing failed. Invalid HTML.")
		}
		return ProcessedData{Title: title, Description: description, Keywords: keywords}, nil
	}
}

func splitKeywords(keywordsStr string) []string {
	var result []string
	keywordsStr = strings.ReplaceAll(keywordsStr, ";", ",")
	for _, kw := range strings.Split(keywordsStr, ",") {
		if trimmed := strings.TrimSpace(kw); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
func main() {
	fetcher := &HTTPFetcher{http.DefaultClient}
	processor := &HTMLProcessor{}
	aggregator := NewContentAggregator(fetcher, processor, 3, 2)
	if aggregator == nil {
		log.Fatal("Failed to create content aggregator")
	}
	defer aggregator.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	urls := []string{"", "https://example.org", "https://example.net", "https://example.net/error"}
	startTime := time.Now()
	results, err := aggregator.FetchAndProcess(ctx, urls)
	processingTime := time.Since(startTime)
	if err != nil {
		log.Printf("Error in aggregate operation: %v", err)
	}
	for _, result := range results {
		fmt.Println(processingTime, result.Source, result.Title, result.Description, result.Keywords)
	}
}
