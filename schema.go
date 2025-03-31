package jobj

import (
	"encoding/json"
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
	UseXML      bool
	Fields      []*Field
}

func (r *Schema) GetDescription() string {
	return r.Description
}

func (r *Schema) GetFields() []*Field {
	return r.Fields
}

func (r *Schema) GetSchemaString() string {
	if r.UseXML {
		return r.GetXMLSchemaString()
	}

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

// Validate checks if the response fields are present in the container struct, which must be passed into it since the
// embedded Schema struct does not have access to the container struct. Validation proves that the schema requested
// is able to be Marshalled by the container struct. Mostly used in testing. Returns true if all fields are present,
func (r *Schema) Validate(container interface{}) bool {
	containerValue := reflect.ValueOf(container)
	containerType := containerValue.Type()
	if containerType.Kind() != reflect.Struct && !(containerType.Kind() == reflect.Ptr && containerType.Elem().Kind() == reflect.Struct) {
		return false
	}
	if containerType.Kind() == reflect.Ptr {
		containerValue = containerValue.Elem()
		containerType = containerType.Elem()
	}
	fieldNames := make(map[string]struct{}, containerType.NumField())
	for i := 0; i < containerType.NumField(); i++ {
		field := containerType.Field(i)
		if field.Name != "Schema" {
			fieldNames[field.Tag.Get("json")] = struct{}{}
		}
	}
	for _, field := range r.Fields {
		fieldName := strings.ToLower(field.ValueName)
		if _, exists := fieldNames[fieldName]; !exists {
			return false
		}
	}
	return true
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
