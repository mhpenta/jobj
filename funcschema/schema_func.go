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
	// Use reflect.TypeOf with a typed nil to get the type even for pointer types
	inputType := reflect.TypeOf((*T)(nil)).Elem()

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
	// Use reflect.TypeOf with a typed nil to get the type even for pointer types
	outputType := reflect.TypeOf((*R)(nil)).Elem()

	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}

	// Handle different return types
	if outputType.Kind() == reflect.Struct {
		// Struct return type - use Fields (existing behavior)
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
	} else {
		// Non-struct return type - use RootField (new behavior)
		rootField := createFieldFromType(outputType, "result")
		if rootField == nil {
			return jobj.Schema{}, jobj.Schema{}, fmt.Errorf(
				"unsupported return type %v", outputType,
			)
		}

		typeName := outputType.Name()
		if typeName == "" {
			// For unnamed types like []string or map[string]int
			typeName = outputType.String()
		}

		output = jobj.Schema{
			Name:        typeName,
			Description: fmt.Sprintf("Output schema for %s function return value", typeName),
			RootField:   rootField,
		}
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
// createFieldFromType creates a Field from a reflect.Type (for non-struct return types)
// This is used when the return type is an array, map, or primitive rather than a struct
func createFieldFromType(typ reflect.Type, name string) *jobj.Field {
	var jobjField *jobj.Field

	switch typ.Kind() {
	case reflect.String:
		jobjField = jobj.Text(name)
	case reflect.Bool:
		jobjField = jobj.Bool(name)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		jobjField = jobj.Int(name)
	case reflect.Float32, reflect.Float64:
		jobjField = jobj.Float(name)
	case reflect.Slice, reflect.Array:
		elemType := typ.Elem()
		if elemType.Kind() == reflect.Struct {
			// Array of structs
			subFields := make([]*jobj.Field, 0)
			for i := 0; i < elemType.NumField(); i++ {
				subField := createFieldFromStructField(elemType.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField = jobj.Array(name, subFields)
		} else {
			// Array of primitives
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
				slog.Warn("Unsupported array element type", "type", typ, "elemType", elemType.Kind())
				return nil
			}
			jobjField = jobj.ArrayOf(name, itemType)
		}
	case reflect.Map:
		// Handle map types - maps become objects with additionalProperties
		valueType := typ.Elem()

		// Create an object field with additionalProperties set
		jobjField = &jobj.Field{
			ValueName:            name,
			ValueType:            jobj.TypeObject,
			AdditionalProperties: true,
		}

		// Determine the value type for additionalProperties
		switch valueType.Kind() {
		case reflect.String:
			jobjField.AdditionalPropertiesType = jobj.TypeString
		case reflect.Bool:
			jobjField.AdditionalPropertiesType = jobj.TypeBoolean
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			jobjField.AdditionalPropertiesType = jobj.TypeInteger
		case reflect.Float32, reflect.Float64:
			jobjField.AdditionalPropertiesType = jobj.TypeNumber
		case reflect.Struct:
			// Map with struct values
			subFields := make([]*jobj.Field, 0)
			for i := 0; i < valueType.NumField(); i++ {
				subField := createFieldFromStructField(valueType.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField.AdditionalPropertiesField = &jobj.Field{
				ValueType: jobj.TypeObject,
				SubFields: subFields,
			}
		case reflect.Interface:
			// Map with interface{} values
			jobjField.AdditionalPropertiesField = &jobj.Field{
				ValueType: jobj.TypeObject,
				SubFields: nil,
			}
		default:
			slog.Warn("Unsupported map value type", "type", typ, "valueType", valueType.Kind())
			return nil
		}
	default:
		slog.Warn("Unsupported return type", "type", typ, "kind", typ.Kind())
		return nil
	}

	return jobjField
}

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
	case reflect.Ptr:
		// Handle pointer fields by unwrapping and processing the underlying type
		elemType := field.Type.Elem()
		switch elemType.Kind() {
		case reflect.String:
			jobjField = jobj.Text(fieldName)
		case reflect.Bool:
			jobjField = jobj.Bool(fieldName)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			jobjField = jobj.Int(fieldName)
		case reflect.Float32, reflect.Float64:
			jobjField = jobj.Float(fieldName)
		case reflect.Struct:
			if elemType.String() == "time.Time" {
				jobjField = jobj.Date(fieldName)
			} else {
				subFields := make([]*jobj.Field, 0)
				for i := 0; i < elemType.NumField(); i++ {
					subField := createFieldFromStructField(elemType.Field(i))
					if subField != nil {
						subFields = append(subFields, subField)
					}
				}
				jobjField = jobj.Object(fieldName, subFields)
			}
		default:
			slog.Warn("Unsupported pointer element type", "field", field.Name, "elemType", elemType.Kind())
			return nil
		}
		// Pointer fields are inherently optional, so we don't mark them as required by default
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
	case reflect.Map:
		// Handle map types - maps become objects with additionalProperties
		valueType := field.Type.Elem()

		// Create an object field with additionalProperties set
		jobjField = &jobj.Field{
			ValueName:            fieldName,
			ValueType:            jobj.TypeObject,
			AdditionalProperties: true,
		}

		// Determine the value type for additionalProperties
		switch valueType.Kind() {
		case reflect.String:
			jobjField.AdditionalPropertiesType = jobj.TypeString
		case reflect.Bool:
			jobjField.AdditionalPropertiesType = jobj.TypeBoolean
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			jobjField.AdditionalPropertiesType = jobj.TypeInteger
		case reflect.Float32, reflect.Float64:
			jobjField.AdditionalPropertiesType = jobj.TypeNumber
		case reflect.Struct:
			// Map with struct values
			subFields := make([]*jobj.Field, 0)
			for i := 0; i < valueType.NumField(); i++ {
				subField := createFieldFromStructField(valueType.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField.AdditionalPropertiesField = &jobj.Field{
				ValueType: jobj.TypeObject,
				SubFields: subFields,
			}
		case reflect.Ptr:
			// Map with pointer values - unwrap and process
			elemType := valueType.Elem()
			if elemType.Kind() == reflect.Struct {
				subFields := make([]*jobj.Field, 0)
				for i := 0; i < elemType.NumField(); i++ {
					subField := createFieldFromStructField(elemType.Field(i))
					if subField != nil {
						subFields = append(subFields, subField)
					}
				}
				jobjField.AdditionalPropertiesField = &jobj.Field{
					ValueType: jobj.TypeObject,
					SubFields: subFields,
				}
			} else {
				slog.Warn("Unsupported map pointer value type", "field", field.Name, "valueType", elemType.Kind())
				return nil
			}
		case reflect.Interface:
			// Map with interface{} values - treat as generic object
			// This is common for metadata fields: map[string]interface{}
			// We can't introspect the actual type, so we allow any value type
			jobjField.AdditionalPropertiesField = &jobj.Field{
				ValueType: jobj.TypeObject,
				SubFields: nil, // Empty SubFields means any properties allowed
			}
		default:
			slog.Warn("Unsupported map value type", "field", field.Name, "valueType", valueType.Kind())
			return nil
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
