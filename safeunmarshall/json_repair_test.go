// Package safeunmarshall provides utilities for safely unmarshalling JSON data.
//
// Copyright (C) 2025 mhpenta (https://github.com/mhpenta)
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.
package safeunmarshall

// These tests are adapted from the Python JSON repair library tests
// to ensure our MIT implementation handles a wide range of edge cases.

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func Test_repairJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Invalid json with valid json",
			input:    `Here is a obj: {"name": "John", "age": 30}`,
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Invalid json with valid json array",
			input:    `Here is a obj: [{"name": "John", "age": 30}, {"name": "Marc", "age": 36}]`,
			expected: `[{"name": "John", "age": 30}, {"name": "Marc", "age": 36}]`,
			wantErr:  false,
		},
		{
			name:     "Already valid JSON",
			input:    `{"name": "John", "age": 30}`,
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "Code block markdown syntax",
			input:    "```json\n{\"name\": \"John\", \"age\": 30}\n```",
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Missing quotes around keys",
			input:    `{name: "John", age: 30}`,
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Single quotes instead of double quotes",
			input:    `{'name': 'John', 'age': 30}`,
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Trailing comma in object",
			input:    `{"name": "John", "age": 30,}`,
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Trailing comma in array",
			input:    `[1, 2, 3,]`,
			expected: `[1,2,3]`,
			wantErr:  false,
		},
		{
			name:     "Missing closing brace",
			input:    `{"name": "John", "age": 30`,
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Missing closing bracket",
			input:    `[1, 2, 3`,
			expected: `[1,2,3]`,
			wantErr:  false,
		},
		{
			name:     "Unquoted string values",
			input:    `{"name": John, "age": 30}`,
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Mixed single and double quotes",
			input:    `{"name": 'John', 'age': 30}`,
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Nested objects with errors",
			input:    `{"person": {name: "John", age: 30}}`,
			expected: `{"person":{"name":"John","age":30}}`,
			wantErr:  false,
		},
		{
			name:     "Array of objects with errors",
			input:    `[{name: "John", age: 30}, {name: "Jane", age: 25}]`,
			expected: `[{"name":"John","age":30},{"name":"Jane","age":25}]`,
			wantErr:  false,
		},
		{
			name:     "Complex structure with multiple issues",
			input:    `{users: [{name: 'John', hobbies: ['reading', 'sports',]}, {name: 'Jane', hobbies: ['painting', 'music',]}],}`,
			expected: `{"users":[{"name":"John","hobbies":["reading","sports"]},{"name":"Jane","hobbies":["painting","music"]}]}`,
			wantErr:  false,
		},
		{
			name:     "Boolean values",
			input:    `{active: true, verified: false}`,
			expected: `{"active":true,"verified":false}`,
			wantErr:  false,
		},
		{
			name:     "Null values",
			input:    `{name: "John", address: null}`,
			expected: `{"name":"John","address":null}`,
			wantErr:  false,
		},
		{
			name:     "Escaped quotes",
			input:    `{"message": "He said \"Hello\""}`,
			expected: `{"message":"He said \"Hello\""}`,
			wantErr:  false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name+"_"+strconv.Itoa(i+1), func(t *testing.T) {
			repaired, err := repairJSON(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("repairJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !jsonEqual(repaired, tt.expected) {
				t.Errorf("repairJSON() = %v, want %v", repaired, tt.expected)
			}
		})
	}
}

