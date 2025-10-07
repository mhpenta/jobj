package funcschema

import (
	"context"
)

// SafeSchemaFromFunc generates a properties map from a function's parameter structure.
// This function provides a type-safe wrapper around schema generation using generics.
//
// Type parameters:
//   - T: The type of the second parameter (must be a struct or pointer to struct)
//   - R: The return type of the function (can be any type)
//
// The function accepts handlers with the signature:
//
//	func(context.Context, T) (R, error)
//
// It returns a map containing the JSON Schema properties derived from T's structure.
// If schema generation fails (e.g., T is not a struct or has no exported fields),
// an error is returned and the properties map will be nil.
//
// This function is a convenient alternative to calling NewSchemaFromFuncV2 and
// GetPropertiesMap separately while maintaining compile-time type checking.
//
// Since SafeSchemaFromFunc uses jobj, we use a description field tag of "desc" for
// field descriptions. Both "desc" and "description" tags are supported and treated
// equivalently for convenience.
//
// Internally, we use this to transform Go functions into "Tools" for LLM Agents.
func SafeSchemaFromFunc[T any, R any](function func(context.Context, T) (R, error)) (map[string]interface{}, error) {
	schema, err := NewSchemaFromFuncV2(function)
	if err != nil {
		return nil, err
	}
	return GetPropertiesMap(schema), nil
}

// SafeSchemasFromFunc generates both input and output properties maps from a function's
// parameter and return type structures. This function provides a type-safe wrapper around
// schema generation using generics.
//
// Type parameters:
//   - T: The type of the second parameter (must be a struct or pointer to struct)
//   - R: The return type of the function (must be a struct or pointer to struct)
//
// The function accepts handlers with the signature:
//
//	func(context.Context, T) (R, error)
//
// It returns two maps:
//  1. Input properties map: JSON Schema properties derived from T's structure
//  2. Output properties map: JSON Schema properties derived from R's structure
//
// If schema generation fails (e.g., T or R is not a struct, or has no exported fields),
// an error is returned and both properties maps will be nil.
//
// This function is a convenient alternative to calling NewSchemasFromFunc and
// GetPropertiesMap separately while maintaining compile-time type checking. Use this
// when you need to validate both the input parameters and output structure, which is
// particularly useful for bidirectional LLM tool validation.
//
// Since SafeSchemasFromFunc uses jobj, we use a description field tag of "desc" for
// field descriptions. Both "desc" and "description" tags are supported and treated
// equivalently for convenience.
//
// Internally, we use this to transform Go functions into "Tools" for LLM Agents where
// both input parameter validation and output structure validation are required.
func SafeSchemasFromFunc[T any, R any](function func(context.Context, T) (R, error)) (map[string]interface{}, map[string]interface{}, error) {
	schemaIn, schemaOut, err := NewSchemasFromFunc(function)
	if err != nil {
		return nil, nil, err
	}
	return GetPropertiesMap(schemaIn), GetPropertiesMap(schemaOut), nil
}
