package funcschema

import (
	"context"
	"errors"
	"log/slog"
)

var ErrSchemaGeneration = errors.New("failed to generate schema from function, review the function signature")

// SafeSchemaFromFunc attempts to generate a schema from a function and handles errors uniformly
// It returns the properties map and a boolean indicating if the operation was successful
func SafeSchemaFromFunc[T any, R any](function func(context.Context, T) (R, error)) (map[string]interface{}, error) {
	schema, err := NewSchemaFromFuncV2(function)
	if err != nil {
		slog.Error("Failed to parse function schema", "error", err)
		return nil, ErrSchemaGeneration
	}
	return GetPropertiesMap(schema), nil
}
