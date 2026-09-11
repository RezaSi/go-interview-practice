// Package challenge11 contains the solution for Challenge 11.
package challenge11

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
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
	fetcher           ContentFetcher
	processor         ContentProcessor
	workerCount       int
	requestsPerSecond int

	// shutdownCh is closed by Shutdown to signal that no further work should be
	// started and any in-flight batch should be cancelled. shutdownOnce makes
	// Shutdown safe to call multiple times (closing a closed channel panics).
	shutdownCh   chan struct{}
	shutdownOnce sync.Once
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	// Reject configurations that can't do useful work: no components, or
	// non-positive worker count / rate limit.
	if fetcher == nil || processor == nil || workerCount <= 0 || requestsPerSecond <= 0 {
		return nil
	}

	return &ContentAggregator{
		fetcher:           fetcher,
		processor:         processor,
		workerCount:       workerCount,
		requestsPerSecond: requestsPerSecond,
		shutdownCh:        make(chan struct{}),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	// Refuse new work once the aggregator has been shut down.
	select {
	case <-ca.shutdownCh:
		return nil, errors.New("content aggregator is shut down")
	default:
	}

	// Nothing to do: return an empty (non-nil) slice so callers can range safely.
	if len(urls) == 0 {
		return []ProcessedData{}, nil
	}

	// Delegate the concurrency to the fan-out / fan-in helper. It returns every
	// result it managed to produce plus every error it collected.
	results, errs := ca.fanOut(ctx, urls)

	// The challenge says never lose error information: if anything failed, join
	// all of the errors into one and drop the partial results.
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return results, nil
}

// Shutdown performs cleanup and ensures all resources are properly released.
// It is idempotent: calling it more than once is safe and returns nil.
func (ca *ContentAggregator) Shutdown() error {
	ca.shutdownOnce.Do(func() {
		close(ca.shutdownCh)
	})
	return nil
}

// workerPool runs ca.workerCount worker goroutines that pull URLs off `jobs`,
// fetch + process each one, and send the outcome to `results` or `errors`.
// It blocks until every worker has finished (jobs closed or ctx cancelled).
// It does NOT close `results` / `errors` — the caller owns those channels.
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errors chan<- error,
) {
	// Rate limiting: one shared ticker emits a token every 1/requestsPerSecond.
	// Every worker receives from the same channel, so the cap is global to the
	// pool, not per worker.
	rps := ca.requestsPerSecond
	if rps <= 0 {
		rps = 1
	}
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	defer ticker.Stop()

	var wg sync.WaitGroup
	for i := 0; i < ca.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				// Wait for a rate-limit token, but stay responsive to cancellation.
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}

				content, err := ca.fetcher.Fetch(ctx, url)
				if err != nil {
					errors <- fmt.Errorf("fetch %q: %w", url, err)
					continue
				}

				data, err := ca.processor.Process(ctx, content)
				if err != nil {
					errors <- fmt.Errorf("process %q: %w", url, err)
					continue
				}
				data.Source = url

				select {
				case <-ctx.Done():
					return
				case results <- data:
				}
			}
		}()
	}
	wg.Wait()
}

// fanOut wires up the pipeline: a feeder goroutine distributes `urls` onto a
// jobs channel (fan-out), the worker pool processes them concurrently, and this
// function collects every result and error (fan-in) before returning.
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	// Local cancellable context so the feeder unwinds cleanly on any early return.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Cancel this batch if Shutdown is called while it is running. The goroutine
	// also exits when the batch finishes (defer cancel above), so it never leaks.
	go func() {
		select {
		case <-ca.shutdownCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	jobs := make(chan string)
	results := make(chan ProcessedData)
	errCh := make(chan error, len(urls)) // buffered so a worker never blocks reporting

	// Fan-out: feed URLs in, then close `jobs` so idle workers leave their range loop.
	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- url:
			}
		}
	}()

	// Run the pool in the background; close the output channels once it's drained.
	go func() {
		ca.workerPool(ctx, jobs, results, errCh)
		close(results)
		close(errCh)
	}()

	// Fan-in: read both channels until each is closed.
	var (
		out  []ProcessedData
		errs []error
	)
	for results != nil || errCh != nil {
		select {
		case d, ok := <-results:
			if !ok {
				results = nil
				continue
			}
			out = append(out, d)
		case e, ok := <-errCh:
			if !ok {
				errCh = nil
				continue
			}
			errs = append(errs, e)
		}
	}
	return out, errs
}

