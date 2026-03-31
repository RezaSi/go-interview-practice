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
	"golang.org/x/net/html"
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
	ErrShutdown        = errors.New("aggregator is shut down")
	ErrTooManyFailures = errors.New("too many failures")
)

// HTTPStatusError reports a non-success HTTP response code.
type HTTPStatusError struct {
	StatusCode int
	URL        string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("unexpected status code %d for %s", e.StatusCode, e.URL)
}

// ProcessedData represents structured data extracted from raw content
type ProcessedData struct {
	Title       string
	Description string
	Keywords    []string
	Timestamp   time.Time
	Source      string
}

// MultiError preserves all errors produced during concurrent processing.
type MultiError struct {
	errs []error
}

func (m MultiError) Error() string {
	if len(m.errs) == 0 {
		return ""
	}
	if len(m.errs) == 1 {
		return m.errs[0].Error()
	}

	var b strings.Builder
	b.WriteString("multiple errors:")
	for _, err := range m.errs {
		if err == nil {
			continue
		}
		b.WriteString("\n- ")
		b.WriteString(err.Error())
	}

	return b.String()
}

func (m MultiError) Unwrap() []error {
	return m.errs

}
func (m MultiError) HasErrors() bool {
	return len(m.errs) > 0
}

// Aggregator - the main orchestrator
// =================================
// ContentAggregator manages the concurrent fetching and processing of content
type ContentAggregator struct {
	fetcher           ContentFetcher
	processor         ContentProcessor
	workerCount       int
	requestsPerSecond int
	limiter           *rate.Limiter
	activeRuns        sync.WaitGroup
	mu                sync.RWMutex
	shutdown          bool
	shutdownCh        chan struct{}
	closeOnce         sync.Once
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
		limiter:           rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
		shutdownCh:        make(chan struct{}),
	}
}

