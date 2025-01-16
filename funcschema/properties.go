package funcschema

import "github.com/mhpenta/jobj"

func GetPropertiesMap(schema jobj.Schema) map[string]interface{} {
	return map[string]interface{}{
		"type":                 "object",
		"properties":           schema.FieldsJson(),
		"required":             schema.RequiredFields(),
		"additionalProperties": false,
	}
}
