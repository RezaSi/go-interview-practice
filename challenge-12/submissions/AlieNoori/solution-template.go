// Package challenge12 contains the solution for Challenge 12.
package challenge12

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	// Add any necessary imports here
)

// Reader defines an interface for data sources
type Reader interface {
	Read(ctx context.Context) ([]byte, error)
}

// Validator defines an interface for data validation
type Validator interface {
	Validate(data []byte) error
}

// Transformer defines an interface for data transformation
type Transformer interface {
	Transform(data []byte) ([]byte, error)
}

// Writer defines an interface for data destinations
type Writer interface {
	Write(ctx context.Context, data []byte) error
}

// ValidationError represents an error during data validation
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

// Error returns a string representation of the ValidationError
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation error: %s: %s: %s", e.Message, e.Field, e.Err)
	}
	return fmt.Sprintf("validation error: %s: %s", e.Message, e.Field)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// TransformError represents an error during data transformation
type TransformError struct {
	Stage string
	Err   error
}

// Error returns a string representation of the TransformError
func (e *TransformError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("transformation error: %s: %s", e.Stage, e.Err)
	}
	return fmt.Sprintf("transformation error: %s", e.Stage)
}

// Unwrap returns the underlying error
func (e *TransformError) Unwrap() error {
	return e.Err
}

// PipelineError represents an error in the processing pipeline
type PipelineError struct {
	Stage string
	Err   error
}

// Error returns a string representation of the PipelineError
func (e *PipelineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("pipeline error: %s: %s", e.Stage, e.Err)
	}
	return fmt.Sprintf("pipeline error: %s", e.Stage)
}

// Unwrap returns the underlying error
func (e *PipelineError) Unwrap() error {
	return e.Err
}

// Sentinel errors for common error conditions
var (
	ErrInvalidFormat    = errors.New("invalid data format")
	ErrMissingField     = errors.New("required field missing")
	ErrProcessingFailed = errors.New("processing failed")
	ErrDestinationFull  = errors.New("destination is full")
)

// Pipeline orchestrates the data processing flow
type Pipeline struct {
	Reader       Reader
	Validators   []Validator
	Transformers []Transformer
	Writer       Writer
}

// NewPipeline creates a new processing pipeline with specified components
func NewPipeline(r Reader, v []Validator, t []Transformer, w Writer) *Pipeline {
	if r == nil || w == nil {
		return nil
	}

	return &Pipeline{
		Reader:       r,
		Validators:   v,
		Transformers: t,
		Writer:       w,
	}
}

// Process runs the complete pipeline
func (p *Pipeline) Process(ctx context.Context) error {
	// Stage 1: Read
	// data, err := p.Reader.Read(ctx)
	// if err != nil {
	// 	return &PipelineError{Stage: "read", Err: err}
	// }

	errStream := make(chan error)

	readDataStream := p.read(ctx, errStream)

	// Stage 2: Validate
	// for i, validator := range p.Validators {
	// 	if err := validator.Validate(data); err != nil {
	// 		return &PipelineError{
	// 			Stage: fmt.Sprintf("validate_%d", i),
	// 			Err:   err,
	// 		}
	// 	}
	// }

	validateDataStream := p.validate(ctx, readDataStream, errStream)

	// Stage 3: Transform
	// for i, transformer := range p.Transformers {
	// 	data, err = transformer.Transform(data)
	// 	if err != nil {
	// 		return &PipelineError{
	// 			Stage: fmt.Sprintf("transform_%d", i),
	// 			Err:   err,
	// 		}
	// 	}
	// }

	transformDataStream := p.transform(ctx, validateDataStream, errStream)

	// Stage 4: Write
	// if err := p.Writer.Write(ctx, data); err != nil {
	// 	return &PipelineError{Stage: "write", Err: err}
	// }
	p.write(ctx, transformDataStream, errStream)

	err := p.handleErrors(ctx, errStream)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) read(ctx context.Context, errs chan<- error) <-chan []byte {
	result := make(chan []byte)

	go func() {
		defer close(result)

		select {
		case <-ctx.Done():
			return

		default:
			data, err := p.Reader.Read(ctx)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errs <- &PipelineError{Stage: "read", Err: err}:
					return
				}
			}

			select {
			case <-ctx.Done():
			case result <- data:
			}
		}
	}()

	return result
}

