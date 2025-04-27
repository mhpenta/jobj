package funcschema

import (
	"fmt"
	"github.com/mhpenta/jobj"
	"reflect"
)

// SchemaFromStruct generates a jobj.Schema from any struct type using generics.
// This provides a direct way to create schemas without going through function signatures,
// making it easier to use when you just have a struct type rather than a function.
//
// The generic implementation offers several advantages:
// - No need to create an instance of the struct
// - Type safety through compile-time type checking
// - Clear syntax that expresses intent: SchemaFromStruct[User]()
//
// Example:
//
//	type User struct {
//	    Name string `json:"name" desc:"User's full name" required:"true"`
//	    Age  int    `json:"age" desc:"User's age"`
//	}
//
//	schema, err := SchemaFromStruct[User]()
func SchemaFromStruct[T any]() (jobj.Schema, error) {
	var zero T
	return createSchemaFromType(reflect.TypeOf(zero))
}

// createSchemaFromType generates a jobj.Schema from a reflect.Type.
// This function expects a struct type and will return an error if provided with
// any other kind of type.
//
// It creates a schema with fields corresponding to each exported field in the struct,
// processing tags for descriptions and required fields.
//
// This is the underlying implementation used by SchemaFromStruct and the function
// schema generators.
func createSchemaFromType(t reflect.Type) (jobj.Schema, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return jobj.Schema{}, fmt.Errorf("expected struct type, got %v", t.Kind())
	}

	schema := jobj.Schema{
		Name:        t.Name(),
		Description: fmt.Sprintf("Schema for %s", t.Name()),
		Fields:      make([]*jobj.Field, 0, t.NumField()),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

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
			t.Name(),
		)
	}

	return schema, nil
}
