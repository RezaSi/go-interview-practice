package challenge12

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

// MockReader implements the Reader interface for testing
type MockReader struct {
	data  []byte
	err   error
	reads int
}

func NewMockReader(data []byte, err error) *MockReader {
	return &MockReader{
		data:  data,
		err:   err,
		reads: 0,
	}
}

func (mr *MockReader) Read(ctx context.Context) ([]byte, error) {
	mr.reads++
	return mr.data, mr.err
}

func (mr *MockReader) GetReadCount() int {
	return mr.reads
}

// MockValidator implements the Validator interface for testing
type MockValidator struct {
	validData [][]byte
	err       error
	validates int
}

func NewMockValidator(validData [][]byte, err error) *MockValidator {
	return &MockValidator{
		validData: validData,
		err:       err,
		validates: 0,
	}
}

func (mv *MockValidator) Validate(data []byte) error {
	mv.validates++

	// Check if data matches any valid data
	for _, valid := range mv.validData {
		if reflect.DeepEqual(data, valid) {
			return nil
		}
	}

	return mv.err
}

func (mv *MockValidator) GetValidateCount() int {
	return mv.validates
}

// MockTransformer implements the Transformer interface for testing
type MockTransformer struct {
	input      []byte
	output     []byte
	err        error
	transforms int
}

func NewMockTransformer(input, output []byte, err error) *MockTransformer {
	return &MockTransformer{
		input:      input,
		output:     output,
		err:        err,
		transforms: 0,
	}
}

var ErrUnexpectedTransform = errors.New("unexpected transform")

func (mt *MockTransformer) Transform(data []byte) ([]byte, error) {
	mt.transforms++

	// If we have predefined error just return it - imitate error occurring during transformation
	if mt.err != nil {
		return nil, mt.err
	}

	if reflect.DeepEqual(data, mt.input) {
		return mt.output, nil
	}

	return nil, ErrUnexpectedTransform
}

func (mt *MockTransformer) GetTransformCount() int {
	return mt.transforms
}

// MockWriter implements the Writer interface for testing
type MockWriter struct {
	expectedData []byte
	err          error
	writes       int
}

func NewMockWriter(expectedData []byte, err error) *MockWriter {
	return &MockWriter{
		expectedData: expectedData,
		err:          err,
		writes:       0,
	}
}

var ErrUnexpectedWrite = errors.New("unexpected write")

func (mw *MockWriter) Write(_ context.Context, data []byte) error {
	mw.writes++

	// If we have predefined error just return it - imitate error occurring during writing data to the destination
	if mw.err != nil {
		return mw.err
	}

	if reflect.DeepEqual(data, mw.expectedData) {
		return nil
	}

	return ErrUnexpectedWrite
}

func (mw *MockWriter) GetWriteCount() int {
	return mw.writes
}

// TestValidationError tests the ValidationError implementation
func TestValidationError(t *testing.T) {
	underlying := errors.New("underlying error")
	valErr := &ValidationError{
		Field:   "username",
		Message: "must be at least 3 characters",
		Err:     underlying,
	}

	// Test Error() method
	errMsg := valErr.Error()
	if !strings.Contains(errMsg, "username") || !strings.Contains(errMsg, "must be at least 3 characters") {
		t.Errorf("Error message should contain field and message, got: %s", errMsg)
	}

	// Test Unwrap() method
	unwrapped := valErr.Unwrap()
	if unwrapped != underlying {
		t.Errorf("Unwrap() should return the underlying error")
	}
}

// TestTransformError tests the TransformError implementation
func TestTransformError(t *testing.T) {
	underlying := errors.New("underlying error")
	transErr := &TransformError{
		Stage: "normalization",
		Err:   underlying,
	}

	// Test Error() method
	errMsg := transErr.Error()
	if !strings.Contains(errMsg, "normalization") {
		t.Errorf("Error message should contain stage, got: %s", errMsg)
	}

	// Test Unwrap() method
	unwrapped := transErr.Unwrap()
	if unwrapped != underlying {
		t.Errorf("Unwrap() should return the underlying error")
	}
}

// TestPipelineError tests the PipelineError implementation
func TestPipelineError(t *testing.T) {
	underlying := errors.New("underlying error")
	pipeErr := &PipelineError{
		Stage: "validation",
		Err:   underlying,
	}

	// Test Error() method
	errMsg := pipeErr.Error()
	if !strings.Contains(errMsg, "validation") {
		t.Errorf("Error message should contain stage, got: %s", errMsg)
	}

	// Test Unwrap() method
	unwrapped := pipeErr.Unwrap()
	if unwrapped != underlying {
		t.Errorf("Unwrap() should return the underlying error")
	}
}