// RateLimiter is the subset of *golang.org/x/time/rate.Limiter that HTTPFetcher
// needs. Declaring it as an interface keeps this package dependency-free: callers
// can plug in a real *rate.Limiter, and a nil value disables client-side limiting.
type RateLimiter interface {
	Wait(ctx context.Context) error
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP.
type HTTPFetcher struct {
	// Client is the HTTP client used for requests. If nil, http.DefaultClient
	// is used. Set Client.Timeout (or rely on the request context) for deadlines.
	Client *http.Client
	// Limiter, if non-nil, is consulted before every request so a single
	// HTTPFetcher can cap the rate at which it hits remote servers.
	Limiter RateLimiter
}

// Fetch retrieves content from a URL via HTTP, honouring the context for
// cancellation/timeout and returning an error for any non-2xx response.
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	client := hf.Client
	if client == nil {
		client = http.DefaultClient
	}

	// Block until the rate limiter allows another request (or ctx is done).
	if hf.Limiter != nil {
		if err := hf.Limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter wait for %q: %w", url, err)
		}
	}

	// NewRequestWithContext ties the request's lifetime to ctx: if ctx is
	// cancelled or times out, the in-flight request is aborted.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %q: %w", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %q: %w", url, err)
	}
	defer resp.Body.Close()

	// Treat anything outside 2xx as a failure so callers never process an
	// error page as if it were real content.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch %q: unexpected status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body for %q: %w", url, err)
	}
	return body, nil
}

// Regexes used to pull metadata out of HTML. This is a deliberately simple,
// dependency-free approach (no golang.org/x/net/html): good enough for well
// formed <head> metadata, not a general HTML parser.
var (
	titleRe       = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	metaTagRe     = regexp.MustCompile(`(?is)<meta\b[^>]*>`)
	metaNameRe    = regexp.MustCompile(`(?is)\bname\s*=\s*["']([^"']*)["']`)
	metaContentRe = regexp.MustCompile(`(?is)\bcontent\s*=\s*["']([^"']*)["']`)
)

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content.
type HTMLProcessor struct{}

// Process extracts structured data (title, description, keywords) from HTML.
// It returns an error for empty input or content with no <title>, so the
// aggregator won't treat a junk response as a successfully processed page.
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	if err := ctx.Err(); err != nil {
		return ProcessedData{}, err
	}

	html := strings.TrimSpace(string(content))
	if html == "" {
		return ProcessedData{}, errors.New("cannot process empty content")
	}

	titleMatch := titleRe.FindStringSubmatch(html)
	if titleMatch == nil {
		return ProcessedData{}, errors.New("no <title> element found: not valid HTML content")
	}

	data := ProcessedData{
		Title:     strings.TrimSpace(titleMatch[1]),
		Timestamp: time.Now(),
	}

	// Walk every <meta> tag and pick out the ones we care about.
	for _, tag := range metaTagRe.FindAllString(html, -1) {
		nameMatch := metaNameRe.FindStringSubmatch(tag)
		if nameMatch == nil {
			continue
		}
		contentMatch := metaContentRe.FindStringSubmatch(tag)
		if contentMatch == nil {
			continue
		}
		value := strings.TrimSpace(contentMatch[1])

		switch strings.ToLower(strings.TrimSpace(nameMatch[1])) {
		case "description":
			data.Description = value
		case "keywords":
			data.Keywords = splitKeywords(value)
		}
	}

	return data, nil
}

// splitKeywords turns a comma-separated keyword string into a trimmed slice,
// dropping empty entries. Returns nil when there are no usable keywords.
func splitKeywords(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
