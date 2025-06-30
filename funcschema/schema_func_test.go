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

// TestJSONTagNames tests that JSON tag names are properly used in schema generation
func TestJSONTagNames(t *testing.T) {
	type XBRLToolParams struct {
		// Operation specifies what action to perform
		Operation string `json:"operation" desc:"Operation to perform" required:"true"`

		// CIK is the company CIK number
		CIK string `json:"cik,omitempty" desc:"Company CIK number"`

		// AccessionNumber is the SEC accession number
		AccessionNumber string `json:"accession_number,omitempty" desc:"SEC accession number"`

		// DocumentType filters documents by type
		DocumentType string `json:"document_type,omitempty" desc:"Filter documents by type"`

		// ReportNumber is the report/table number
		ReportNumber int `json:"report_number,omitempty" desc:"Report/table number"`

		// IgnoredField should not appear in schema
		IgnoredField string `json:"-"`

		// NoJSONTag should use the Go field name
		NoJSONTag string `desc:"Field without JSON tag"`
	}

	handler := func(ctx context.Context, params XBRLToolParams) (string, error) {
		return "processed", nil
	}

	// Test SafeSchemaFromFunc
	schema, err := SafeSchemaFromFunc(handler)
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	// Verify the schema structure
	assert.Equal(t, "object", schema["type"])
	assert.Equal(t, false, schema["additionalProperties"])

	properties, ok := schema["properties"].(map[string]interface{})
	assert.True(t, ok, "properties should be a map")

	// Check that JSON tag names are used
	_, hasOperation := properties["operation"]
	assert.True(t, hasOperation, "should have 'operation' field (from json tag)")

	_, hasCIK := properties["cik"]
	assert.True(t, hasCIK, "should have 'cik' field (from json tag)")

	_, hasAccessionNumber := properties["accession_number"]
	assert.True(t, hasAccessionNumber, "should have 'accession_number' field (from json tag)")

	_, hasDocumentType := properties["document_type"]
	assert.True(t, hasDocumentType, "should have 'document_type' field (from json tag)")

	_, hasReportNumber := properties["report_number"]
	assert.True(t, hasReportNumber, "should have 'report_number' field (from json tag)")

	// Check that json:"-" field is excluded
	_, hasIgnored := properties["IgnoredField"]
	assert.False(t, hasIgnored, "should not have 'IgnoredField' (json:\"-\")")

	// Check that field without JSON tag uses Go field name
	_, hasNoJSONTag := properties["NoJSONTag"]
	assert.True(t, hasNoJSONTag, "should have 'NoJSONTag' field (no json tag)")

	// Verify field types and descriptions
	operationField := properties["operation"].(map[string]string)
	assert.Equal(t, "string", operationField["type"])
	assert.Equal(t, "Operation to perform", operationField["description"])

	cikField := properties["cik"].(map[string]string)
	assert.Equal(t, "string", cikField["type"])
	assert.Equal(t, "Company CIK number", cikField["description"])

	reportNumberField := properties["report_number"].(map[string]string)
	assert.Equal(t, "integer", reportNumberField["type"])
	assert.Equal(t, "Report/table number", reportNumberField["description"])

	// Check required fields
	required, ok := schema["required"].([]string)
	assert.True(t, ok)
	assert.Contains(t, required, "operation", "operation should be required")
	assert.NotContains(t, required, "cik", "cik should not be required")
	assert.NotContains(t, required, "accession_number", "accession_number should not be required")
}