// Aggregator methods
// ==================
// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	if err := ca.validate(); err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return []ProcessedData{}, nil
	}

	ca.mu.RLock()
	if ca.shutdown {
		ca.mu.RUnlock()
		return nil, fmt.Errorf("aggregator error: %w", ErrShutdown)
	}
	ca.activeRuns.Add(1)
	ca.mu.RUnlock()
	defer ca.activeRuns.Done()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		select {
		case <-ca.shutdownCh:
			cancel()
		case <-runCtx.Done():
		}
	}()

	jobs := make(chan string)
	results := make(chan ProcessedData)
	errCh := make(chan error)

	// Write jobs (URLs) into jobs channel
	go func() {
		for _, url := range urls {
			select {
			case <-runCtx.Done():
				close(jobs)
				return
			case jobs <- url:
			}
		}
		close(jobs)
	}()

	ca.workerPool(runCtx, jobs, results, errCh)

	collected := make([]ProcessedData, 0, len(urls))
	var allErrs MultiError

	for {
		if results == nil && errCh == nil {
			if allErrs.HasErrors() {
				return collected, allErrs
			}
			return collected, nil
		}

		select {
		case res, ok := <-results:
			if !ok {
				results = nil
				continue
			}
			collected = append(collected, res)
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
				continue
			}
			if err != nil {
				allErrs.errs = append(allErrs.errs, err)
			}
		case <-runCtx.Done():
			if allErrs.HasErrors() {
				allErrs.errs = append(allErrs.errs, runCtx.Err())
				return collected, allErrs
			}
			return collected, runCtx.Err()
		}
	}
}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	if err := ca.validate(); err != nil {
		return err
	}

	ca.closeOnce.Do(func() {
		ca.mu.Lock()
		ca.shutdown = true
		close(ca.shutdownCh)
		ca.mu.Unlock()
	})
	ca.activeRuns.Wait()
	return nil
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
		go func() {
			defer wg.Done()
			for {
				var job string
				var ok bool
				select {
				case <-ctx.Done():
					return
				case job, ok = <-jobs:
					if !ok {
						return
					}
				}

				if ca.limiter != nil {
					if err := ca.limiter.Wait(ctx); err != nil {
						select {
						case <-ctx.Done():
							return
						case errors <- err:
						}
						return
					}
				}

				content, err := ca.fetcher.Fetch(ctx, job)
				if err != nil {
					select {
					case <-ctx.Done():
						return
					case errors <- err:
					}
					continue
				}
				processedContent, err := ca.processor.Process(ctx, content)
				if err != nil {
					select {
					case <-ctx.Done():
						return
					case errors <- err:
					}
					continue
				}
				if processedContent.Source == "" {
					processedContent.Source = job
				}
				if processedContent.Timestamp.IsZero() {
					processedContent.Timestamp = time.Now().UTC()
				}
				select {
				case <-ctx.Done():
					return
				case results <- processedContent:
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently
// The method is not used in the main flow but can be used for testing or as an alternative approach to workerPool
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	if ca == nil {
		return nil, []error{fmt.Errorf("aggregator error: %w", ErrNilReceiver)}
	}
	if len(urls) == 0 {
		return []ProcessedData{}, nil
	}

	resChan := make(chan ProcessedData, len(urls))
	errs := make(chan error, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		url := url
		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			content, err := ca.fetcher.Fetch(ctx, url)
			if err != nil {
				select {
				case <-ctx.Done():
				case errs <- err:
				}
				return
			}

			processedContent, err := ca.processor.Process(ctx, content)
			if err != nil {
				select {
				case <-ctx.Done():
				case errs <- err:
				}
				return
			}
			if processedContent.Source == "" {
				processedContent.Source = url
			}
			if processedContent.Timestamp.IsZero() {
				processedContent.Timestamp = time.Now().UTC()
			}
			select {
			case <-ctx.Done():
			case resChan <- processedContent:
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resChan)
		close(errs)
	}()

	processedData := make([]ProcessedData, 0, len(urls))
	for res := range resChan {
		processedData = append(processedData, res)
	}

	errList := make([]error, 0)
	for err := range errs {
		errList = append(errList, err)
	}

	if ctx.Err() != nil {
		errList = append(errList, ctx.Err())
	}

	return processedData, errList
}

func (ca *ContentAggregator) validate() error {
	if ca == nil {
		return fmt.Errorf("aggregator error: %w", ErrNilReceiver)
	}
	if ca.fetcher == nil || ca.processor == nil || ca.workerCount <= 0 || ca.requestsPerSecond <= 0 || ca.shutdownCh == nil {
		return fmt.Errorf("aggregator error: %w", ErrBadParam)
	}
	return nil
}

// Fetcher
// ========
// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
}

// Fetch retrieves content from a URL via HTTP not exceeding rate limit
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	// Validation input parameters
	if hf == nil {
		return nil, fmt.Errorf("fetcher error: %w", ErrNilReceiver)
	}
	if url == "" {
		return nil, fmt.Errorf("fetcher error: %w", ErrBadUrl)
	}
	client := hf.Client
	if client == nil {
		client = http.DefaultClient
	}

	// New request with context added
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fetcher error: %w: %v", ErrBadUrl, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}

	// Request body defer close
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &HTTPStatusError{
			StatusCode: resp.StatusCode,
			URL:        url,
		}
	}

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

