// Package challenge11 contains the solution for Challenge 11.
package challenge11

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/time/rate"
)

// ContentFetcher fetches raw bytes from a URL.
//
// Implementations must honor ctx for cancellation and deadlines, and return a
// non-nil error if the URL cannot be retrieved.
type ContentFetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// ContentProcessor turns raw fetched bytes into a structured ProcessedData.
//
// Implementations must honor ctx for cancellation and return an error for
// inputs they cannot interpret.
type ContentProcessor interface {
	Process(ctx context.Context, content []byte) (ProcessedData, error)
}

// ProcessedData is the structured result extracted from a single fetched
// source. Which fields are populated depends on the ContentProcessor
// implementation used.
type ProcessedData struct {
	Title       string
	Description string
	Keywords    []string
	Timestamp   time.Time
	Source      string
}

// ContentAggregator coordinates concurrent fetching and processing of URLs.
//
// It owns a worker pool sized by workerCount and rate-limits outbound requests
// at requestsPerSecond. The aggregator's lifetime is bounded by an internal
// shutdown context: once Shutdown is called, in-flight FetchAndProcess calls
// are cancelled and subsequent calls fail fast.
type ContentAggregator struct {
	fetcher           ContentFetcher
	processor         ContentProcessor
	workerCount       int
	requestsPerSecond int
	mu                sync.Mutex
	wg                sync.WaitGroup
	limiter           *rate.Limiter
	shutdownCtx       context.Context
	shutdownCtxCancel context.CancelFunc
	shutdownOnce      sync.Once
	done              chan struct{}
}

// NewContentAggregator constructs a ContentAggregator with the given fetcher,
// processor, and pool configuration.
//
// It returns nil if fetcher or processor is nil, or if workerCount or
// requestsPerSecond is non-positive.
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

	ctx, cancel := context.WithCancel(context.Background())

	return &ContentAggregator{
		fetcher:           fetcher,
		processor:         processor,
		workerCount:       workerCount,
		requestsPerSecond: requestsPerSecond,
		limiter:           rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
		shutdownCtx:       ctx,
		shutdownCtxCancel: cancel,
		done:              make(chan struct{}),
	}
}

// FetchAndProcess fetches and processes each URL concurrently and returns the
// successfully-processed results.
//
// The returned slice contains every ProcessedData the workers emitted. If any
// fetch or process step failed, FetchAndProcess additionally returns a non-nil
// aggregate error. Both ctx and the aggregator's shutdown context cancel
// in-flight work; after Shutdown has been called, this method returns
// immediately with an error.
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	ca.mu.Lock()
	if ca.shutdownCtx.Err() != nil {
		ca.mu.Unlock()
		return []ProcessedData{}, errors.New("ErrShutdown")
	}
	ca.wg.Add(1)
	ca.mu.Unlock()
	defer ca.wg.Done()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-ca.shutdownCtx.Done():
			cancel()
			return
		}
	}()

	jobs := make(chan string)
	results := make(chan ProcessedData)
	errs := make(chan error)

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

	ca.workerPool(ctx, jobs, results, errs)

	var (
		allResults []ProcessedData
		allErrors  []error
	)
	wg1 := &sync.WaitGroup{}
	wg1.Add(2)
	go func() {
		defer wg1.Done()
		for res := range results {
			allResults = append(allResults, res)
		}
	}()
	go func() {
		defer wg1.Done()
		for err := range errs {
			allErrors = append(allErrors, err)
		}
	}()

	wg1.Wait()
	if len(allErrors) > 0 {
		return allResults, errors.New("Data fetched or processed with errors")
	}
	return allResults, nil
}

