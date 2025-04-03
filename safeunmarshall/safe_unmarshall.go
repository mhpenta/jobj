// Package safeunmarshall provides utilities for safely unmarshalling JSON data.
//
// Copyright (C) 2025 mhpenta (https://github.com/mhpenta)
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.
package safeunmarshall

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

// ErrExpectedJSONArray is returned when the response being unmarshalled is not a JSON array in our unmarshaller
// utilities
var ErrExpectedJSONArray = errors.New("expected JSON array for array type")

// ToAny attempts to unmarshal a JSON string into a value of type T.
//
// This function is more flexible than SafeUnmarshallToStructPtr as it can handle
// various types, including arrays and slices. It first attempts to unmarshal the
// provided JSON string directly. If that fails, it tries to repair the JSON before
// attempting to unmarshal again.
//
// Parameters:
//   - data: A string containing the JSON data to be parsed.
//
// Returns:
//   - T: The unmarshalled value of type T.
//   - error: An error if the unmarshalling process fails, or nil if successful.
//     Notably, it returns ErrExpectedJSONArray (wrapped in a fmt.Errorf) if the
//     target type is an array or slice but the input is not a JSON array.
//
// The function uses the following process:
//  1. Prepares the JSON string for unmarshalling.
//  2. Checks if the input is empty.
//  3. Determines if the target type is an array or slice.
//  4. Attempts to unmarshal the JSON into the value.
//  5. If unmarshalling fails and the target is an array/slice, checks if the input is a JSON array.
//     If not, it returns ErrExpectedJSONArray.
//  6. If repair is needed, attempts to repair and unmarshal the repaired JSON.
//
// Usage:
//
//	type MyStruct struct {
//	    // fields
//	}
//	result, err := SafeUnmarshallToAny[MyStruct](jsonString)
//	if err != nil {
//	    if errors.Is(err, ErrExpectedJSONArray) {
//	        // Handle case where array was expected but not received
//	    } else {
//	        // Handle other errors
//	    }
//	}
func ToAny[T any](raw []byte) (T, error) {
	data := prepareStringWithJsonForUnmarshalling(string(raw))

	// remove newlines - unnecessary for parsing and can cause issues with multi-line strings
	data = strings.ReplaceAll(data, "\n", "")

	if len(data) == 0 {
		var zero T
		return zero, fmt.Errorf("empty input string")
	}

	var response T
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {

		var temp T
		valueType := reflect.ValueOf(temp).Type()
		isArray := valueType.Kind() == reflect.Array || valueType.Kind() == reflect.Slice

		if isArray && !isJSONArray(data) {
			var zero T
			return zero, fmt.Errorf("%w: got %s", ErrExpectedJSONArray, data)
		}

		repairedData, repairErr := repairJSON(data)
		if repairErr != nil {
			var zero T
			return zero, fmt.Errorf("failed to repair JSON: %w", repairErr)
		}

		err = json.Unmarshal([]byte(repairedData), &response)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("failed to parse repaired JSON into struct: %w", err)
		}
	}
	return response, nil
}

// isJSONArray checks if the input string represents a JSON array.
//
// This function scans the input string, skipping any leading whitespace,
// to determine if it starts with an opening square bracket '[', which
// indicates the beginning of a JSON array.
//
// Parameters:
//   - data: A string containing the JSON data to be checked.
//
// Returns:
//   - bool: true if the input represents a JSON array, false otherwise.
//
// Note: This function only checks the first non-whitespace character
// and does not validate the entire JSON structure.
func isJSONArray(data string) bool {
	for _, c := range data {
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			continue
		}
		return c == '['
	}
	return false
}

// prepareStringWithJsonForUnmarshalling takes a string which may contain a JSON object and returns the string
// with the left and right trimmed to the outermost brackets.
func prepareStringWithJsonForUnmarshalling(data string) string {
	trimmedData := strings.TrimSpace(data)

	if len(trimmedData) == 0 {
		return ""
	}

	// Check if the first character is '{' and the last character is '}'
	if (trimmedData[0] == '{' && trimmedData[len(trimmedData)-1] == '}') ||
		(trimmedData[0] == '[' && trimmedData[len(trimmedData)-1] == ']') {
		return trimmedData
	} else {
		// Find the first occurrence of a JSON object.
		startIndex := -1
		braceCount := 0
		for i, char := range data {
			if char == '{' {
				braceCount++
				if startIndex == -1 {
					startIndex = i
				}
			} else if char == '}' {
				braceCount--
				if braceCount == 0 && startIndex != -1 {
					return data[startIndex : i+1]
				}
			}
		}

		slog.Error("Error parsing JSON from string", "data", data)
		return ""
	}
}
