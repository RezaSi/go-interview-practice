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

	// Add any necessary imports here
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/time/rate"
)

// Interfaces
// ===========
// ContentFetcher defines an interface for fetching content from URLs
type ContentFetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// ContentProcessor defines an interface for processing raw content
type ContentProcessor interface {
	Process(ctx context.Context, content []byte) (ProcessedData, error)
}

// Custom errors
// =============
var (
	ErrBadParam        = errors.New("bad parameter")
	ErrNilReceiver     = errors.New("nil receiver")
	ErrBadUrl          = errors.New("bad or empty url")
	ErrTooManyFailures = errors.New("too many failures")
)

// ProcessedData represents structured data extracted from raw content
type ProcessedData struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Keywords    []string  `json:"keywords"`
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"`
}

// Aggregator - the main orchestrator
// =================================
// ContentAggregator manages the concurrent fetching and processing of content
type ContentAggregator struct {
	fetcher           ContentFetcher
	processor         ContentProcessor
	workerCount       int
	requestsPerSecond int
	wg                sync.WaitGroup
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	if fetcher == nil || processor == nil {
		return nil
	}
	if workerCount <= 0 || requestsPerSecond <= 0 {
		return nil
	}

	return &ContentAggregator{
		fetcher:           fetcher,
		processor:         processor,
		workerCount:       workerCount,
		requestsPerSecond: requestsPerSecond,
	}
}

// Aggregator methods
// ==================
// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	// TODO: Implement concurrent fetching and processing with proper error handling
	// TODO
	// ca.fetcher.limiter = rate.NewLimiter()
	jobs := make(chan string)
	results := make(chan ProcessedData)
	errors := make(chan error)

	// Write jobs (URLs) into jobs channel
	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	ca.workerPool(ctx, jobs, results, errors)

	var result []ProcessedData
	for res := range results {
		result = append(result, res)
	}

	for {
		select {
		case res, ok := <-results:
			if !ok {
				return result, nil
			}
			result = append(result, res)
		case <-errors:
		case <-ctx.Done():
			err := ca.Shutdown()
			return nil, err
		}
	}
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
	for i := 0; i < ca.workerCount; i++ {
		ca.wg.Add(1)
		go func() {
			defer ca.wg.Done()
			for job := range jobs {
				content, err := ca.fetcher.Fetch(ctx, job)
				if err != nil {
					errors <- err
					continue
				}
				processedContent, err := ca.processor.Process(ctx, content)
				if err != nil {
					errors <- err
					continue
				}
				results <- processedContent
			}
		}()
	}
	go func() {
		ca.wg.Wait()
		close(results)
		close(errors)
	}()
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	resChan := make(chan ProcessedData)
	errs := make(chan error)
	// TODO: Implement fan-out, fan-in pattern
	wg := sync.WaitGroup{}
	for _, url := range urls {
		url := url
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			content, err := ca.fetcher.Fetch(ctx, url)
			if err != nil {
				errs <- err
				return
			}
			processedContent, err := ca.processor.Process(ctx, content)
			if err != nil {
				errs <- err
				return
			}
			resChan <- processedContent
			close(resChan)
		}(url)
	}
	wg.Wait()
	close(errs)
	var processedData []ProcessedData
	for res := range resChan {
		processedData = append(processedData, res)
	}
	var errList []error
	for err := range errs {
		errList = append(errList, err)
	}
	return processedData, errList
}

// Fetcher
// ========
// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client  *http.Client
	Limiter *rate.Limiter
}

// Fetch retrieves content from a URL via HTTP not exccided rate limit
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	// Validation input parameters
	if hf == nil {
		return nil, fmt.Errorf("fetcher error: %w", ErrNilReceiver)
	}
	if url == "" {
		return nil, fmt.Errorf("fetcher error: %w", ErrBadUrl)
	}

	// New request with context added
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	// Limiter init
	// TODO refactoring
	limiter := rate.NewLimiter(rate.Limit(1), 5)
	hf.Limiter = limiter
	// Send the requst throug the limiter
	if err := hf.Limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("fetcher rate limiter: %w", err)
	}
	resp, err := hf.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}

	// Request body gefer close
	defer resp.Body.Close()

	// Read the body of request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body: %w", err)
	}

	//Return body as []byte
	return body, nil
}

// processor
// =========
// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct {
}

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	// Input parameter validation
	if hp == nil {
		return ProcessedData{}, fmt.Errorf("processor error: %w", ErrNilReceiver)
	}
	if len(content) == 0 {
		return ProcessedData{}, fmt.Errorf("processor error: %w", ErrBadParam)
	}

	// Data processing
	var res ProcessedData
	r := strings.NewReader(string(content))
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return ProcessedData{}, err
	}
	title := doc.Find("Title").Text()
	res.Title = title
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		content, _ := s.Attr("content")
		switch {
		case name == "description":
			res.Description = content
		case name == "keywords":
			res.Keywords = append(res.Keywords, content)
		}
	})
	return res, nil

}

// Retrier

// Circuit breaker
// ===============
// CircuitBreaker struct represents simple circuit breaker
type CircuitBreaker struct {
	failMax      int           // Max failes to open breaker
	failCount    int           // fail counter
	resetTime    time.Duration // Time before breaker close
	lastFailTime time.Time     // Time to start resetTime
	mu           sync.Mutex    // Mutex
}

func NewCircuitBreaker(failMax int, resetTime time.Duration) *CircuitBreaker {
	if failMax <= 0 {
		failMax = 3
	}
	if resetTime <= 0 {
		resetTime = 1 * time.Second
	}
	return &CircuitBreaker{
		failMax:   failMax,
		resetTime: resetTime,
	}
}

// Execute check circuit braker state and execute callback func
func (cb *CircuitBreaker) Exexute(
	ctx context.Context,
	url string,
	operate func(context.Context, string) ([]byte, error),
) (
	[]byte, error) {
	cb.mu.Lock()
	// Check for circuit braker state
	if cb.failCount >= cb.failMax {
		if time.Since(cb.lastFailTime) > cb.resetTime {
			cb.failCount = 0
		} else {
			// If open returns error
			cb.mu.Unlock()
			return nil, fmt.Errorf("circuit braker: %w", ErrTooManyFailures)
		}
	}
	// If closd execute callback
	cb.mu.Unlock()
	res, err := operate(ctx, url)
	// If operate returns error renew fail counter and time of last fail
	if err != nil {
		cb.mu.Lock()
		cb.failCount++
		cb.lastFailTime = time.Now()
		cb.mu.Unlock()
	}
	// Returns result of callback
	return res, err
}

// Cache
