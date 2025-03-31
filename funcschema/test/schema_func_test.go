package test

import (
	"context"
	"fmt"
	"github.com/mhpenta/jobj/funcschema"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SearchTool struct{}

type SearchToolParams struct {
	ID    int    `desc:"ID of item to search" required:"true" `
	Query string `desc:"Query to search for, e.g., xyz" required:"true"`
}

func (f *SearchTool) SearchForData(ctx context.Context, params SearchToolParams) (string, error) {
	return fmt.Sprintf("Searching for %s", params.Query), nil
}

func TestSearchTool_Parameters(t *testing.T) {
	searchTool := &SearchTool{}

	schema, err := funcschema.NewSchemaFromFunc(searchTool.SearchForData)
	if err != nil {
		t.Error(err)
	}

	correct := `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "definitions": {
    "SearchToolParams": {
      "additionalProperties": false,
      "properties": {
        "ID": {
          "description": "ID of item to search",
          "type": "integer"
        },
        "Query": {
          "description": "Query to search for, e.g., xyz",
          "type": "string"
        }
      },
      "required": [
        "ID",
        "Query"
      ],
      "type": "object"
    }
  },
  "$ref": "#/definitions/SearchToolParams"
}`

	assert.Equal(t, correct, schema.GetSchemaString())

	paramsMap := funcschema.GetPropertiesMap(schema)

	if paramsMap["properties"].(map[string]interface{})["ID"].(map[string]string)["type"] != "integer" {
		t.Error("ID field should be integer type")
	}
	if paramsMap["properties"].(map[string]interface{})["Query"].(map[string]string)["type"] != "string" {
		t.Error("Query field should be string type")
	}
	if paramsMap["type"] != "object" {
		t.Error("type should be object")
	}
	if paramsMap["additionalProperties"] != false {
		t.Error("additionalProperties should be false")
	}
	if _, ok := paramsMap["properties"].(map[string]interface{}); !ok {
		t.Error("properties should be a map")
	}
}

func TestSearchTool_ParametersV2(t *testing.T) {
	searchTool := &SearchTool{}

	schema, err := funcschema.NewSchemaFromFuncV2(searchTool.SearchForData)
	if err != nil {
		t.Error(err)
	}

	correct := `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "definitions": {
    "SearchToolParams": {
      "additionalProperties": false,
      "properties": {
        "ID": {
          "description": "ID of item to search",
          "type": "integer"
        },
        "Query": {
          "description": "Query to search for, e.g., xyz",
          "type": "string"
        }
      },
      "required": [
        "ID",
        "Query"
      ],
      "type": "object"
    }
  },
  "$ref": "#/definitions/SearchToolParams"
}`

	assert.Equal(t, correct, schema.GetSchemaString())

	paramsMap := funcschema.GetPropertiesMap(schema)

	if paramsMap["properties"].(map[string]interface{})["ID"].(map[string]string)["type"] != "integer" {
		t.Error("ID field should be integer type")
	}
	if paramsMap["properties"].(map[string]interface{})["Query"].(map[string]string)["type"] != "string" {
		t.Error("Query field should be string type")
	}
	if paramsMap["type"] != "object" {
		t.Error("type should be object")
	}
	if paramsMap["additionalProperties"] != false {
		t.Error("additionalProperties should be false")
	}
	if _, ok := paramsMap["properties"].(map[string]interface{}); !ok {
		t.Error("properties should be a map")
	}
}
