package jobj

import (
	"fmt"
	"log/slog"
	"reflect"
)

// ParseFuncSchema generates a Schema from a function's parameter struct.
//
// It expects a function with signature func(context.Context, any) (string, error)
//
// It creates a schema based on the fields of the second parameter's struct type.
//
// Returns an error if the function signature doesn't match the expected format.
func ParseFuncSchema(function interface{}) (Schema, error) {

	if function == nil {
		return Schema{}, fmt.Errorf("received nil function; must provide a valid function")
	}

	funcType := reflect.TypeOf(function)
	if funcType.Kind() != reflect.Func {
		return Schema{}, fmt.Errorf("received %v, expected a function type", funcType.Kind())
	}

	if funcType.NumIn() != 2 || funcType.NumOut() != 2 {
		return Schema{}, fmt.Errorf(
			"invalid function signature: got %d inputs and %d outputs, expected signature: func(context.Context, any) (string, error)",
			funcType.NumIn(),
			funcType.NumOut(),
		)
	}

	if funcType.In(0).String() != "context.Context" {
		return Schema{}, fmt.Errorf(
			"first parameter must be context.Context, got %s",
			funcType.In(0).String(),
		)
	}

	paramType := funcType.In(1)

	if paramType.Kind() == reflect.Ptr {
		paramType = paramType.Elem()
	}

	if paramType.Kind() != reflect.Struct {
		return Schema{}, fmt.Errorf(
			"second parameter must be a struct or pointer to struct, got %v. Consider wrapping your parameter in a struct",
			paramType,
		)
	}

	schema := Schema{
		Name:        paramType.Name(),
		Description: fmt.Sprintf("Schema for %s function parameters", paramType.Name()),
		Fields:      make([]*Field, 0),
		UseXML:      false,
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
		return Schema{}, fmt.Errorf(
			"no valid fields found in struct %s. Ensure fields are exported and of supported types",
			paramType.Name(),
		)
	}

	return schema, nil
}

// createFieldFromStructField converts a reflect.StructField to a Field
func createFieldFromStructField(field reflect.StructField) *Field {
	var jobjField *Field

	switch field.Type.Kind() {
	case reflect.String:
		jobjField = Text(field.Name)
	case reflect.Bool:
		jobjField = Bool(field.Name)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		jobjField = Int(field.Name)
	case reflect.Float32, reflect.Float64:
		jobjField = Float(field.Name)
	case reflect.Struct:
		if field.Type.String() == "time.Time" {
			jobjField = Date(field.Name)
		} else {
			subFields := make([]*Field, 0)
			for i := 0; i < field.Type.NumField(); i++ {
				subField := createFieldFromStructField(field.Type.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField = Object(field.Name, subFields)
		}
	case reflect.Slice, reflect.Array:
		elemType := field.Type.Elem()
		if elemType.Kind() == reflect.Struct {
			subFields := make([]*Field, 0)
			for i := 0; i < elemType.NumField(); i++ {
				subField := createFieldFromStructField(elemType.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField = Array(field.Name, subFields)
		}
	default:
		slog.Warn("Unsupported field type", "field", field.Name, "type", field.Type.Kind())
		return nil
	}

	if jobjField != nil {
		if desc, ok := field.Tag.Lookup("desc"); ok {
			jobjField.Desc(desc)
		}

		if req, ok := field.Tag.Lookup("required"); ok && req == "true" {
			jobjField.Required()
		}
	}

	return jobjField
}