// TestNewPipeline tests the NewPipeline constructor
func TestNewPipeline(t *testing.T) {
	reader := NewMockReader([]byte("test data"), nil)
	validator := NewMockValidator([][]byte{[]byte("test data")}, nil)
	transformer := NewMockTransformer([]byte("test data"), []byte("transformed data"), nil)
	writer := NewMockWriter([]byte("transformed data"), nil)

	tests := []struct {
		name         string
		reader       Reader
		validators   []Validator
		transformers []Transformer
		writer       Writer
		expectNil    bool
	}{
		{
			name:         "Valid pipeline",
			reader:       reader,
			validators:   []Validator{validator},
			transformers: []Transformer{transformer},
			writer:       writer,
			expectNil:    false,
		},
		{
			name:         "No reader",
			reader:       nil,
			validators:   []Validator{validator},
			transformers: []Transformer{transformer},
			writer:       writer,
			expectNil:    true,
		},
		{
			name:         "No validators",
			reader:       reader,
			validators:   nil,
			transformers: []Transformer{transformer},
			writer:       writer,
			expectNil:    false, // Validators are optional
		},
		{
			name:         "Empty validators",
			reader:       reader,
			validators:   []Validator{},
			transformers: []Transformer{transformer},
			writer:       writer,
			expectNil:    false, // Empty validators are optional
		},
		{
			name:         "No transformers",
			reader:       reader,
			validators:   []Validator{validator},
			transformers: nil,
			writer:       writer,
			expectNil:    false, // Transformers are optional
		},
		{
			name:         "Empty transformers",
			reader:       reader,
			validators:   []Validator{validator},
			transformers: []Transformer{},
			writer:       writer,
			expectNil:    false, // Empty transformers are optional
		},
		{
			name:         "No writer",
			reader:       reader,
			validators:   []Validator{validator},
			transformers: []Transformer{transformer},
			writer:       nil,
			expectNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipeline := NewPipeline(tt.reader, tt.validators, tt.transformers, tt.writer)

			if tt.expectNil && pipeline != nil {
				t.Errorf("Expected nil pipeline, but got non-nil")
			}

			if !tt.expectNil && pipeline == nil {
				t.Errorf("Expected non-nil pipeline, but got nil")
			}

			if pipeline != nil {
				if pipeline.Reader != tt.reader {
					t.Errorf("Pipeline reader not set correctly")
				}

				if !reflect.DeepEqual(pipeline.Validators, tt.validators) {
					t.Errorf("Pipeline validators not set correctly")
				}

				if !reflect.DeepEqual(pipeline.Transformers, tt.transformers) {
					t.Errorf("Pipeline transformers not set correctly")
				}

				if pipeline.Writer != tt.writer {
					t.Errorf("Pipeline writer not set correctly")
				}
			}
		})
	}
}