var standardHTMLElements = map[string]struct{}{
	"a":          {},
	"abbr":       {},
	"address":    {},
	"area":       {},
	"article":    {},
	"aside":      {},
	"audio":      {},
	"b":          {},
	"base":       {},
	"bdi":        {},
	"bdo":        {},
	"blockquote": {},
	"body":       {},
	"br":         {},
	"button":     {},
	"canvas":     {},
	"caption":    {},
	"cite":       {},
	"code":       {},
	"col":        {},
	"colgroup":   {},
	"data":       {},
	"datalist":   {},
	"dd":         {},
	"del":        {},
	"details":    {},
	"dfn":        {},
	"dialog":     {},
	"div":        {},
	"dl":         {},
	"dt":         {},
	"em":         {},
	"embed":      {},
	"fieldset":   {},
	"figcaption": {},
	"figure":     {},
	"footer":     {},
	"form":       {},
	"h1":         {},
	"h2":         {},
	"h3":         {},
	"h4":         {},
	"h5":         {},
	"h6":         {},
	"head":       {},
	"header":     {},
	"hgroup":     {},
	"hr":         {},
	"html":       {},
	"i":          {},
	"iframe":     {},
	"img":        {},
	"input":      {},
	"ins":        {},
	"kbd":        {},
	"label":      {},
	"legend":     {},
	"li":         {},
	"link":       {},
	"main":       {},
	"map":        {},
	"mark":       {},
	"meta":       {},
	"meter":      {},
	"nav":        {},
	"noscript":   {},
	"object":     {},
	"ol":         {},
	"optgroup":   {},
	"option":     {},
	"output":     {},
	"p":          {},
	"picture":    {},
	"pre":        {},
	"progress":   {},
	"q":          {},
	"rp":         {},
	"rt":         {},
	"ruby":       {},
	"s":          {},
	"samp":       {},
	"script":     {},
	"search":     {},
	"section":    {},
	"select":     {},
	"slot":       {},
	"small":      {},
	"source":     {},
	"span":       {},
	"strong":     {},
	"style":      {},
	"sub":        {},
	"summary":    {},
	"sup":        {},
	"table":      {},
	"tbody":      {},
	"td":         {},
	"template":   {},
	"textarea":   {},
	"tfoot":      {},
	"th":         {},
	"thead":      {},
	"time":       {},
	"title":      {},
	"tr":         {},
	"track":      {},
	"u":          {},
	"ul":         {},
	"var":        {},
	"video":      {},
	"wbr":        {},
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
	if err := validateHTML(content); err != nil {
		return ProcessedData{}, fmt.Errorf("processor error: %w", err)
	}

	// Data processing
	var res ProcessedData
	r := strings.NewReader(string(content))
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return ProcessedData{}, err
	}
	title := doc.Find("title").Text()
	res.Title = title
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		content, _ := s.Attr("content")
		switch {
		case name == "description":
			res.Description = content
		case name == "keywords":
			for _, keyword := range strings.Split(content, ",") {
				keyword = strings.TrimSpace(keyword)
				if keyword != "" {
					res.Keywords = append(res.Keywords, keyword)

				}
			}
		}
	})
	return res, nil

}

func validateHTML(content []byte) error {
	if len(strings.TrimSpace(string(content))) == 0 {
		return ErrBadParam
	}

	if _, err := html.Parse(strings.NewReader(string(content))); err != nil {
		return err
	}

	tokenizer := html.NewTokenizer(strings.NewReader(string(content)))
	foundRecognizedElement := false

	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			err := tokenizer.Err()
			if errors.Is(err, io.EOF) {
				if !foundRecognizedElement {
					return ErrBadParam
				}
				return nil
			}
			return err
		case html.StartTagToken:
			token := tokenizer.Token()
			if isRecognizedHTMLElement(token.Data) {
				foundRecognizedElement = true
			}
		case html.SelfClosingTagToken:
			token := tokenizer.Token()
			if isRecognizedHTMLElement(token.Data) {
				foundRecognizedElement = true
			}
		}
	}
}

func isRecognizedHTMLElement(tag string) bool {
	tag = strings.ToLower(tag)
	if _, ok := standardHTMLElements[tag]; ok {
		return true
	}

	// Custom elements must contain a hyphen by HTML spec.
	return strings.Contains(tag, "-")
}

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
func (cb *CircuitBreaker) Execute(
	ctx context.Context,
	url string,
	operate func(context.Context, string) ([]byte, error),
) (
	[]byte, error) {
	if cb == nil {
		return nil, fmt.Errorf("circuit breaker error: %w", ErrNilReceiver)
	}
	if operate == nil {
		return nil, fmt.Errorf("circuit breader error: %w", ErrBadParam)
	}

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
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			cb.mu.Lock()
			cb.failCount++
			cb.lastFailTime = time.Now()
			cb.mu.Unlock()
		}
	} else {
		cb.mu.Lock()
		cb.failCount = 0
		cb.mu.Unlock()
	}
	// Returns result of callback
	return res, err
}

// Cache
