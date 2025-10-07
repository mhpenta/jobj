package funcschema

import "github.com/mhpenta/jobj"

// GetPropertiesMap returns a map of properties for a schema, often useful when constructing schemas for LLM tool calls
func GetPropertiesMap(schema jobj.Schema) map[string]interface{} {
	// Check if this is a non-struct return type (has RootField)
	if schema.RootField != nil {
		return generateSchemaForField(schema.RootField)
	}

	// Default: struct type with Fields
	return map[string]interface{}{
		"type":                 "object",
		"properties":           schema.FieldsJson(),
		"required":             schema.RequiredFields(),
		"additionalProperties": false,
	}
}

// generateSchemaForField creates a JSON schema for a single field (used for non-struct return types)
func generateSchemaForField(field *jobj.Field) map[string]interface{} {
	schema := make(map[string]interface{})

	switch field.ValueType {
	case jobj.TypeArray:
		schema["type"] = "array"
		if field.ArrayItemType != "" {
			// Array of primitives
			schema["items"] = map[string]interface{}{
				"type": string(field.ArrayItemType),
			}
		} else if field.SubFields != nil {
			// Array of objects
			schema["items"] = map[string]interface{}{
				"type":       "object",
				"properties": generatePropertiesForFields(field.SubFields),
			}
		}
	case jobj.TypeObject:
		schema["type"] = "object"
		if field.AdditionalProperties {
			// This is a map
			if field.AdditionalPropertiesType != "" {
				// Map with primitive values
				schema["additionalProperties"] = map[string]interface{}{
					"type": string(field.AdditionalPropertiesType),
				}
			} else if field.AdditionalPropertiesField != nil {
				// Map with complex values
				if field.AdditionalPropertiesField.SubFields == nil {
					// interface{} case
					schema["additionalProperties"] = true
				} else {
					// Struct case
					schema["additionalProperties"] = map[string]interface{}{
						"type":       "object",
						"properties": generatePropertiesForFields(field.AdditionalPropertiesField.SubFields),
					}
				}
			}
		} else if field.SubFields != nil {
			// Regular object with defined properties
			schema["properties"] = generatePropertiesForFields(field.SubFields)
		}
	default:
		// Primitive types
		schema["type"] = string(field.ValueType)
	}

	if field.ValueDescription != "" {
		schema["description"] = field.ValueDescription
	}

	return schema
}

// generatePropertiesForFields is a helper to create properties map from fields
func generatePropertiesForFields(fields []*jobj.Field) map[string]interface{} {
	properties := make(map[string]interface{})
	for _, f := range fields {
		properties[f.ValueName] = generateSchemaForField(f)
	}
	return properties
}
