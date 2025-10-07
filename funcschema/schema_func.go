package funcschema

import (
	"context"
	"fmt"
	"github.com/mhpenta/jobj"
	"log/slog"
	"reflect"
	"strings"
)

// NewSchemaFromFuncV2 creates a jobj.Schema from a function's second parameter type.
// Unlike NewSchemaFromFunc, this version uses generics to enforce the function signature
// at compile time.
//
// Type parameters:
//   - T: The type of the second parameter (must be a struct type)
//   - R: The return type of the function (can be any type)
//
// The function accepts handlers with the signature:
//
//	func(context.Context, T) (R, error)
//
// Returns a Schema describing the structure of type T and any error encountered.
// An error is returned if T is not a struct type or if T has no exported fields
// of supported types.
func NewSchemaFromFuncV2[T any, R any](function func(context.Context, T) (R, error)) (jobj.Schema, error) {
	var zero T
	paramType := reflect.TypeOf(zero)

	if paramType.Kind() == reflect.Ptr {
		paramType = paramType.Elem()
	}

	if paramType.Kind() != reflect.Struct {
		return jobj.Schema{}, fmt.Errorf("second parameter must be a struct")
	}

	schema := jobj.Schema{
		Name:        paramType.Name(),
		Description: fmt.Sprintf("Schema for %s function parameters", paramType.Name()),
		Fields:      make([]*jobj.Field, 0),
	}

	for i := 0; i < paramType.NumField(); i++ {
		field := paramType.Field(i)

		if !field.IsExported() {
			continue
		}

		jobjField := createFieldFromStructField(field)
		if jobjField != nil {
			schema.Fields = append(schema.Fields, jobjField)
		}
	}

	if len(schema.Fields) == 0 {
		return jobj.Schema{}, fmt.Errorf(
			"no valid fields found in struct %s. Ensure fields are exported and of supported types",
			paramType.Name(),
		)
	}

	return schema, nil
}

// NewSchemasFromFunc creates jobj.Schemas from a function's input and output types.
// Uses generics to enforce the function signature at compile time.
//
// Type parameters:
//   - T: The type of the second parameter (must be a struct type)
//   - R: The return type of the function (must be a struct type)
//
// The function accepts handlers with the signature:
//
//	func(context.Context, T) (R, error)
//
// Returns input and output Schemas describing the structure of types T and R respectively,
// and any error encountered. An error is returned if T or R are not struct types, or if
// they have no exported fields of supported types.
func NewSchemasFromFunc[T any, R any](function func(context.Context, T) (R, error)) (input jobj.Schema, output jobj.Schema, err error) {
	// Create input schema from T
	var zeroInput T
	inputType := reflect.TypeOf(zeroInput)

	if inputType.Kind() == reflect.Ptr {
		inputType = inputType.Elem()
	}

	if inputType.Kind() != reflect.Struct {
		return jobj.Schema{}, jobj.Schema{}, fmt.Errorf("input parameter type must be a struct")
	}

	input = jobj.Schema{
		Name:        inputType.Name(),
		Description: fmt.Sprintf("Input schema for %s function parameters", inputType.Name()),
		Fields:      make([]*jobj.Field, 0),
	}

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)

		if !field.IsExported() {
			continue
		}

		jobjField := createFieldFromStructField(field)
		if jobjField != nil {
			input.Fields = append(input.Fields, jobjField)
		}
	}

	if len(input.Fields) == 0 {
		return jobj.Schema{}, jobj.Schema{}, fmt.Errorf(
			"no valid fields found in input struct %s. Ensure fields are exported and of supported types",
			inputType.Name(),
		)
	}

	// Create output schema from R
	var zeroOutput R
	outputType := reflect.TypeOf(zeroOutput)

	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}

	if outputType.Kind() != reflect.Struct {
		return jobj.Schema{}, jobj.Schema{}, fmt.Errorf("output return type must be a struct")
	}

	output = jobj.Schema{
		Name:        outputType.Name(),
		Description: fmt.Sprintf("Output schema for %s function return value", outputType.Name()),
		Fields:      make([]*jobj.Field, 0),
	}

	for i := 0; i < outputType.NumField(); i++ {
		field := outputType.Field(i)

		if !field.IsExported() {
			continue
		}

		jobjField := createFieldFromStructField(field)
		if jobjField != nil {
			output.Fields = append(output.Fields, jobjField)
		}
	}

	if len(output.Fields) == 0 {
		return jobj.Schema{}, jobj.Schema{}, fmt.Errorf(
			"no valid fields found in output struct %s. Ensure fields are exported and of supported types",
			outputType.Name(),
		)
	}

	return input, output, nil
}

