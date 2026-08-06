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
	// TODO: Initialize the ContentAggregator with the provided components
	if fetcher == nil {
	    return nil
	}
	if processor == nil {
	    return nil
	}
	if workerCount<1 {
	    return nil
	}
	if requestsPerSecond <1 {
	    return nil
	}
	return &ContentAggregator{
	    fetcher,
	    processor,
	    workerCount,
	    requestsPerSecond,
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	// TODO: Implement concurrent fetching and processing with proper error handling
	processedData, errs := ca.fanOut(ctx, urls)
	return processedData, errors.Join(errs...)
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
	// TODO: Implement worker pool logic
	for {
    	select {
    	    case <-ctx.Done():
    	        // fmt.Printf("Worker closing (ctx closed)\n")
    	        return
    	    case url, ok := <-jobs:
    	        if !ok {
    	            // fmt.Printf("Worker closing (jobs channel closed)\n")
    	            return
    	        } 
    	        // fmt.Printf("Fetching content for %s\n", url)
    	        content, err := ca.fetcher.Fetch(ctx, url)
    	        if err != nil {
    	            errors<-err
    	            continue
    	        }
    	        // fmt.Printf("Processing content %s for %s\n", string(content), url)
    	        processedData, err := ca.processor.Process(ctx, content)
    	        if err != nil {
    	            errors<-err
    	            continue
    	        }
    	        // fmt.Printf("Sending processed data for %s\n", url)
    	        results<-processedData
    	        continue
    	}
	}
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	// TODO: Implement fan-out, fan-in pattern
	
	jobs := make(chan string)
	results := make(chan ProcessedData)
	errors := make(chan error)
	
	// WaitGroup to wait for workers to finish
	var wg sync.WaitGroup
	
	// Spin up worker gorountines
	for i:=0; i< ca.workerCount;i++ {
	    wg.Go(func() {
	        ca.workerPool(ctx, jobs, results, errors)
	   })
	}
	
	// Close results after all workers finish
	go func() {
	    wg.Wait()
	    // fmt.Printf("Closing Channels (results/errors)\n")
	    close(results)
	    close(errors) // TODO - this may leak
	}()
	
	// Send jobs to the job channel
	go func() {
    	for _,url := range urls {
    	    jobs <- url
    	}
    	
	    // Close the job channel to indicate no more jobs will be sent
	    // fmt.Printf("Closing Channels (jobs)\n")
	    close(jobs)
	} ()
	
	var resultData []ProcessedData
	var resultErrors []error
	
	for {
    	select {
    	    case <-ctx.Done():
    	        // TODO
    	        return resultData, resultErrors
    	    case data, ok := <-results:
    	        if !ok {
    	            return resultData, resultErrors
    	        }
    	        // fmt.Printf("Received result %v\n", data)
    	        resultData = append(resultData, data)
    	        continue
    	    case dataError, ok := <-errors:
    	        if !ok {
    	            return resultData, resultErrors
    	        }
    	        // fmt.Printf("Received error %v\n", dataError)
    	        resultErrors = append(resultErrors, dataError)
    	        continue
    	}
	}
	
	return resultData, resultErrors
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
	// TODO: Add fields for rate limiting, etc.
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	// TODO: Implement HTTP-based content fetching with context support
	// fmt.Println("Fetch")
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
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
	    return nil, err
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
	// fmt.Println("Process")
	
	r := bytes.NewReader(content)
	
	fmt.Printf("parsing %s\n", string(content))
	doc, err := html.Parse(r)
	if err != nil {
	    fmt.Println("parse error %w", err)
		return ProcessedData{}, err
	}
	
	title, err := hp.findTitle(doc)
	if err != nil {
	    return ProcessedData{}, err
	}
	description := hp.findDescription(doc)
	keywords := strings.Split(hp.findKeywords(doc), ",")
	
	return ProcessedData{
	    Title: title,
	    Description: description,
	    Keywords: keywords,
	}, nil
} 

func (hp *HTMLProcessor) findTitle(n *html.Node) (string, error) {
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		return n.FirstChild.Data, nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title, _ := hp.findTitle(c)
		if title != "" {
			return title, nil
		}
	}

	return "", fmt.Errorf("no title found")
}

func (hp *HTMLProcessor) findDescription(doc *html.Node) string {
	var walk func(*html.Node) string
	walk = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string

			for _, a := range n.Attr {
				switch strings.ToLower(a.Key) {
				case "name":
					name = strings.ToLower(a.Val)
				case "content":
					content = a.Val
				}
			}

			if name == "description" {
				return content
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if desc := walk(c); desc != "" {
				return desc
			}
		}

		return ""
	}

	return walk(doc)
}

func (hp *HTMLProcessor) findKeywords(doc *html.Node) string {
	var walk func(*html.Node) string
	walk = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string

			for _, a := range n.Attr {
				switch strings.ToLower(a.Key) {
				case "name":
					name = strings.ToLower(a.Val)
				case "content":
					content = a.Val
				}
			}

			if name == "keywords" {
				return content
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if desc := walk(c); desc != "" {
				return desc
			}
		}

		return ""
	}

	return walk(doc)
}