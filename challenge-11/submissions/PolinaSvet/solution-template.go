package challenge11

//package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/time/rate"
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
	limiter      *rate.Limiter
	wg           sync.WaitGroup
	shutdown     chan struct{}
	shutdownOnce sync.Once
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
		fetcher:     fetcher,
		processor:   processor,
		workerCount: workerCount,
		limiter:     rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
		shutdown:    make(chan struct{}),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {

	jobs := make(chan string, len(urls))
	results := make(chan ProcessedData, len(urls))
	errors := make(chan error, len(urls))
	// Start workers
	ca.workerPool(ctx, jobs, results, errors)
	// Send jobs
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
	// Collect results
	// Implementation here...

	// Create channels for collecting results
	var allResults []ProcessedData
	var allErrors []error
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
		// All URLs processed
		if allErrors != nil {
			return allResults, fmt.Errorf("errors")
		}
		return allResults, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}

}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	ca.shutdownOnce.Do(func() {
		close(ca.shutdown)
		ca.wg.Wait() // Wait for all workers to finish
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
		ca.wg.Add(1)
		go func() {
			defer ca.wg.Done()
			for {
				select {
				case url, ok := <-jobs:
					if !ok {
						return
					}
					// Process URL here

					if err := ca.limiter.Wait(ctx); err != nil {
						errors <- fmt.Errorf("rate limit error for %s: %w", url, err)
						continue
					}

					// Получаем контент
					content, err := ca.fetcher.Fetch(ctx, url)
					if err != nil {
						errors <- fmt.Errorf("fetch error for %s: %w", url, err)
						continue
					}

					// Обрабатываем контент
					processed, err := ca.processor.Process(ctx, content)
					if err != nil {
						errors <- fmt.Errorf("process error for %s: %w", url, err)
						continue
					}

					processed.Source = url
					processed.Timestamp = time.Now()

					results <- processed

				case <-ctx.Done():
					return
				}
			}
		}()
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
		return nil, fmt.Errorf("create request failed: %w", err)
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

	// Используем io.ReadAll вместо resp.Body.Read
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	return content, nil

}

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct {
}

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	select {
	case <-ctx.Done():
		return ProcessedData{}, ctx.Err()
	default:

		if len(content) <= 0 {
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

		return ProcessedData{
			Title:       title,
			Description: description,
			Keywords:    keywords,
		}, nil
	}
}

func splitKeywords(keywordsStr string) []string {
	var keywords []string
	start := 0
	for i, char := range keywordsStr {
		if char == ',' || char == ';' {
			if i > start {
				keywords = append(keywords, keywordsStr[start:i])
			}
			start = i + 1
		}
	}
	if start < len(keywordsStr) {
		keywords = append(keywords, keywordsStr[start:])
	}
	return keywords
}

// =====================================
func main() {
	fmt.Println("🚀 Starting Content Aggregator")
	fmt.Println("===================================")

	httpFetcher := &HTTPFetcher{Client: http.DefaultClient}
	htmlProcessor := &HTMLProcessor{}

	aggregator := NewContentAggregator(
		httpFetcher,
		htmlProcessor,
		3,
		2,
	)

	defer aggregator.Shutdown()

	urls := []string{
		"https://httpbin.org/html",
		"https://httpbin.org/status/404",
		"https://httpbin.org/delay/2",
		"https://httpbin.org/xml",
		"https://example.net",
	}

	fmt.Printf("📋 Processing %d URLs...\n", len(urls))
	fmt.Println("⏰ Setting timeout to 15 seconds")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime := time.Now()
	results, err := aggregator.FetchAndProcess(ctx, urls)
	processingTime := time.Since(startTime)

	fmt.Printf("\n📊 Processing completed in %v\n", processingTime)
	fmt.Println("===================================")

	if err != nil {
		fmt.Printf("⚠️  Completed with errors: %v\n", err)
	} else {
		fmt.Println("✅ All URLs processed successfully!")
	}

	fmt.Printf("\n📈 Results Summary:\n")
	fmt.Printf("   Total URLs: %d\n", len(urls))
	fmt.Printf("   Successfully processed: %d\n", len(results))
	if err != nil {
		if aggErr, ok := err.(interface{ Unwrap() []error }); ok {
			fmt.Printf("   Errors: %d\n", len(aggErr.Unwrap()))
		}
	}

	fmt.Println("\n🔍 Detailed Results:")
	for i, result := range results {
		fmt.Printf("%d. Source: %s\n", i+1, result.Source)
		fmt.Printf("   Title: %s\n", truncate(result.Title, 50))
		fmt.Printf("   Description: %s\n", truncate(result.Description, 70))
		fmt.Printf("   Keywords: %v\n", truncateKeywords(result.Keywords, 3))
		fmt.Printf("   Timestamp: %s\n", result.Timestamp.Format("15:04:05"))
		fmt.Println("   ---")
	}

	fmt.Println("\n✨ Demo completed!")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func truncateKeywords(keywords []string, maxCount int) []string {
	if len(keywords) <= maxCount {
		return keywords
	}
	return append(keywords[:maxCount], "...")
}
