package funcschema

import "github.com/mhpenta/jobj"

// GetPropertiesMap returns a map of properties for a schema, often useful when constructing schemas for LLM tool calls
func GetPropertiesMap(schema jobj.Schema) map[string]interface{} {
	return map[string]interface{}{
		"type":                 "object",
		"properties":           schema.FieldsJson(),
		"required":             schema.RequiredFields(),
		"additionalProperties": false,
	}
}