func (p *Pipeline) validate(ctx context.Context, dataStream <-chan []byte, errStream chan error) <-chan []byte {
	result := make(chan []byte)

	go func() {
		defer close(result)

		select {
		case <-ctx.Done():
			return

		case err := <-errStream:
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errStream <- err:
					return
				}
			}

		default:
			for data := range dataStream {
				for i, v := range p.Validators {
					err := v.Validate(data)
					if err != nil {
						select {
						case <-ctx.Done():
							return
						case errStream <- &PipelineError{Stage: fmt.Sprintf("validate_%d", i), Err: err}:
							return
						}
					}

					select {
					case <-ctx.Done():
						return
					case result <- data:
					}

				}
			}

		}
	}()

	return result
}

func (p *Pipeline) transform(ctx context.Context, dataStream <-chan []byte, errStream chan error) <-chan []byte {
	result := make(chan []byte)

	go func() {
		defer close(result)

		select {
		case <-ctx.Done():
			return
		case err := <-errStream:
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errStream <- err:
					return
				}
			}

		default:
			for data := range dataStream {
				for i, t := range p.Transformers {
					transformedData, err := t.Transform(data)
					if err != nil {
						select {
						case <-ctx.Done():
							return
						case errStream <- &PipelineError{Stage: fmt.Sprintf("transform_%d", i), Err: err}:
							return
						}
					}

					select {
					case <-ctx.Done():
						return
					case result <- transformedData:
					}
				}
			}

		}
	}()

	return result
}

func (p *Pipeline) write(ctx context.Context, dataStream <-chan []byte, errStream chan error) {
	go func() {
		defer close(errStream)

		select {
		case <-ctx.Done():
			return

		case err := <-errStream:
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errStream <- err:
					return
				}
			}

		case data := <-dataStream:
			err := p.Writer.Write(ctx, data)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errStream <- &PipelineError{Stage: "write", Err: err}:
					return
				}
			}
		}
	}()
}

// handleErrors consolidates errors from concurrent operations
func (p *Pipeline) handleErrors(ctx context.Context, errs <-chan error) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled: %w", ctx.Err())

	case err, ok := <-errs:
		if !ok {
			return err
		}
		if err != nil {
			return fmt.Errorf("pipeline error: %w", err)
		}
		return nil
	}
}

// FileReader implements the Reader interface for file sources
type FileReader struct {
	Filename string
}

// NewFileReader creates a new file reader
func NewFileReader(filename string) *FileReader {
	if filename == "" {
		return nil
	}
	return &FileReader{
		Filename: filename,
	}
}

// Read reads data from a file
func (fr *FileReader) Read(ctx context.Context) ([]byte, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return os.ReadFile(fr.Filename)
}

// JSONValidator implements the Validator interface for JSON validation
type JSONValidator struct{}

// NewJSONValidator creates a new JSON validator
func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

// Validate validates JSON data
func (jv *JSONValidator) Validate(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("data is empty: %w", ErrInvalidFormat)
	}
	if !json.Valid(data) {
		return fmt.Errorf("data is not valid JSON: %w", ErrInvalidFormat)
	}
	return nil
}

// SchemaValidator implements the Validator interface for schema validation
type SchemaValidator struct {
	Schema []byte
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator(schema []byte) *SchemaValidator {
	if len(schema) == 0 {
		return nil
	}
	return &SchemaValidator{
		Schema: schema,
	}
}

// Validate validates data against a schema
func (sv *SchemaValidator) Validate(data []byte) error {
	if !bytes.Equal(sv.Schema, data) {
		return fmt.Errorf("schema in not valid: %w", ErrInvalidFormat)
	}
	return nil
}

// FieldTransformer implements the Transformer interface for field transformations
type FieldTransformer struct {
	FieldName     string
	TransformFunc func(string) string
}

// NewFieldTransformer creates a new field transformer
func NewFieldTransformer(fieldName string, transformFunc func(string) string) *FieldTransformer {
	if fieldName == "" {
		return nil
	}
	if transformFunc == nil {
		return nil
	}
	return &FieldTransformer{
		FieldName:     fieldName,
		TransformFunc: transformFunc,
	}
}

// Transform transforms a specific field in the data
func (ft *FieldTransformer) Transform(data []byte) ([]byte, error) {
	return nil, nil
}

// FileWriter implements the Writer interface for file destinations
type FileWriter struct {
	Filename string
}

// NewFileWriter creates a new file writer
func NewFileWriter(filename string) *FileWriter {
	if filename == "" {
		return nil
	}
	return &FileWriter{
		Filename: filename,
	}
}

// Write writes data to a file
func (fw *FileWriter) Write(ctx context.Context, data []byte) error {
	return os.WriteFile(fw.Filename, data, 0o644)
}