// Shutdown cancels the aggregator's internal context to signal all in-flight
// FetchAndProcess calls to stop, then waits for them to drain.
//
// It is safe to call concurrently and repeatedly: the cancellation runs
// exactly once and every caller observes the same completion signal. Returns
// an error if in-flight work does not finish within 10 seconds.
func (ca *ContentAggregator) Shutdown() error {
	ca.mu.Lock()
	ca.shutdownOnce.Do(func() {
		ca.shutdownCtxCancel()
		ca.mu.Unlock()
		go func() {
			ca.wg.Wait()
			close(ca.done)
		}()
	})
	t := time.NewTimer(10 * time.Second)
	defer t.Stop()
	select {
	case <-ca.done:
		return nil
	case <-t.C:
		return errors.New("shutdown time out error")
	}
}

// workerPool spawns workerCount goroutines that consume URLs from jobs,
// throttle on the aggregator's rate limiter, fetch and process each URL, and
// emit results or errors on the corresponding channels.
//
// The function returns once all workers are started. A separate goroutine
// closes results and errs after every worker has exited, so callers can range
// over both channels until they drain.
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errs chan<- error,
) {
	wg2 := &sync.WaitGroup{}
	for i := 0; i < ca.workerCount; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for url := range jobs {
				if err := ca.limiter.Wait(ctx); err != nil {
					errs <- err
					if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
						return
					}
					continue
				}
				data, err := ca.fetcher.Fetch(ctx, url)
				if err != nil {
					errs <- err
					if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
						return
					}
					continue
				}
				processed, err := ca.processor.Process(ctx, data)
				if err != nil {
					errs <- err
					if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
						return
					}
					continue
				}

				if processed.Timestamp.IsZero() {
					processed.Timestamp = time.Now()
				}

				if processed.Source == "" {
					processed.Source = url
				}
				results <- processed
			}
		}()
	}
	go func() {
		wg2.Wait()
		close(results)
		close(errs)
	}()
}

// HTTPFetcher is a ContentFetcher backed by an *http.Client.
type HTTPFetcher struct {
	Client *http.Client
}

// Fetch issues an HTTP GET against url using the configured client and returns
// the response body bytes.
//
// It returns an error if the request cannot be constructed, the client fails,
// the response status is not 200 OK, or the body cannot be read. The response
// body is always closed before Fetch returns.
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if hf.Client == nil {
		hf.Client = http.DefaultClient
	}
	response, err := hf.Client.Do(r)
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []byte{}, errors.New("status code" + response.Status)
	}

	res, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// HTMLProcessor is a ContentProcessor that extracts a Title, Description, and
// Keywords from HTML content by walking its tokens.
type HTMLProcessor struct {
}

// Process tokenizes the HTML content and populates a ProcessedData with the
// document's <title> text and any <meta name="description"> /
// <meta name="keywords"> values it encounters.
//
// It returns an error for empty input, when the tokenizer fails before EOF,
// or when no <title> is present in the parsed document.
func (hp *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	if len(content) == 0 {
		return ProcessedData{}, errors.New("empty content")
	}

	z := html.NewTokenizer(bytes.NewReader(content))
	var (
		data    ProcessedData
		inTitle bool
	)

	for {
		switch z.Next() {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				if data.Title == "" {
					return ProcessedData{}, errors.New("no title found")
				}
				return data, nil
			}
			return ProcessedData{}, z.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			name, hasAttr := z.TagName()
			switch string(name) {
			case "title":
				inTitle = true
			case "meta":
				if !hasAttr {
					continue
				}
				var nameVal, contentVal string
				for {
					key, val, more := z.TagAttr()
					switch string(key) {
					case "name":
						nameVal = string(val)
					case "content":
						contentVal = string(val)
					}
					if !more {
						break
					}
				}
				switch nameVal {
				case "description":
					data.Description = contentVal
				case "keywords":
					data.Keywords = strings.Split(contentVal, ",")
				}
			}

		case html.TextToken:
			if inTitle {
				data.Title = strings.TrimSpace(string(z.Text()))
			}

		case html.EndTagToken:
			name, _ := z.TagName()
			if string(name) == "title" {
				inTitle = false
			}
		}
	}
}