func Test_replaceQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No quotes to replace",
			input:    `"name":"John"`,
			expected: `"name":"John"`,
		},
		{
			name:     "Replace single quotes",
			input:    `'name':'John'`,
			expected: `"name":"John"`,
		},
		{
			name:     "Mixed quotes",
			input:    `'name':"John"`,
			expected: `"name":"John"`,
		},
		{
			name:     "Escaped quotes",
			input:    `'It\\'s a test'`,
			expected: `"It\\"s a test"`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name+"_"+strconv.Itoa(i+1), func(t *testing.T) {
			result := replaceQuotes(tt.input)
			if result != tt.expected {
				t.Errorf("replaceQuotes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func Test_fixUnquotedKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No unquoted keys",
			input:    `{"name": "John"}`,
			expected: `{"name": "John"}`,
		},
		{
			name:     "Single unquoted key",
			input:    `{name: "John"}`,
			expected: `{"name": "John"}`,
		},
		{
			name:     "Multiple unquoted keys",
			input:    `{name: "John", age: 30}`,
			expected: `{"name": "John", "age": 30}`,
		},
		{
			name:     "Nested unquoted keys",
			input:    `{person: {name: "John"}}`,
			expected: `{"person": {"name": "John"}}`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name+"_"+strconv.Itoa(i+1), func(t *testing.T) {
			result := fixUnquotedKeys(tt.input)
			if result != tt.expected {
				t.Errorf("fixUnquotedKeys() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func Test_removeTrailingCommas(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No trailing commas",
			input:    `{"items": [1, 2, 3]}`,
			expected: `{"items": [1, 2, 3]}`,
		},
		{
			name:     "Trailing comma in object",
			input:    `{"name": "John", "age": 30,}`,
			expected: `{"name": "John", "age": 30}`,
		},
		{
			name:     "Trailing comma in array",
			input:    `[1, 2, 3,]`,
			expected: `[1, 2, 3]`,
		},
		{
			name:     "Trailing commas in nested structures",
			input:    `{"people": [{"name": "John",}, {"name": "Jane",},]}`,
			expected: `{"people": [{"name": "John"}, {"name": "Jane"}]}`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name+"_"+strconv.Itoa(i+1), func(t *testing.T) {
			result := removeTrailingCommas(tt.input)
			if result != tt.expected {
				t.Errorf("removeTrailingCommas() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func Test_balanceBrackets(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Already balanced",
			input:    `{"name": "John"}`,
			expected: `{"name": "John"}`,
		},
		{
			name:     "Missing closing brace",
			input:    `{"name": "John"`,
			expected: `{"name": "John"}`,
		},
		{
			name:     "Missing closing bracket",
			input:    `[1, 2, 3`,
			expected: `[1, 2, 3]`,
		},
		{
			name:     "Multiple missing closures",
			input:    `{"people": [{"name": "John"`,
			expected: `{"people": [{"name": "John"}]}`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name+"_"+strconv.Itoa(i+1), func(t *testing.T) {
			result := balanceBrackets(tt.input)
			if result != tt.expected {
				t.Errorf("balanceBrackets() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test edge cases for bracket handling - adapted from Python tests
func Test_BracketsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unbalanced brackets",
			input:    "[{]",
			expected: "[{}]",
		},
		{
			name:     "Empty object with whitespace",
			input:    "   {  }   ",
			expected: "{}",
		},
		{
			name:     "Only opening bracket",
			input:    "[",
			expected: "[]",
		},
		{
			name:     "Only closing bracket",
			input:    "]",
			expected: `""`,
		},
		{
			name:     "Only opening brace",
			input:    "{",
			expected: "{}",
		},
		{
			name:     "Only closing brace",
			input:    "}",
			expected: `""`,
		},
		{
			name:     "Incomplete object start",
			input:    "{\"",
			expected: "{}",
		},
		{
			name:     "Incomplete array start",
			input:    "[\"",
			expected: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repairJSON(tt.input)
			if err != nil {
				t.Errorf("repairJSON() error = %v", err)
				return
			}

			if !jsonEqual(result, tt.expected) {
				t.Errorf("repairJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test array edge cases - adapted from Python tests
func Test_ArrayEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Trailing comma in array",
			input:    "[1, 2, 3,",
			expected: "[1,2,3]",
		},
		{
			name:     "Ellipsis at the end",
			input:    "[1, 2, 3, ...]",
			expected: "[1,2,3]",
		},
		{
			name:     "Unquoted strings in array",
			input:    "[\"a\" \"b\" \"c\" 1",
			expected: "[\"a\",\"b\",\"c\",1]",
		},
		{
			name:     "Nested incomplete structures",
			input:    "{\"key1\": {\"key2\": [1, 2, 3",
			expected: "{\"key1\":{\"key2\":[1,2,3]}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repairJSON(tt.input)
			if err != nil {
				t.Errorf("repairJSON() error = %v", err)
				return
			}

			if !jsonEqual(result, tt.expected) {
				t.Errorf("repairJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test general edge cases - adapted from Python tests
func Test_GeneralEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Only quote",
			input:    "\"",
			expected: `""`,
		},
		{
			name:     "Only newline",
			input:    "\n",
			expected: `""`,
		},
		{
			name:     "Only space",
			input:    " ",
			expected: `""`,
		},
		{
			name:     "Nested array with whitespace",
			input:    "[[1\n\n]",
			expected: "[[1]]",
		},
		{
			name:     "Plain string",
			input:    "string",
			expected: `""`,
		},
		{
			name:     "String before object",
			input:    "stringbeforeobject {}",
			expected: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repairJSON(tt.input)
			if err != nil {
				t.Errorf("repairJSON() error = %v", err)
				return
			}

			if !jsonEqual(result, tt.expected) {
				t.Errorf("repairJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// jsonEqual checks if two JSON strings are semantically equal
// by unmarshaling them to interface{} and comparing the results
func jsonEqual(json1, json2 string) bool {
	var obj1, obj2 interface{}

	if json1 == json2 {
		return true
	}

	err1 := json.Unmarshal([]byte(json1), &obj1)
	err2 := json.Unmarshal([]byte(json2), &obj2)

	if err1 != nil || err2 != nil {
		return false
	}

	return reflect.DeepEqual(obj1, obj2)
}

// TestTo tests the generic To function for unmarshalling JSON into various types
func TestTo(t *testing.T) {
	// Define test struct types
	type SimpleStruct struct {
		Name  string  `json:"name"`
		Age   int     `json:"age"`
		Score float64 `json:"score"`
	}

	type NestedStruct struct {
		Title  string       `json:"title"`
		Person SimpleStruct `json:"person"`
		Active bool         `json:"active"`
	}

	type ArrayContainer struct {
		Items []SimpleStruct `json:"items"`
	}

	// Test cases
	tests := []struct {
		name           string
		input          string
		expectedOutput interface{}
		expectErr      bool
		expectedErrIs  error
	}{
		{
			name:  "Valid simple struct",
			input: `{"name": "John", "age": 30, "score": 95.5}`,
			expectedOutput: SimpleStruct{
				Name:  "John",
				Age:   30,
				Score: 95.5,
			},
		},
		{
			name:  "Valid nested struct",
			input: `{"title": "Profile", "person": {"name": "Jane", "age": 25, "score": 98.2}, "active": true}`,
			expectedOutput: NestedStruct{
				Title: "Profile",
				Person: SimpleStruct{
					Name:  "Jane",
					Age:   25,
					Score: 98.2,
				},
				Active: true,
			},
		},
		{
			name:  "Valid array",
			input: `{"items": [{"name": "Alice", "age": 35, "score": 92.7}, {"name": "Bob", "age": 40, "score": 88.3}]}`,
			expectedOutput: ArrayContainer{
				Items: []SimpleStruct{
					{Name: "Alice", Age: 35, Score: 92.7},
					{Name: "Bob", Age: 40, Score: 88.3},
				},
			},
		},
		{
			name:  "Single-quoted keys and values",
			input: `{'name': 'John', 'age': 30, 'score': 95.5}`,
			expectedOutput: SimpleStruct{
				Name:  "John",
				Age:   30,
				Score: 95.5,
			},
		},
		{
			name:  "Unquoted keys",
			input: `{name: "John", age: 30, score: 95.5}`,
			expectedOutput: SimpleStruct{
				Name:  "John",
				Age:   30,
				Score: 95.5,
			},
		},
		{
			name:  "Trailing commas",
			input: `{"name": "John", "age": 30, "score": 95.5,}`,
			expectedOutput: SimpleStruct{
				Name:  "John",
				Age:   30,
				Score: 95.5,
			},
		},
		{
			name:  "JSON with code block markdown",
			input: "```json\n{\"name\": \"John\", \"age\": 30, \"score\": 95.5}\n```",
			expectedOutput: SimpleStruct{
				Name:  "John",
				Age:   30,
				Score: 95.5,
			},
		},
		{
			name:      "Empty input",
			input:     "",
			expectErr: true,
		},
		{
			name:  "Invalid JSON structure",
			input: `{"name": "John", "age": 30, "score": 95.5}`,
			expectedOutput: SimpleStruct{
				Name:  "John",
				Age:   30,
				Score: 95.5,
			},
			expectErr: false,
		},
		{
			name:          "Non-array when expecting array",
			input:         `{"name": "John"}`,
			expectErr:     true,
			expectedErrIs: ErrExpectedJSONArray,
			expectedOutput: []SimpleStruct{
				{Name: "John"},
			},
		},
		{
			name:      "Invalid field type",
			input:     `{"name": "John", "age": "thirty", "score": 95.5}`,
			expectErr: true,
		},
		{
			name:  "Repair array with mixed quotes",
			input: `{"items": [{'name': "Alice", 'age': 35, 'score': 92.7}, {'name': "Bob", 'age': 40, 'score': 88.3}]}`,
			expectedOutput: ArrayContainer{
				Items: []SimpleStruct{
					{Name: "Alice", Age: 35, Score: 92.7},
					{Name: "Bob", Age: 40, Score: 88.3},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch expectedOutput := tt.expectedOutput.(type) {
			case SimpleStruct:
				result, err := To[SimpleStruct]([]byte(tt.input))
				if (err != nil) != tt.expectErr {
					t.Errorf("To() error = %v, expectErr %v", err, tt.expectErr)
					return
				}
				if tt.expectedErrIs != nil && !errors.Is(err, tt.expectedErrIs) {
					t.Errorf("Expected error to be %v, got %v", tt.expectedErrIs, err)
				}
				if err == nil && !reflect.DeepEqual(result, expectedOutput) {
					t.Errorf("To() = %v, want %v", result, expectedOutput)
				}
			case NestedStruct:
				result, err := To[NestedStruct]([]byte(tt.input))
				if (err != nil) != tt.expectErr {
					t.Errorf("To() error = %v, expectErr %v", err, tt.expectErr)
					return
				}
				if tt.expectedErrIs != nil && !errors.Is(err, tt.expectedErrIs) {
					t.Errorf("Expected error to be %v, got %v", tt.expectedErrIs, err)
				}
				if err == nil && !reflect.DeepEqual(result, expectedOutput) {
					t.Errorf("To() = %v, want %v", result, expectedOutput)
				}
			case ArrayContainer:
				result, err := To[ArrayContainer]([]byte(tt.input))
				if (err != nil) != tt.expectErr {
					t.Errorf("To() error = %v, expectErr %v", err, tt.expectErr)
					return
				}
				if tt.expectedErrIs != nil && !errors.Is(err, tt.expectedErrIs) {
					t.Errorf("Expected error to be %v, got %v", tt.expectedErrIs, err)
				}
				if err == nil && !reflect.DeepEqual(result, expectedOutput) {
					t.Errorf("To() = %v, want %v", result, expectedOutput)
				}
			case []SimpleStruct:
				result, err := To[[]SimpleStruct]([]byte(tt.input))
				if (err != nil) != tt.expectErr {
					t.Errorf("To() error = %v, expectErr %v", err, tt.expectErr)
					return
				}
				if tt.expectedErrIs != nil && !errors.Is(err, tt.expectedErrIs) {
					t.Errorf("Expected error to be %v, got %v", tt.expectedErrIs, err)
				}
				if err == nil && !reflect.DeepEqual(result, expectedOutput) {
					t.Errorf("To() = %v, want %v", result, expectedOutput)
				}
			}
		})
	}
}
