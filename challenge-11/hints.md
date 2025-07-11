# Hints for Challenge 11: Concurrent Web Content Aggregator

## Hint 1: Understanding the Core Structure
Start by implementing the basic ContentAggregator struct with the necessary fields:
```go
type ContentAggregator struct {
    fetcher         ContentFetcher
    processor       ContentProcessor
    workerCount     int
    requestLimiter  *rate.Limiter  // for rate limiting
    wg              sync.WaitGroup
    shutdown        chan struct{}
    shutdownOnce    sync.Once
}
```

## Hint 2: Rate Limiting Implementation
Use Go's `golang.org/x/time/rate` package for rate limiting:
```go
import "golang.org/x/time/rate"

// In constructor
requestLimiter := rate.NewLimiter(rate.Limit(requestsPerSecond), 1)

// Before making requests
err := requestLimiter.Wait(ctx)
if err != nil {
    return err // context cancelled or deadline exceeded
}
```

## Hint 3: Worker Pool Pattern
Create a worker pool that processes jobs from a channel:
```go
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
                    // Process URL here
                case <-ctx.Done():
                    return
                }
            }
        }()
    }
}
```

## Hint 4: Fan-Out, Fan-In Pattern
Distribute URLs to workers and collect results:
```go
func (ca *ContentAggregator) FetchAndProcess(ctx context.Context, urls []string) ([]ProcessedData, error) {
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
}
```

## Hint 5: Context Propagation and Error Handling
Always pass context down the call chain and handle cancellation:
```go
// In worker processing
content, err := ca.fetcher.Fetch(ctx, url)
if err != nil {
    select {
    case errors <- fmt.Errorf("fetch error for %s: %w", url, err):
    case <-ctx.Done():
    }
    return
}

processedData, err := ca.processor.Process(ctx, content)
if err != nil {
    select {
    case errors <- fmt.Errorf("process error for %s: %w", url, err):
    case <-ctx.Done():
    }
    return
}
```

## Hint 6: Graceful Shutdown
Implement proper cleanup in the shutdown method:
```go
func (ca *ContentAggregator) Shutdown() error {
    ca.shutdownOnce.Do(func() {
        close(ca.shutdown)
        ca.wg.Wait() // Wait for all workers to finish
    })
    return nil
}
```

## Hint 7: Result Collection Pattern
Use a separate goroutine to collect results and handle the done signal:
```go
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
case <-ctx.Done():
    return nil, ctx.Err()
}
```

## Key Concepts to Remember:
- **Context**: Always propagate context and check for cancellation
- **Channels**: Use buffered channels to avoid blocking
- **WaitGroup**: Coordinate goroutine completion
- **Rate Limiting**: Respect rate limits to avoid overwhelming servers
- **Error Handling**: Collect and return meaningful errors
- **Resource Cleanup**: Always close channels and wait for goroutines 