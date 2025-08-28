package challenge12

import (
	"context"
	"errors"
	"fmt"
	"os"
	"encoding/json"
)

type Reader interface {
	Read(ctx context.Context) ([]byte, error)
}

type Validator interface {
	Validate(data []byte) error
}

type Transformer interface {
	Transform(data []byte) ([]byte, error)
}

type Writer interface {
	Write(ctx context.Context, data []byte) error
}

type ValidationError struct {
	Field   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s error : %s", e.Field, e.Message, e.Err.Error())
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

type TransformError struct {
	Stage string
	Err   error
}

func (e *TransformError) Error() string {
	return fmt.Sprintf("transform error, stage: %s, %v", e.Stage, e.Err)
}

func (e *TransformError) Unwrap() error {
	return e.Err
}

type PipelineError struct {
	Stage string
	Err   error
}

func (e *PipelineError) Error() string {
	return fmt.Sprintf("pipeline error, stage: %s, %v", e.Stage, e.Err)
}

func (e *PipelineError) Unwrap() error {
	return e.Err
}

var (
	ErrInvalidFormat        = errors.New("invalid data format")
	ErrMissingField         = errors.New("required field missing")
	ErrProcessingFailed     = errors.New("processing failed")
	ErrTransformationFailed = errors.New("transform error")
	ErrDestinationFull      = errors.New("destination is full")
)

type Pipeline struct {
	Reader       Reader
	Validators   []Validator
	Transformers []Transformer
	Writer       Writer
}

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

func (p *Pipeline) Process(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		data, err := p.Reader.Read(ctx)
		if err != nil {
			return err
		}

		for _, v := range(p.Validators) {
			if err = v.Validate(data); err != nil {
				return err
			}
		}

		for _, t := range(p.Transformers) {
			if data, err = t.Transform(data); err != nil {
				return err
			}
			if data == nil {
				return ErrTransformationFailed
			}
		}

		if err = p.Writer.Write(ctx, nil); err != nil {
			return err
		}
		return nil
	}
}

//
// NOTE: unused method
//
// func (p *Pipeline) handleErrors(ctx context.Context, errs <-chan error) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case err := <-errs:
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// }

type FileReader struct {
	Filename string
}

func NewFileReader(filename string) *FileReader {
	return &FileReader{Filename: filename}
}

func (fr *FileReader) Read(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		data, err := os.ReadFile(fr.Filename)
		if err != nil {
			return nil, fmt.Errorf("file read error: %w", err)
		}
		return data, nil
	}
}

type JSONValidator struct{}

func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

func (jv *JSONValidator) Validate(data []byte) error {
	var tmp map[string]any
	if err := json.Unmarshal(data, &tmp); err != nil {
		return &ValidationError{
			Field: "",
			Message: err.Error(),
			Err: ErrInvalidFormat,
		}
	}
	return nil
}

type SchemaValidator struct {
	Schema []byte
}

func NewSchemaValidator(schema []byte) *SchemaValidator {
	return &SchemaValidator{Schema: schema}
}

func (sv *SchemaValidator) Validate(data []byte) error {
	var tmp map[string]any
	if err := json.Unmarshal(data, &tmp); err != nil {
		return &ValidationError{
			Field: "",
			Message: err.Error(),
			Err: ErrInvalidFormat,
		}
	}
	return nil
}

type FieldTransformer struct {
	FieldName     string
	TransformFunc func(string) string
}

func NewFieldTransformer(fieldName string, transformFunc func(string) string) *FieldTransformer {
	return &FieldTransformer{FieldName: fieldName, TransformFunc: transformFunc}
}

func (ft *FieldTransformer) Transform(data []byte) ([]byte, error) {
	var parsedData map[string]any
	if err := json.Unmarshal(data, &parsedData); err != nil {
		return nil, &TransformError{Stage: ft.FieldName, Err: ErrInvalidFormat}
	}

	if val, ok := parsedData[ft.FieldName]; ok {
		if strVal, ok := val.(string); ok {
			parsedData[ft.FieldName] = ft.TransformFunc(strVal)
		} else {
			return nil, &TransformError{
				Stage: ft.FieldName,
				Err:   ErrInvalidFormat,
			}
		}
	} else {
		return nil, &TransformError{
			Stage: ft.FieldName,
			Err:   ErrMissingField,
		}
	}

	result, err := json.Marshal(parsedData)
	if err != nil {
		return nil, &TransformError{Stage: ft.FieldName, Err: err}
	}
	return result, nil
}

type FileWriter struct {
	Filename string
}

func NewFileWriter(filename string) *FileWriter {
	return &FileWriter{Filename: filename}
}

func (fw *FileWriter) Write(ctx context.Context, data []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		info, err := os.Stat(fw.Filename)
		if err != nil {
			return err
		}
		if info.Size() > 0 {
			return ErrDestinationFull
		}

		err = os.WriteFile(fw.Filename, data, 0644)
		if err != nil {
			return err
		}
		return nil
	}
} 