// NewSchemaFromFunc creates a Schema from a function's second parameter type.
// Returns an error if the function doesn't match signature func(context.Context, any)
// or if the second parameter is not a struct type.
func NewSchemaFromFunc(function interface{}) (jobj.Schema, error) {
	if function == nil {
		return jobj.Schema{}, fmt.Errorf("received nil function; must provide a valid function")
	}

	funcType := reflect.TypeOf(function)
	if funcType.Kind() != reflect.Func {
		return jobj.Schema{}, fmt.Errorf("received %v, expected a function type", funcType.Kind())
	}

	if funcType.In(0).String() != "context.Context" {
		return jobj.Schema{}, fmt.Errorf(
			"first parameter must be context.Context, got %s",
			funcType.In(0).String(),
		)
	}

	paramType := funcType.In(1)

	if paramType.Kind() == reflect.Ptr {
		paramType = paramType.Elem()
	}

	if paramType.Kind() != reflect.Struct {
		return jobj.Schema{}, fmt.Errorf(
			"second parameter must be a struct or pointer to struct, got %v. Consider wrapping your parameter in a struct",
			paramType,
		)
	}

	schema := jobj.Schema{
		Name:        paramType.Name(),
		Description: fmt.Sprintf("Schema for %s function parameters", paramType.Name()),
		Fields:      make([]*jobj.Field, 0),
	}

	for i := 0; i < paramType.NumField(); i++ {
		field := paramType.Field(i)

		if !field.IsExported() {
			continue
		}

		jobjField := createFieldFromStructField(field)
		if jobjField != nil {
			schema.Fields = append(schema.Fields, jobjField)
		}
	}

	if len(schema.Fields) == 0 {
		return jobj.Schema{}, fmt.Errorf(
			"no valid fields found in struct %s. Ensure fields are exported and of supported types",
			paramType.Name(),
		)
	}

	return schema, nil
}

// createFieldFromStructField converts a reflect.StructField to a Field
func createFieldFromStructField(field reflect.StructField) *jobj.Field {
	var jobjField *jobj.Field

	// Get the field name from JSON tag if present, otherwise use the Go field name
	fieldName := field.Name
	if jsonTag, ok := field.Tag.Lookup("json"); ok {
		// Parse the json tag to get the field name (before any comma)
		if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
			fieldName = jsonTag[:commaIdx]
		} else {
			fieldName = jsonTag
		}
		// Skip field if json tag is "-"
		if fieldName == "-" {
			return nil
		}
	}

	switch field.Type.Kind() {
	case reflect.String:
		jobjField = jobj.Text(fieldName)
	case reflect.Bool:
		jobjField = jobj.Bool(fieldName)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		jobjField = jobj.Int(fieldName)
	case reflect.Float32, reflect.Float64:
		jobjField = jobj.Float(fieldName)
	case reflect.Struct:
		if field.Type.String() == "time.Time" {
			jobjField = jobj.Date(fieldName)
		} else {
			subFields := make([]*jobj.Field, 0)
			for i := 0; i < field.Type.NumField(); i++ {
				subField := createFieldFromStructField(field.Type.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField = jobj.Object(fieldName, subFields)
		}
	case reflect.Slice, reflect.Array:
		elemType := field.Type.Elem()
		if elemType.Kind() == reflect.Struct {
			// Array of structs
			subFields := make([]*jobj.Field, 0)
			for i := 0; i < elemType.NumField(); i++ {
				subField := createFieldFromStructField(elemType.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField = jobj.Array(fieldName, subFields)
		} else {
			// Array of primitives - use ArrayOf with the appropriate item type
			var itemType jobj.DataType
			switch elemType.Kind() {
			case reflect.String:
				itemType = jobj.TypeString
			case reflect.Bool:
				itemType = jobj.TypeBoolean
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				itemType = jobj.TypeInteger
			case reflect.Float32, reflect.Float64:
				itemType = jobj.TypeNumber
			default:
				slog.Warn("Unsupported array element type", "field", field.Name, "elemType", elemType.Kind())
				return nil
			}
			jobjField = jobj.ArrayOf(fieldName, itemType)
		}
	default:
		slog.Warn("Unsupported field type", "field", field.Name, "type", field.Type.Kind())
		return nil
	}

	if jobjField != nil {
		// Support both "desc" and "description" tags, with "desc" taking precedence
		if desc, ok := field.Tag.Lookup("desc"); ok {
			jobjField.Desc(desc)
		} else if desc, ok := field.Tag.Lookup("description"); ok {
			jobjField.Desc(desc)
		}

		if req, ok := field.Tag.Lookup("required"); ok && req == "true" {
			jobjField.Required()
		}
	}

	return jobjField
}
