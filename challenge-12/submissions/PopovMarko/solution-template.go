// Package challenge12 contains the solution for Challenge 12.
package challenge12

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
)

// Pipeline interfaces
// ===================
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

// Custom pipeline Errors
// ValidationError represents an error during data validation
// ==========================================================
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

// Error returns a string representation of the ValidationError
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("invalid: %s: %s: %s", e.Field, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("invalid: %s: %s", e.Field, e.Message)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return nil
}

// TransformError represents an error during data transformation
// =============================================================
type TransformError struct {
	Stage string
	Err   error
}

// Error returns a string representation of the TransformError
func (e *TransformError) Error() string {
	strErr := ""
	if e.Err != nil {
		return fmt.Sprintf("transformation stage: %s error: %s", e.Stage, strErr)
	}
	return fmt.Sprintf("transformation stage: %s", e.Stage)
}

// Unwrap returns the underlying error
func (e *TransformError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return nil
}

// PipelineError represents an error in the processing pipeline
// ============================================================
type PipelineError struct {
	Stage string
	Err   error
}

// Error returns a string representation of the PipelineError
func (e *PipelineError) Error() string {
	strErr := ""
	if e.Err != nil {
		return fmt.Sprintf("pipeline stage: %s error: %s", e.Stage, strErr)
	}
	return fmt.Sprintf("pipeline stage: %s", e.Stage)
}

// Unwrap returns the underlying error
func (e *PipelineError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return nil
}

// Sentinel errors for common error conditions
// ===========================================
var (
	ErrInvalidFormat    = errors.New("invalid data format")
	ErrMissingField     = errors.New("required field missing")
	ErrProcessingFailed = errors.New("processing failed")
	ErrDestinationFull  = errors.New("destination is full")
	ErrNilReceiver      = errors.New("nil receiver")
	ErrFieldType        = errors.New("type mismatch")
)

// Pipeline orchestrates the data processing flow
// ==============================================
type Pipeline struct {
	Reader       Reader
	Validators   []Validator
	Transformers []Transformer
	Writer       Writer
}

// NewPipeline creates a new processing pipeline with specified components
func NewPipeline(r Reader, v []Validator, t []Transformer, w Writer) *Pipeline {
	if r == nil {
		return nil
	}
	if w == nil {
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
func (p *Pipeline) Process(ctx context.Context) (err error) {
	// Clean up the resources before exit
	defer func() {
		if err != nil {
			p.CleanUp()
		}
	}()
	// Read the data
	data, err := p.Reader.Read(ctx)
	if err != nil {
		return &PipelineError{
			Stage: "read source",
			Err:   err,
		}
	}
	// Validate the data
	for _, v := range p.Validators {
		if err := v.Validate(data); err != nil {
			return &PipelineError{
				Stage: "validation",
				Err:   err,
			}
		}
	}
	// Transform the data

	for _, t := range p.Transformers {
		data, err = t.Transform(data)
		if err != nil {
			return &PipelineError{
				Stage: "transformation",
				Err:   err,
			}
		}
	}
	// Write the data
	if err := p.Writer.Write(ctx, data); err != nil {
		return &PipelineError{
			Stage: "write data",
			Err:   err,
		}
	}
	return nil
}

// TODO method implemented for future use in concurrent operation
// handleErrors consolidates errors from concurrent operations
func (p *Pipeline) handleErrors(ctx context.Context, errs <-chan error) error {
	select {
	case err, ok := <-errs:
		if ok {
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// FileReader implements the Reader interface for file sources
// ===========================================================
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
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		data, err := os.ReadFile(fr.Filename)
		if err != nil {
			return nil, fmt.Errorf("reader: can't read source: %w", err)
		}
		return data, nil
	}
}

// Validators implementation
// =========================
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
	}
	return ErrInvalidFormat
}

// SchemaValidator implements the Validator interface for schema validation
type SchemaValidator struct {
	Schema []byte
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator(schema []byte) *SchemaValidator {
	if !json.Valid(schema) {
		return nil
	}
	return &SchemaValidator{
		Schema: schema,
	}
}

// Validate validates data against a schema
func (sv *SchemaValidator) Validate(data []byte) error {
	if sv == nil {
		return fmt.Errorf("schema validator error: %w", ErrNilReceiver)
	}
	if !json.Valid(data) {
		return fmt.Errorf("invalid json data : %w", ErrInvalidFormat)
	}
	var mapData, mapSchema map[string]any
	if err := json.Unmarshal(sv.Schema, &mapSchema); err != nil {
		return fmt.Errorf("unmarshal schema error: %w", err)
	}
	if err := json.Unmarshal(data, &mapData); err != nil {
		return fmt.Errorf("unmarshal data error: %w", err)
	}
	for key, value := range mapSchema {
		if _, exist := mapData[key]; !exist {
			return &ValidationError{
				Field:   key,
				Message: "missing field",
				Err:     ErrMissingField,
			}
		}
		if reflect.TypeOf(value) != reflect.TypeOf(mapData[key]) {
			return &ValidationError{
				Field:   key,
				Message: "field type mismatch",
				Err:     ErrFieldType,
			}
		}
	}
	return nil
}

// Transformers implemetation
// ==========================
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
	if ft == nil {
		return nil, &TransformError{
			Stage: "transformer receiver",
			Err:   ErrNilReceiver,
		}
	}
	var mapData map[string]any
	if err := json.Unmarshal(data, &mapData); err != nil {
		return nil, &TransformError{
			Stage: "unmarshal data",
			Err:   err,
		}
	}
	if _, exist := mapData[ft.FieldName]; !exist {
		return nil, &TransformError{
			Stage: "map data",
			Err:   ErrMissingField,
		}
	}
	field, ok := mapData[ft.FieldName].(string)
	if !ok {
		return nil, &TransformError{
			Stage: "field type assertion",
			Err:   ErrFieldType,
		}
	}
	transformedField := ft.TransformFunc(field)
	mapData[ft.FieldName] = transformedField
	res, err := json.Marshal(mapData)
	if err != nil {
		return nil, &TransformError{
			Stage: "field transformer",
			Err:   err,
		}
	}
	return res, nil
}

// Writer implementation
// =====================
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
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := os.WriteFile(fw.Filename, data, 0644)
		if err != nil {
			return err
		}
		return nil
	}
}

// CleanUp clean the resources
func (p *Pipeline) CleanUp() {
	// TODO implement this method in case concurrent operation
}

