package safeunmarshal

import (
	"errors"
	"testing"
)

// TestTo_BasicTypes tests basic types wrapped in objects
// Note: prepareJSONForUnmarshalling is designed for objects and arrays,
// not bare primitives
func TestTo_BasicTypes(t *testing.T) {
	type StringWrapper struct {
		Value string `json:"value"`
	}
	type IntWrapper struct {
		Value int `json:"value"`
	}
	type BoolWrapper struct {
		Value bool `json:"value"`
	}
	type FloatWrapper struct {
		Value float64 `json:"value"`
	}

	tests := []struct {
		name    string
		input   []byte
		want    interface{}
		wantErr bool
	}{
		{
			name:  "string in object",
			input: []byte(`{"value":"hello"}`),
			want:  StringWrapper{Value: "hello"},
		},
		{
			name:  "integer in object",
			input: []byte(`{"value":42}`),
			want:  IntWrapper{Value: 42},
		},
		{
			name:  "boolean in object",
			input: []byte(`{"value":true}`),
			want:  BoolWrapper{Value: true},
		},
		{
			name:  "float in object",
			input: []byte(`{"value":3.14}`),
			want:  FloatWrapper{Value: 3.14},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.want.(type) {
			case StringWrapper:
				got, err := To[StringWrapper](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != v {
					t.Errorf("To() = %v, want %v", got, v)
				}
			case IntWrapper:
				got, err := To[IntWrapper](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != v {
					t.Errorf("To() = %v, want %v", got, v)
				}
			case BoolWrapper:
				got, err := To[BoolWrapper](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != v {
					t.Errorf("To() = %v, want %v", got, v)
				}
			case FloatWrapper:
				got, err := To[FloatWrapper](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != v {
					t.Errorf("To() = %v, want %v", got, v)
				}
			}
		})
	}
}

func TestTo_Structs(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	tests := []struct {
		name    string
		input   []byte
		want    TestStruct
		wantErr bool
	}{
		{
			name:  "valid struct",
			input: []byte(`{"name":"John","age":30,"email":"john@example.com"}`),
			want: TestStruct{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
		},
		{
			name:  "struct with whitespace",
			input: []byte(`  {"name":"Jane","age":25,"email":"jane@example.com"}  `),
			want: TestStruct{
				Name:  "Jane",
				Age:   25,
				Email: "jane@example.com",
			},
		},
		{
			name:  "struct with newlines",
			input: []byte("{\n\"name\":\"Bob\",\n\"age\":35,\n\"email\":\"bob@example.com\"\n}"),
			want: TestStruct{
				Name:  "Bob",
				Age:   35,
				Email: "bob@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := To[TestStruct](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("To() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTo_Arrays(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantLen int
		wantErr bool
		errType error
	}{
		{
			name:    "valid int array",
			input:   []byte(`[1,2,3,4,5]`),
			wantLen: 5,
		},
		{
			name:    "valid string array",
			input:   []byte(`["a","b","c"]`),
			wantLen: 3,
		},
		{
			name:    "empty array",
			input:   []byte(`[]`),
			wantLen: 0,
		},
		{
			name:    "array with whitespace",
			input:   []byte(`  [1, 2, 3]  `),
			wantLen: 3,
		},
		{
			name:    "non-array for array type",
			input:   []byte(`{"key": "value"}`),
			wantErr: true,
			errType: ErrExpectedJSONArray,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid int array" || tt.name == "array with whitespace" {
				got, err := To[[]int](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("To() length = %v, want %v", len(got), tt.wantLen)
				}
			} else if tt.name == "valid string array" {
				got, err := To[[]string](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("To() length = %v, want %v", len(got), tt.wantLen)
				}
			} else if tt.name == "empty array" {
				got, err := To[[]int](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("To() length = %v, want %v", len(got), tt.wantLen)
				}
			} else if tt.name == "non-array for array type" {
				_, err := To[[]string](tt.input)
				if err == nil {
					t.Errorf("To() expected error but got none")
					return
				}
				if !errors.Is(err, tt.errType) {
					t.Errorf("To() error = %v, want error type %v", err, tt.errType)
				}
			}
		})
	}
}

func TestTo_Slices(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name    string
		input   []byte
		wantLen int
		wantErr bool
	}{
		{
			name:    "valid slice of structs",
			input:   []byte(`[{"name":"Alice","age":30},{"name":"Bob","age":25}]`),
			wantLen: 2,
		},
		{
			name:    "empty slice",
			input:   []byte(`[]`),
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := To[[]Person](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("To() length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestTo_EmptyInput(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "empty byte slice",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   []byte("   "),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := To[string](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTo_PrepareJSON(t *testing.T) {
	type TestStruct struct {
		Value string `json:"value"`
	}

	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:  "json with prefix text",
			input: []byte(`some text before {"value":"test"}`),
			want:  "test",
		},
		{
			name:  "clean json object",
			input: []byte(`{"value":"clean"}`),
			want:  "clean",
		},
		{
			name:  "json array",
			input: []byte(`[{"value":"arr1"},{"value":"arr2"}]`),
			want:  "arr1", // First element
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "json array" {
				got, err := To[[]TestStruct](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) > 0 && got[0].Value != tt.want {
					t.Errorf("To() = %v, want %v", got[0].Value, tt.want)
				}
			} else {
				got, err := To[TestStruct](tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got.Value != tt.want {
					t.Errorf("To() = %v, want %v", got.Value, tt.want)
				}
			}
		})
	}
}

func TestIsJSONArray(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  bool
	}{
		{
			name:  "valid array",
			input: []byte(`[1,2,3]`),
			want:  true,
		},
		{
			name:  "array with leading whitespace",
			input: []byte(`  [1,2,3]`),
			want:  true,
		},
		{
			name:  "object",
			input: []byte(`{"key":"value"}`),
			want:  false,
		},
		{
			name:  "empty string",
			input: []byte(``),
			want:  false,
		},
		{
			name:  "string",
			input: []byte(`"test"`),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isJSONArray(tt.input)
			if got != tt.want {
				t.Errorf("isJSONArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
