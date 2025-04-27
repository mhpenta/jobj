package jobj

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

type CreatableSchema interface {
	CreateDescription() CreatableSchema
	CreateFields() CreatableSchema
	GetDescription() string
	GetFields() []*Field
}

type Schema struct {
	Name        string
	Description string
	Fields      []*Field
}

func (r *Schema) GetDescription() string {
	return r.Description
}

func (r *Schema) GetFields() []*Field {
	return r.Fields
}

func (r *Schema) GetSchemaString() string {
	schema := struct {
		Schema      string                 `json:"$schema"`
		Definitions map[string]interface{} `json:"definitions"`
		Reference   string                 `json:"$ref"`
	}{
		Schema: "http://json-schema.org/draft-07/schema#",
		Definitions: map[string]interface{}{
			r.Name: map[string]interface{}{
				"properties":           r.FieldsJson(),
				"type":                 "object",
				"required":             r.RequiredFields(),
				"additionalProperties": false,
			},
		},
		Reference: "#/definitions/" + r.Name,
	}

	schemaJson, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		// In theory, this could be problematic - in practice, however, there are very few ways we could experience
		// an error: (1) the system ran out of memory, (2) a field values contained invalid UTF-8 characters
		// or (3) a type was added that implements a custom MarshalJSON method that returns an error.
		//
		// Since these are unlikely, we return an empty string and log the error.
		slog.Error("Error marshalling JSON schema", "err", err)
		return ""
	}
	return string(schemaJson)
}

func (r *Schema) FieldsJson() map[string]interface{} {
	properties := make(map[string]interface{}, len(r.Fields))
	for _, field := range r.Fields {
		if field.ValueAnyOf != nil {
			anyOf := make([]map[string]interface{}, 0, len(field.ValueAnyOf))
			for _, enum := range field.ValueAnyOf {
				anyOf = append(anyOf, map[string]interface{}{
					"const":       enum.Const,
					"description": enum.Description,
				})
			}

			fieldProps := map[string]interface{}{
				"anyOf": anyOf,
			}
			if field.ValueDescription != "" {
				fieldProps["description"] = field.ValueDescription
			}

			properties[field.ValueName] = fieldProps
			continue
		}

		if field.ValueType == "array" && field.SubFields != nil {
			arrayFieldProperties := make(map[string]interface{})
			requiredFields := []string{}

			for _, subField := range field.SubFields {
				if subField.ValueRequired {
					requiredFields = append(requiredFields, subField.ValueName)
				}
			}

			for _, subField := range field.SubFields {
				if subField.ValueAnyOf != nil {
					anyOf := make([]map[string]interface{}, 0, len(subField.ValueAnyOf))
					for _, enum := range subField.ValueAnyOf {
						anyOf = append(anyOf, map[string]interface{}{
							"const":       enum.Const,
							"description": enum.Description,
						})
					}

					fieldProps := map[string]interface{}{
						"anyOf": anyOf,
					}
					if subField.ValueDescription != "" {
						fieldProps["description"] = subField.ValueDescription
					}

					arrayFieldProperties[subField.ValueName] = fieldProps
					continue
				}

				if subField.ValueType == "object" && subField.SubFields != nil {
					objectFieldProperties := processObjectFields(subField.SubFields)
					arrayFieldProperties[subField.ValueName] = map[string]interface{}{
						"type":        subField.ValueType,
						"description": subField.ValueDescription,
						"properties":  objectFieldProperties,
						"required":    subField.getRequiredFields(),
					}
					continue
				}

				arrayFieldProperties[subField.ValueName] = map[string]string{
					"type":        string(subField.ValueType),
					"description": subField.ValueDescription,
				}
			}

			properties[field.ValueName] = map[string]interface{}{
				"type":                 field.ValueType,
				"description":          field.ValueDescription,
				"additionalProperties": field.AdditionalProperties,
				"items": map[string]interface{}{
					"type":       "object",
					"properties": arrayFieldProperties,
					"required":   requiredFields,
				},
			}
			continue
		}

		if field.ValueType == "object" {
			objectFieldProperties := processObjectFields(field.SubFields)
			properties[field.ValueName] = map[string]interface{}{
				"type":        field.ValueType,
				"description": field.ValueDescription,
				"properties":  objectFieldProperties,
				"required":    field.getRequiredFields(),
			}
			continue
		}

		properties[field.ValueName] = map[string]string{
			"type":        string(field.ValueType),
			"description": field.ValueDescription,
		}
	}
	return properties
}

func (r *Schema) RequiredFields() []string {
	required := make([]string, 0, len(r.Fields))
	for _, field := range r.Fields {
		if field.ValueType == "array" && field.ValueRequired {
			required = append(required, field.ValueName)
			continue
		}

		if field.ValueRequired {
			required = append(required, field.ValueName)
		}
	}

	if len(required) == 0 {
		return nil
	}
	return required
}