// TestProcess tests the complete pipeline processing
func TestProcess(t *testing.T) {
	ctx := context.Background()

	// Success case
	t.Run("Success path", func(t *testing.T) {
		inputData := []byte(`{"name":"test","value":123}`)
		validatedData := inputData
		transformedData := []byte(`{"name":"TEST","value":123}`)

		reader := NewMockReader(inputData, nil)
		validator := NewMockValidator([][]byte{inputData}, nil)
		transformer := NewMockTransformer(validatedData, transformedData, nil)
		writer := NewMockWriter(transformedData, nil)

		pipeline := NewPipeline(reader, []Validator{validator}, []Transformer{transformer}, writer)
		if pipeline == nil {
			t.Fatal("Failed to create pipeline")
		}

		err := pipeline.Process(ctx)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		if reader.GetReadCount() != 1 {
			t.Errorf("Expected 1 read, got %d", reader.GetReadCount())
		}

		if validator.GetValidateCount() != 1 {
			t.Errorf("Expected 1 validation, got %d", validator.GetValidateCount())
		}

		if transformer.GetTransformCount() != 1 {
			t.Errorf("Expected 1 transformation, got %d", transformer.GetTransformCount())
		}

		if writer.GetWriteCount() != 1 {
			t.Errorf("Expected 1 write, got %d", writer.GetWriteCount())
		}
	})

	// Error cases
	t.Run("Reader error", func(t *testing.T) {
		reader := NewMockReader(nil, errors.New("read error"))
		validator := NewMockValidator([][]byte{}, nil)
		transformer := NewMockTransformer(nil, nil, nil)
		writer := NewMockWriter(nil, nil)

		pipeline := NewPipeline(reader, []Validator{validator}, []Transformer{transformer}, writer)
		if pipeline == nil {
			t.Fatal("Failed to create pipeline")
		}

		err := pipeline.Process(ctx)
		if err == nil {
			t.Errorf("Expected error, but got none")
		}

		if !strings.Contains(err.Error(), "read error") {
			t.Errorf("Error should contain read error message, got: %v", err)
		}
	})

	t.Run("Validation error", func(t *testing.T) {
		inputData := []byte(`{"name":"","value":123}`)

		reader := NewMockReader(inputData, nil)
		validator := NewMockValidator([][]byte{}, &ValidationError{Field: "name", Message: "required", Err: errors.New("validation error")})
		transformer := NewMockTransformer(nil, nil, nil)
		writer := NewMockWriter(nil, nil)

		pipeline := NewPipeline(reader, []Validator{validator}, []Transformer{transformer}, writer)
		if pipeline == nil {
			t.Fatal("Failed to create pipeline")
		}

		err := pipeline.Process(ctx)
		if err == nil {
			t.Errorf("Expected error, but got none")
		}

		if !strings.Contains(err.Error(), "validation") {
			t.Errorf("Error should contain validation info, got: %v", err)
		}
	})

	t.Run("Transformation error", func(t *testing.T) {
		inputData := []byte(`{"name":"test","value":123}`)

		reader := NewMockReader(inputData, nil)
		validator := NewMockValidator([][]byte{inputData}, nil)
		transformer := NewMockTransformer(inputData, nil, &TransformError{Stage: "uppercase", Err: errors.New("transform error")})
		writer := NewMockWriter(nil, nil)

		pipeline := NewPipeline(reader, []Validator{validator}, []Transformer{transformer}, writer)
		if pipeline == nil {
			t.Fatal("Failed to create pipeline")
		}

		err := pipeline.Process(ctx)
		if err == nil {
			t.Errorf("Expected error, but got none")
		}

		if !strings.Contains(err.Error(), "transform") {
			t.Errorf("Error should contain transform info, got: %v", err)
		}
	})

	t.Run("Writer error", func(t *testing.T) {
		inputData := []byte(`{"name":"test","value":123}`)
		transformedData := []byte(`{"name":"TEST","value":123}`)

		reader := NewMockReader(inputData, nil)
		validator := NewMockValidator([][]byte{inputData}, nil)
		transformer := NewMockTransformer(inputData, transformedData, nil)
		writer := NewMockWriter(transformedData, errors.New("write error"))

		pipeline := NewPipeline(reader, []Validator{validator}, []Transformer{transformer}, writer)
		if pipeline == nil {
			t.Fatal("Failed to create pipeline")
		}

		err := pipeline.Process(ctx)
		if err == nil {
			t.Errorf("Expected error, but got none")
		}

		if !strings.Contains(err.Error(), "write") {
			t.Errorf("Error should contain write info, got: %v", err)
		}
	})
}

// TestFileReader tests the FileReader implementation
func TestFileReader(t *testing.T) {
	// Create a temporary test file
	content := []byte("test file content")
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	t.Run("Valid file", func(t *testing.T) {
		reader := NewFileReader(tmpfile.Name())
		if reader == nil {
			t.Fatal("Failed to create FileReader")
		}

		data, err := reader.Read(context.Background())
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		if !reflect.DeepEqual(data, content) {
			t.Errorf("Expected %s, got %s", content, data)
		}
	})

	t.Run("Non-existent file", func(t *testing.T) {
		reader := NewFileReader("non-existent-file.txt")
		if reader == nil {
			t.Fatal("Failed to create FileReader")
		}

		_, err := reader.Read(context.Background())
		if err == nil {
			t.Errorf("Expected error for non-existent file, but got none")
		}
	})

	t.Run("Context cancellation", func(t *testing.T) {
		reader := NewFileReader(tmpfile.Name())
		if reader == nil {
			t.Fatal("Failed to create FileReader")
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately

		_, err := reader.Read(ctx)
		if err == nil || err != context.Canceled {
			t.Errorf("Expected context canceled error, but got: %v", err)
		}
	})
}

// TestJSONValidator tests the JSONValidator implementation
func TestJSONValidator(t *testing.T) {
	validator := NewJSONValidator()
	if validator == nil {
		t.Fatal("Failed to create JSONValidator")
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "Valid JSON",
			data:    []byte(`{"name":"test","value":123}`),
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			data:    []byte(`{"name":"test","value":123`), // Missing closing brace
			wantErr: true,
		},
		{
			name:    "Empty JSON",
			data:    []byte(``),
			wantErr: true,
		},
		{
			name:    "Non-JSON data",
			data:    []byte(`Hello, world!`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.data)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
