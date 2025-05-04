package funcschema

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mhpenta/jobj/safeunmarshal"
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

	schema, err := NewSchemaFromFunc(searchTool.SearchForData)
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

	paramsMap := GetPropertiesMap(schema)

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

	schema, err := NewSchemaFromFuncV2(searchTool.SearchForData)
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

	paramsMap := GetPropertiesMap(schema)

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

// TestSafeunmarshalIntegration tests the integration between funcschema and safeunmarshal packages
func TestSafeunmarshalIntegration(t *testing.T) {
	type ToolResult struct {
		Result string `json:"result"`
	}

	executeFunc := func(ctx context.Context, params json.RawMessage) (*ToolResult, error) {
		searchParams, err := safeunmarshal.To[SearchToolParams](params)
		if err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result := &ToolResult{
			Result: fmt.Sprintf("Found results for ID:%d, Query:%s", searchParams.ID, searchParams.Query),
		}
		return result, nil
	}

	// Test with well-formed JSON
	wellFormedJSON := []byte(`{"ID": 42, "Query": "test query"}`)

	// Test with JSON that needs repair (single quotes, unquoted keys)
	needsRepairJSON := []byte(`{ID: 123, 'Query': 'another query'}`)

	// Execute with both inputs
	ctx := context.Background()

	// Test with well-formed JSON
	result1, err := executeFunc(ctx, wellFormedJSON)
	assert.NoError(t, err)
	assert.Equal(t, "Found results for ID:42, Query:test query", result1.Result)

	// Test with JSON that needs repair
	result2, err := executeFunc(ctx, needsRepairJSON)
	assert.NoError(t, err)
	assert.Equal(t, "Found results for ID:123, Query:another query", result2.Result)

	// For SafeSchemaFromFunc with JSON parameters, we need to use a function with a struct parameter
	// rather than json.RawMessage
	searchFunc := func(ctx context.Context, params SearchToolParams) (*ToolResult, error) {
		return &ToolResult{
			Result: fmt.Sprintf("Found results for ID:%d, Query:%s", params.ID, params.Query),
		}, nil
	}

	// Verify schema generation works correctly
	schema, err := SafeSchemaFromFunc(searchFunc)
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	// Verify properties exist in the schema
	properties, ok := schema["properties"].(map[string]interface{})
	assert.True(t, ok)

	idField, ok := properties["ID"].(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, "integer", idField["type"])

	queryField, ok := properties["Query"].(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, "string", queryField["type"])
}