// Validate verifies that the schema fields match the struct fields.
// It checks:
// - Every schema field has a corresponding struct field with matching JSON tag
// - Every required struct field (no omitempty tag) has a corresponding schema field
// - Type compatibility between schema fields and struct fields
// Returns detailed error messages for any mismatches.
func (r *Schema) Validate(structPtr interface{}) error {
	val := reflect.ValueOf(structPtr)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	structType := val.Elem().Type()

	schemaFields := make(map[string]*Field)
	for _, field := range r.Fields {
		schemaFields[field.ValueName] = field
	}

	structFields := make(map[string]reflect.StructField)
	requiredStructFields := make(map[string]bool)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Skip the embedded Schema field
		if field.Anonymous && field.Type.Name() == "Schema" {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		parts := strings.Split(jsonTag, ",")
		jsonName := parts[0]
		isOmitEmpty := false

		for _, opt := range parts[1:] {
			if opt == "omitempty" {
				isOmitEmpty = true
				break
			}
		}

		structFields[jsonName] = field
		if !isOmitEmpty {
			requiredStructFields[jsonName] = true
		}
	}

	var errors []string

	for name, schemaField := range schemaFields {
		structField, exists := structFields[name]
		if !exists {
			errors = append(errors, fmt.Sprintf("schema field %q does not exist in struct", name))
			continue
		}

		if !isTypeCompatible(structField.Type, schemaField.ValueType) {
			errors = append(errors, fmt.Sprintf("schema field %q has type %q but struct field has incompatible type %q",
				name, schemaField.ValueType, structField.Type.String()))
		}
	}

	for name, isRequired := range requiredStructFields {
		if !isRequired {
			continue
		}

		_, exists := schemaFields[name]
		if !exists {
			errors = append(errors, fmt.Sprintf("required struct field %q does not exist in schema", name))
		}
	}

	for name, field := range schemaFields {
		if !field.ValueRequired {
			continue
		}

		structField, exists := structFields[name]
		if !exists {
			continue
		}

		jsonTag := structField.Tag.Get("json")
		if strings.Contains(jsonTag, "omitempty") {
			errors = append(errors, fmt.Sprintf("schema field %q is required but struct field has omitempty tag", name))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("schema validation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (f *Field) getRequiredFields() []string {
	var required []string
	for _, subField := range f.SubFields {
		if subField.ValueRequired {
			required = append(required, subField.ValueName)
		}
	}
	return required
}

func processObjectFields(fields []*Field) map[string]interface{} {
	objectFieldProperties := make(map[string]interface{})
	for _, field := range fields {
		if field.ValueAnyOf != nil {
			anyOf := make([]map[string]interface{}, 0, len(field.ValueAnyOf))
			for _, enum := range field.ValueAnyOf {
				anyOf = append(anyOf, map[string]interface{}{
					"const":       enum.Const,
					"description": enum.Description,
				})
			}

			fieldProps := map[string]interface{}{
				"anyOf": anyOf,
			}
			if field.ValueDescription != "" {
				fieldProps["description"] = field.ValueDescription
			}

			objectFieldProperties[field.ValueName] = fieldProps
			continue
		}

		objectFieldProperties[field.ValueName] = map[string]string{
			"type":        string(field.ValueType),
			"description": field.ValueDescription,
		}
	}
	return objectFieldProperties
}

func isTypeCompatible(goType reflect.Type, schemaType DataType) bool {
	if goType.Kind() == reflect.Ptr {
		goType = goType.Elem()
	}

	switch schemaType {
	case "string":
		return goType.Kind() == reflect.String
	case "number":
		return goType.Kind() == reflect.Float32 || goType.Kind() == reflect.Float64
	case "integer":
		return goType.Kind() == reflect.Int || goType.Kind() == reflect.Int8 ||
			goType.Kind() == reflect.Int16 || goType.Kind() == reflect.Int32 ||
			goType.Kind() == reflect.Int64 || goType.Kind() == reflect.Uint ||
			goType.Kind() == reflect.Uint8 || goType.Kind() == reflect.Uint16 ||
			goType.Kind() == reflect.Uint32 || goType.Kind() == reflect.Uint64
	case "boolean":
		return goType.Kind() == reflect.Bool
	case "array":
		return goType.Kind() == reflect.Slice || goType.Kind() == reflect.Array
	case "object":
		return goType.Kind() == reflect.Struct || goType.Kind() == reflect.Map
	case "anyOf":
		return true
	default:
		// For unknown types, we're permissive
		return true
	}
}
