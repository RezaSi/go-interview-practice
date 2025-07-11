# Hints for Challenge 12: File Processing Pipeline with Advanced Error Handling

## Hint 1: Implementing Custom Error Types
Create error types that provide context and implement the error interface:
```go
type ValidationError struct {
    Field   string
    Message string
    Err     error
}

func (e *ValidationError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("validation failed for field '%s': %s: %v", e.Field, e.Message, e.Err)
    }
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
    return e.Err
}
```

## Hint 2: Error Wrapping and Context
Use fmt.Errorf with %w to wrap errors and preserve the error chain:
```go
func (p *Pipeline) Process(ctx context.Context) error {
    // Read stage
    data, err := p.Reader.Read(ctx)
    if err != nil {
        return &PipelineError{
            Stage: "read",
            Err:   fmt.Errorf("failed to read data: %w", err),
        }
    }
    
    // Validation stage
    for i, validator := range p.Validators {
        if err := validator.Validate(data); err != nil {
            return &PipelineError{
                Stage: fmt.Sprintf("validation_%d", i),
                Err:   fmt.Errorf("validation step %d failed: %w", i, err),
            }
        }
    }
}
```

## Hint 3: Sentinel Error Creation and Usage
Define package-level sentinel errors for common conditions:
```go
var (
    ErrInvalidFormat    = errors.New("invalid data format")
    ErrMissingField     = errors.New("required field missing")
    ErrProcessingFailed = errors.New("processing failed")
    ErrDestinationFull  = errors.New("destination is full")
)

// Usage in validator
func (v *JSONValidator) Validate(data []byte) error {
    if !json.Valid(data) {
        return fmt.Errorf("data is not valid JSON: %w", ErrInvalidFormat)
    }
    return nil
}
```

## Hint 4: Type-Based Error Handling with errors.As
Use errors.As to check for specific error types:
```go
func handlePipelineError(err error) {
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        log.Printf("Validation error in field '%s': %s", validationErr.Field, validationErr.Message)
        return
    }
    
    var transformErr *TransformError
    if errors.As(err, &transformErr) {
        log.Printf("Transform error at stage '%s': %v", transformErr.Stage, transformErr.Err)
        return
    }
    
    // Check for sentinel errors
    if errors.Is(err, ErrInvalidFormat) {
        log.Println("Invalid format detected")
        return
    }
}
```

## Hint 5: Pipeline Structure and Flow
Implement the pipeline to process data through all stages:
```go
type Pipeline struct {
    Reader       Reader
    Validators   []Validator
    Transformers []Transformer
    Writer       Writer
}

func (p *Pipeline) Process(ctx context.Context) error {
    // Stage 1: Read
    data, err := p.Reader.Read(ctx)
    if err != nil {
        return &PipelineError{Stage: "read", Err: err}
    }
    
    // Stage 2: Validate
    for i, validator := range p.Validators {
        if err := validator.Validate(data); err != nil {
            return &PipelineError{
                Stage: fmt.Sprintf("validate_%d", i),
                Err:   err,
            }
        }
    }
    
    // Stage 3: Transform
    for i, transformer := range p.Transformers {
        data, err = transformer.Transform(data)
        if err != nil {
            return &PipelineError{
                Stage: fmt.Sprintf("transform_%d", i),
                Err:   err,
            }
        }
    }
    
    // Stage 4: Write
    if err := p.Writer.Write(ctx, data); err != nil {
        return &PipelineError{Stage: "write", Err: err}
    }
    
    return nil
}
```

## Hint 6: Concurrent Error Handling
Handle errors from multiple goroutines properly:
```go
func (p *Pipeline) handleErrors(ctx context.Context, errs <-chan error) error {
    select {
    case err := <-errs:
        if err != nil {
            return fmt.Errorf("concurrent operation failed: %w", err)
        }
        return nil
    case <-ctx.Done():
        return fmt.Errorf("operation cancelled: %w", ctx.Err())
    }
}

// Example of concurrent processing with error collection
func (p *Pipeline) ProcessConcurrently(ctx context.Context, inputs [][]byte) error {
    errChan := make(chan error, len(inputs))
    
    for _, input := range inputs {
        go func(data []byte) {
            if err := p.processOne(ctx, data); err != nil {
                errChan <- err
                return
            }
            errChan <- nil
        }(input)
    }
    
    // Collect errors
    for i := 0; i < len(inputs); i++ {
        if err := <-errChan; err != nil {
            return err
        }
    }
    
    return nil
}
```

## Hint 7: Resource Cleanup with Error Handling
Implement proper cleanup even when errors occur:
```go
func (p *Pipeline) ProcessWithCleanup(ctx context.Context) (retErr error) {
    // Setup resources
    if setupErr := p.setup(); setupErr != nil {
        return fmt.Errorf("setup failed: %w", setupErr)
    }
    
    // Defer cleanup - this will run even if we return early due to errors
    defer func() {
        if cleanupErr := p.cleanup(); cleanupErr != nil {
            if retErr != nil {
                // Both processing and cleanup failed
                retErr = fmt.Errorf("processing failed: %w, cleanup also failed: %v", retErr, cleanupErr)
            } else {
                // Only cleanup failed
                retErr = fmt.Errorf("cleanup failed: %w", cleanupErr)
            }
        }
    }()
    
    // Main processing
    return p.Process(ctx)
}
```

## Key Error Handling Concepts:
- **Custom Error Types**: Provide context and implement Unwrap()
- **Error Wrapping**: Use %w to preserve error chains
- **Sentinel Errors**: Define package-level errors for common conditions
- **errors.Is**: Check if an error is or wraps a sentinel error
- **errors.As**: Extract specific error types from error chains
- **Context Preservation**: Always provide context about where errors occurred
- **Resource Cleanup**: Use defer to ensure cleanup happens even with errors 