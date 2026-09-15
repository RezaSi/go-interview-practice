// Package challenge12 contains the solution for Challenge 12.
package challenge12

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
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
	if e.Err == nil {
		return fmt.Sprintf("validation error on the field %q: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error on the field %q: %s: %v", e.Field, e.Message, e.Err)
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
	if e.Err == nil {
		return fmt.Sprintf("transform error at stage %q", e.Stage)
	}
	return fmt.Sprintf("transform error at stage %q: %s", e.Stage, e.Err)
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
	if e.Err == nil {
		return fmt.Sprintf("pipeline error at stage \"%s\"", e.Stage)
	}
	return fmt.Sprintf("pipeline error at stage %q: %v", e.Stage, e.Err)
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
	if p == nil || p.Reader == nil || p.Writer == nil {
		return &PipelineError{Stage: "Process", Err: ErrProcessingFailed}

	}

	data, err := p.Reader.Read(ctx)
	if err != nil {
		return &PipelineError{Stage: "Read", Err: err}
	}
	for _, v := range p.Validators {
		select {
		case <-ctx.Done():
			return &PipelineError{Stage: "Validate", Err: ctx.Err()}
		default:
			err = v.Validate(data)
			if err != nil {
				return &PipelineError{Stage: "Validate", Err: err}
			}
		}
	}
	for _, t := range p.Transformers {
		select {
		case <-ctx.Done():
			return &PipelineError{Stage: "Transform", Err: ctx.Err()}
		default:
			data, err = t.Transform(data)
			if err != nil {
				return &PipelineError{Stage: "Transform", Err: err}
			}
		}
	}
	err = p.Writer.Write(ctx, data)
	if err != nil {
		return &PipelineError{Stage: "Write", Err: err}
	}
	return err
}

// handleErrors consolidates errors from concurrent operations
func (p *Pipeline) handleErrors(ctx context.Context, errs <-chan error) error {
	var errsSlice []error
	for {
		select {
		case err, ok := <-errs:
			if err != nil {
				errsSlice = append(errsSlice, err)
			}
			if !ok {
				return joinErrors(errsSlice)
			}
		case <-ctx.Done():
			if ctx.Err() != nil {
				errsSlice = append(errsSlice, ctx.Err())
			} else if len(errsSlice) == 0 {
				return nil
			}
			return joinErrors(errsSlice)
		}
	}
}

// FileReader implements the Reader interface for file sources
type FileReader struct {
	Filename string
}

// NewFileReader creates a new file reader
func NewFileReader(filename string) *FileReader {
	return &FileReader{Filename: filename}
}

// Read reads data from a file
func (fr *FileReader) Read(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return os.ReadFile(fr.Filename)
	}
}

// JSONValidator implements the Validator interface for JSON validation
type JSONValidator struct{}

// NewJSONValidator creates a new JSON validator
func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

// Validate validates JSON data
func (jv *JSONValidator) Validate(data []byte) error {
	if json.Valid(data) {
		return nil
	} else {
		return &ValidationError{
			Field:   "JSON",
			Message: "Json structure invalid",
			Err:     ErrInvalidFormat,
		}
	}
}

// SchemaValidator implements the Validator interface for schema validation
type SchemaValidator struct {
	Schema []byte
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator(schema []byte) *SchemaValidator {
	return &SchemaValidator{Schema: schema}
}

// Validate validates data against a schema
func (sv *SchemaValidator) Validate(data []byte) error {
	// TODO: Implement schema validation
	jsonMap := make(map[string]any)
	err := json.Unmarshal(data, &jsonMap)
	if err != nil {
		return &ValidationError{
			Field:   "Data",
			Message: "failed to parse data JSON",
			Err:     ErrInvalidFormat,
		}
	}
	schemaMap := make(map[string]any)
	err = json.Unmarshal(sv.Schema, &schemaMap)
	if err != nil {
		return &ValidationError{
			Field:   "Schema",
			Message: "failed to parse schema JSON",
			Err:     ErrInvalidFormat,
		}
	}
	for key := range schemaMap {
		if _, ok := jsonMap[key]; !ok {
			return &ValidationError{
				Field:   key,
				Message: fmt.Sprintf("missing required field: %s", key),
				Err:     ErrMissingField,
			}
		}
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
	return &FieldTransformer{
		FieldName:     fieldName,
		TransformFunc: transformFunc,
	}
}

// Transform transforms a specific field in the data
func (ft *FieldTransformer) Transform(data []byte) ([]byte, error) {
	if ft == nil || ft.TransformFunc == nil {
		return data, &TransformError{
			Stage: "Transform",
			Err:   fmt.Errorf("transform function is nil"),
		}
	}

	jsonMap := make(map[string]any)
	err := json.Unmarshal(data, &jsonMap)
	if err != nil {
		return data, &TransformError{
			Stage: "Transform",
			Err:   ErrInvalidFormat,
		}
	}
	if value, ok := jsonMap[ft.FieldName]; ok {
		// Type assertion: check if the underlying value is actually a string
		if strValue, isString := value.(string); isString {
			jsonMap[ft.FieldName] = ft.TransformFunc(strValue)
		} else {
			// If it's not a string, we can return a processing error
			return data, &TransformError{
				Stage: "Transform",
				Err:   fmt.Errorf("field %s is not a string", ft.FieldName),
			}
		}
	}

	// Marshal the modified map back into bytes
	resultData, err := json.Marshal(jsonMap)
	if err != nil {
		return data, &TransformError{
			Stage: "Transform",
			Err:   ErrProcessingFailed,
		}
	}
	return resultData, nil
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
	return &FileWriter{Filename: filename}
}

// Write writes data to a file
func (fw *FileWriter) Write(ctx context.Context, data []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := os.WriteFile(fw.Filename, data, 0644)
		return err
	}
}

type multiError struct {
	errors []error
}

// Join all errors into one
func joinErrors(errs []error) error {
	if len(errs) > 0 {
		return &multiError{errors: errs}
	}
	return nil
}

// Error implements the error interface
func (m *multiError) Error() string {
	var s []string
	for _, err := range m.errors {
		s = append(s, err.Error())
	}
	return strings.Join(s, "\n")
}

// Unwrap implements the interface to allow error unwrapping
func (m *multiError) Unwrap() []error {
	return m.errors
}
