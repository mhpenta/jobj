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
// Internally, we use this to transform Go functions into "Tools" for LLM Agents.
func SafeSchemaFromFunc[T any, R any](function func(context.Context, T) (R, error)) (map[string]interface{}, error) {
	schema, err := NewSchemaFromFuncV2(function)
	if err != nil {
		return nil, err
	}
	return GetPropertiesMap(schema), nil
}
