package funcschema

import (
	"context"
	"fmt"
	"github.com/mhpenta/jobj"
	"log/slog"
	"reflect"
)

// NewSchemaFromFuncV2 creates a Schema from a function's second parameter type.
// Returns an error if the function doesn't match signature func(context.Context, T) (R, error)
// or if the second parameter is not a struct type.
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
		return jobj.Schema{}, fmt.Errorf(
			"no valid fields found in struct %s. Ensure fields are exported and of supported types",
			paramType.Name(),
		)
	}

	return schema, nil
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

	switch field.Type.Kind() {
	case reflect.String:
		jobjField = jobj.Text(field.Name)
	case reflect.Bool:
		jobjField = jobj.Bool(field.Name)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		jobjField = jobj.Int(field.Name)
	case reflect.Float32, reflect.Float64:
		jobjField = jobj.Float(field.Name)
	case reflect.Struct:
		if field.Type.String() == "time.Time" {
			jobjField = jobj.Date(field.Name)
		} else {
			subFields := make([]*jobj.Field, 0)
			for i := 0; i < field.Type.NumField(); i++ {
				subField := createFieldFromStructField(field.Type.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField = jobj.Object(field.Name, subFields)
		}
	case reflect.Slice, reflect.Array:
		elemType := field.Type.Elem()
		if elemType.Kind() == reflect.Struct {
			subFields := make([]*jobj.Field, 0)
			for i := 0; i < elemType.NumField(); i++ {
				subField := createFieldFromStructField(elemType.Field(i))
				if subField != nil {
					subFields = append(subFields, subField)
				}
			}
			jobjField = jobj.Array(field.Name, subFields)
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
