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

// TestNewSchemasFromFunc tests the NewSchemasFromFunc function that returns both input and output schemas
func TestNewSchemasFromFunc(t *testing.T) {
	type UserInput struct {
		UserID   int    `json:"user_id" desc:"User identifier" required:"true"`
		Username string `json:"username" desc:"Username to search for" required:"true"`
		Email    string `json:"email" desc:"Email address"`
	}

	type UserOutput struct {
		Success bool   `json:"success" desc:"Whether operation succeeded" required:"true"`
		Message string `json:"message" desc:"Status message"`
		UserID  int    `json:"user_id" desc:"User identifier"`
	}

	handler := func(ctx context.Context, input UserInput) (UserOutput, error) {
		return UserOutput{
			Success: true,
			Message: fmt.Sprintf("Processed user %s", input.Username),
			UserID:  input.UserID,
		}, nil
	}

	inputSchema, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	// Verify input schema
	assert.Equal(t, "UserInput", inputSchema.Name)
	assert.Equal(t, "Input schema for UserInput function parameters", inputSchema.Description)
	assert.Len(t, inputSchema.Fields, 3)

	// Verify input fields
	inputMap := GetPropertiesMap(inputSchema)
	inputProps := inputMap["properties"].(map[string]interface{})

	userIDField := inputProps["user_id"].(map[string]string)
	assert.Equal(t, "integer", userIDField["type"])
	assert.Equal(t, "User identifier", userIDField["description"])

	usernameField := inputProps["username"].(map[string]string)
	assert.Equal(t, "string", usernameField["type"])
	assert.Equal(t, "Username to search for", usernameField["description"])

	emailField := inputProps["email"].(map[string]string)
	assert.Equal(t, "string", emailField["type"])
	assert.Equal(t, "Email address", emailField["description"])

	// Verify input required fields
	inputRequired := inputMap["required"].([]string)
	assert.Contains(t, inputRequired, "user_id")
	assert.Contains(t, inputRequired, "username")
	assert.NotContains(t, inputRequired, "email")

	// Verify output schema
	assert.Equal(t, "UserOutput", outputSchema.Name)
	assert.Equal(t, "Output schema for UserOutput function return value", outputSchema.Description)
	assert.Len(t, outputSchema.Fields, 3)

	// Verify output fields
	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	successField := outputProps["success"].(map[string]string)
	assert.Equal(t, "boolean", successField["type"])
	assert.Equal(t, "Whether operation succeeded", successField["description"])

	messageField := outputProps["message"].(map[string]string)
	assert.Equal(t, "string", messageField["type"])
	assert.Equal(t, "Status message", messageField["description"])

	outputUserIDField := outputProps["user_id"].(map[string]string)
	assert.Equal(t, "integer", outputUserIDField["type"])
	assert.Equal(t, "User identifier", outputUserIDField["description"])

	// Verify output required fields
	outputRequired := outputMap["required"].([]string)
	assert.Contains(t, outputRequired, "success")
	assert.NotContains(t, outputRequired, "message")
	assert.NotContains(t, outputRequired, "user_id")
}

// TestNewSchemasFromFunc_InvalidInputType tests error handling when input type is not a struct
func TestNewSchemasFromFunc_InvalidInputType(t *testing.T) {
	type ValidOutput struct {
		Result string `json:"result"`
	}

	// Handler with non-struct input (string)
	handler := func(ctx context.Context, input string) (ValidOutput, error) {
		return ValidOutput{Result: input}, nil
	}

	_, _, err := NewSchemasFromFunc(handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "input parameter type must be a struct")
}

// TestNewSchemasFromFunc_InvalidOutputType tests error handling when output type is not a struct
func TestNewSchemasFromFunc_InvalidOutputType(t *testing.T) {
	type ValidInput struct {
		Query string `json:"query"`
	}

	// Handler with non-struct output (string)
	handler := func(ctx context.Context, input ValidInput) (string, error) {
		return "result", nil
	}

	_, _, err := NewSchemasFromFunc(handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "output return type must be a struct")
}

// TestNewSchemasFromFunc_NoValidFields tests error handling when structs have no valid fields
func TestNewSchemasFromFunc_NoValidFields(t *testing.T) {
	type EmptyInput struct {
		privateField string // unexported field
	}

	type EmptyOutput struct {
		privateResult int // unexported field
	}

	handler := func(ctx context.Context, input EmptyInput) (EmptyOutput, error) {
		return EmptyOutput{}, nil
	}

	_, _, err := NewSchemasFromFunc(handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no valid fields found")
}

// TestPointerFields tests that pointer fields are properly handled in schema generation
func TestPointerFields(t *testing.T) {
	type CompactSearchResponse struct {
		CompanyName string `json:"company_name" desc:"Company name"`
		CIK         string `json:"cik" desc:"Company CIK"`
	}

	type FullSearchResponse struct {
		CompanyName string `json:"company_name" desc:"Company name"`
		CIK         string `json:"cik" desc:"Company CIK"`
		Address     string `json:"address" desc:"Company address"`
		Industry    string `json:"industry" desc:"Industry"`
	}

	type AllCompanySearchResponse struct {
		CompactSearchResponse *CompactSearchResponse `json:"compact_search_response,omitempty" desc:"Compact search result"`
		FullSearchResponse    *FullSearchResponse    `json:"full_search_response,omitempty" desc:"Full search result"`
	}

	handler := func(ctx context.Context, input SearchToolParams) (AllCompanySearchResponse, error) {
		return AllCompanySearchResponse{
			CompactSearchResponse: &CompactSearchResponse{
				CompanyName: "Test Company",
				CIK:         "0001234567",
			},
		}, nil
	}

	_, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	// Verify output schema has the pointer fields
	assert.Equal(t, "AllCompanySearchResponse", outputSchema.Name)
	assert.Len(t, outputSchema.Fields, 2)

	// Verify fields are present
	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	// Check compact_search_response field
	compactField, hasCompact := outputProps["compact_search_response"]
	assert.True(t, hasCompact, "should have compact_search_response field")
	compactFieldMap := compactField.(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(compactFieldMap["type"]))
	assert.Equal(t, "Compact search result", compactFieldMap["description"])

	// Check nested properties of CompactSearchResponse
	compactProps := compactFieldMap["properties"].(map[string]interface{})
	assert.Contains(t, compactProps, "company_name")
	assert.Contains(t, compactProps, "cik")

	// Check full_search_response field
	fullField, hasFull := outputProps["full_search_response"]
	assert.True(t, hasFull, "should have full_search_response field")
	fullFieldMap := fullField.(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(fullFieldMap["type"]))
	assert.Equal(t, "Full search result", fullFieldMap["description"])

	// Check nested properties of FullSearchResponse
	fullProps := fullFieldMap["properties"].(map[string]interface{})
	assert.Contains(t, fullProps, "company_name")
	assert.Contains(t, fullProps, "cik")
	assert.Contains(t, fullProps, "address")
	assert.Contains(t, fullProps, "industry")

	// Verify pointer fields are NOT required (they're optional by default)
	requiredFields := outputMap["required"]
	if requiredFields != nil {
		required := requiredFields.([]string)
		assert.NotContains(t, required, "compact_search_response")
		assert.NotContains(t, required, "full_search_response")
	}
}

// TestPointerPrimitives tests pointer fields with primitive types
func TestPointerPrimitives(t *testing.T) {
	type OptionalFieldsParams struct {
		RequiredName string  `json:"required_name" desc:"Required name" required:"true"`
		OptionalAge  *int    `json:"optional_age,omitempty" desc:"Optional age"`
		OptionalCity *string `json:"optional_city,omitempty" desc:"Optional city"`
		OptionalFlag *bool   `json:"optional_flag,omitempty" desc:"Optional flag"`
	}

	handler := func(ctx context.Context, params OptionalFieldsParams) (string, error) {
		return "processed", nil
	}

	schema, err := NewSchemaFromFunc(handler)
	assert.NoError(t, err)

	paramsMap := GetPropertiesMap(schema)
	properties := paramsMap["properties"].(map[string]interface{})

	// Check required field
	requiredField := properties["required_name"].(map[string]string)
	assert.Equal(t, "string", requiredField["type"])

	// Check optional pointer fields exist
	ageField := properties["optional_age"].(map[string]string)
	assert.Equal(t, "integer", ageField["type"])
	assert.Equal(t, "Optional age", ageField["description"])

	cityField := properties["optional_city"].(map[string]string)
	assert.Equal(t, "string", cityField["type"])
	assert.Equal(t, "Optional city", cityField["description"])

	flagField := properties["optional_flag"].(map[string]string)
	assert.Equal(t, "boolean", flagField["type"])
	assert.Equal(t, "Optional flag", flagField["description"])

	// Verify required fields
	required := paramsMap["required"].([]string)
	assert.Contains(t, required, "required_name")
	assert.NotContains(t, required, "optional_age")
	assert.NotContains(t, required, "optional_city")
	assert.NotContains(t, required, "optional_flag")
}

// TestPointerFieldsV2 tests pointer fields with NewSchemaFromFuncV2
func TestPointerFieldsV2(t *testing.T) {
	type ResponseData struct {
		Status  string `json:"status" desc:"Operation status"`
		Message string `json:"message" desc:"Status message"`
	}

	type OptionalResponse struct {
		Success bool          `json:"success" desc:"Whether operation succeeded" required:"true"`
		Data    *ResponseData `json:"data,omitempty" desc:"Response data if available"`
		Error   *string       `json:"error,omitempty" desc:"Error message if failed"`
	}

	handler := func(ctx context.Context, input SearchToolParams) (OptionalResponse, error) {
		return OptionalResponse{Success: true}, nil
	}

	_, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	// Check success field (required)
	successField := outputProps["success"].(map[string]string)
	assert.Equal(t, "boolean", successField["type"])

	// Check data field (pointer to struct)
	dataField := outputProps["data"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(dataField["type"]))
	assert.Equal(t, "Response data if available", dataField["description"])
	dataProps := dataField["properties"].(map[string]interface{})
	assert.Contains(t, dataProps, "status")
	assert.Contains(t, dataProps, "message")

	// Check error field (pointer to string)
	errorField := outputProps["error"].(map[string]string)
	assert.Equal(t, "string", errorField["type"])
	assert.Equal(t, "Error message if failed", errorField["description"])

	// Verify only success is required
	required := outputMap["required"].([]string)
	assert.Contains(t, required, "success")
	assert.NotContains(t, required, "data")
	assert.NotContains(t, required, "error")
}

// TestMixedPointerAndNonPointer tests structs with both pointer and non-pointer fields
func TestMixedPointerAndNonPointer(t *testing.T) {
	type MixedParams struct {
		ID          int     `json:"id" desc:"Required ID" required:"true"`
		Name        string  `json:"name" desc:"Required name" required:"true"`
		Description *string `json:"description,omitempty" desc:"Optional description"`
		Age         *int    `json:"age,omitempty" desc:"Optional age"`
	}

	handler := func(ctx context.Context, params MixedParams) (string, error) {
		return "processed", nil
	}

	schema, err := SafeSchemaFromFunc(handler)
	assert.NoError(t, err)

	properties := schema["properties"].(map[string]interface{})

	// Verify all fields exist
	assert.Contains(t, properties, "id")
	assert.Contains(t, properties, "name")
	assert.Contains(t, properties, "description")
	assert.Contains(t, properties, "age")

	// Verify required fields
	required := schema["required"].([]string)
	assert.Contains(t, required, "id")
	assert.Contains(t, required, "name")
	assert.NotContains(t, required, "description")
	assert.NotContains(t, required, "age")
}

// TestArrayOfPrimitives tests the ArrayOf function and primitive array schema generation
func TestArrayOfPrimitives(t *testing.T) {
	type ArrayTestParams struct {
		Tags        []string  `json:"tags" desc:"List of tags" required:"true"`
		IDs         []int     `json:"ids" desc:"List of IDs"`
		Scores      []float64 `json:"scores" desc:"List of scores"`
		Flags       []bool    `json:"flags" desc:"List of boolean flags"`
		Description string    `json:"description" desc:"Description text"`
	}

	handler := func(ctx context.Context, params ArrayTestParams) (string, error) {
		return "processed", nil
	}

	schema, err := SafeSchemaFromFunc(handler)
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	properties, ok := schema["properties"].(map[string]interface{})
	assert.True(t, ok, "properties should be a map")

	// Test string array field
	tagsField, ok := properties["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Equal(t, "array", fmt.Sprint(tagsField["type"]))
	assert.Equal(t, "List of tags", tagsField["description"])
	tagsItems, ok := tagsField["items"].(map[string]interface{})
	assert.True(t, ok, "tags items should be a map")
	assert.Equal(t, "string", fmt.Sprint(tagsItems["type"]))

	// Test integer array field
	idsField, ok := properties["ids"].(map[string]interface{})
	assert.True(t, ok, "ids should be a map")
	assert.Equal(t, "array", fmt.Sprint(idsField["type"]))
	assert.Equal(t, "List of IDs", idsField["description"])
	idsItems, ok := idsField["items"].(map[string]interface{})
	assert.True(t, ok, "ids items should be a map")
	assert.Equal(t, "integer", fmt.Sprint(idsItems["type"]))

	// Test float array field
	scoresField, ok := properties["scores"].(map[string]interface{})
	assert.True(t, ok, "scores should be a map")
	assert.Equal(t, "array", fmt.Sprint(scoresField["type"]))
	assert.Equal(t, "List of scores", scoresField["description"])
	scoresItems, ok := scoresField["items"].(map[string]interface{})
	assert.True(t, ok, "scores items should be a map")
	assert.Equal(t, "number", fmt.Sprint(scoresItems["type"]))

	// Test boolean array field
	flagsField, ok := properties["flags"].(map[string]interface{})
	assert.True(t, ok, "flags should be a map")
	assert.Equal(t, "array", fmt.Sprint(flagsField["type"]))
	assert.Equal(t, "List of boolean flags", flagsField["description"])
	flagsItems, ok := flagsField["items"].(map[string]interface{})
	assert.True(t, ok, "flags items should be a map")
	assert.Equal(t, "boolean", fmt.Sprint(flagsItems["type"]))

	// Verify required fields
	required, ok := schema["required"].([]string)
	assert.True(t, ok)
	assert.Contains(t, required, "tags", "tags should be required")
	assert.NotContains(t, required, "ids", "ids should not be required")
}

// TestArrayOfPrimitivesWithNewSchemaFromFunc tests primitive arrays work with NewSchemaFromFunc
func TestArrayOfPrimitivesWithNewSchemaFromFunc(t *testing.T) {
	type FilterParams struct {
		Keywords []string `json:"keywords" desc:"Search keywords" required:"true"`
		Years    []int    `json:"years" desc:"Filter by years"`
	}

	handler := func(ctx context.Context, params FilterParams) (string, error) {
		return fmt.Sprintf("Filtering with %d keywords", len(params.Keywords)), nil
	}

	schema, err := NewSchemaFromFunc(handler)
	assert.NoError(t, err)

	paramsMap := GetPropertiesMap(schema)
	properties := paramsMap["properties"].(map[string]interface{})

	// Verify keywords array
	keywordsField, ok := properties["keywords"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "array", fmt.Sprint(keywordsField["type"]))
	keywordsItems := keywordsField["items"].(map[string]interface{})
	assert.Equal(t, "string", fmt.Sprint(keywordsItems["type"]))

	// Verify years array
	yearsField, ok := properties["years"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "array", fmt.Sprint(yearsField["type"]))
	yearsItems := yearsField["items"].(map[string]interface{})
	assert.Equal(t, "integer", fmt.Sprint(yearsItems["type"]))
}

// TestArrayOfPrimitivesWithNewSchemasFromFunc tests primitive arrays with input/output schemas
func TestArrayOfPrimitivesWithNewSchemasFromFunc(t *testing.T) {
	type BatchInput struct {
		Operations []string `json:"operations" desc:"List of operations to perform" required:"true"`
		Targets    []int    `json:"targets" desc:"Target IDs"`
	}

	type BatchOutput struct {
		Completed  []bool   `json:"completed" desc:"Completion status for each operation" required:"true"`
		Errors     []string `json:"errors" desc:"Error messages if any"`
		Percentages []float64 `json:"percentages" desc:"Success percentages"`
	}

	handler := func(ctx context.Context, input BatchInput) (BatchOutput, error) {
		return BatchOutput{
			Completed: []bool{true, true, false},
			Errors:    []string{},
			Percentages: []float64{100.0, 100.0, 0.0},
		}, nil
	}

	inputSchema, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	// Test input schema
	inputMap := GetPropertiesMap(inputSchema)
	inputProps := inputMap["properties"].(map[string]interface{})

	operationsField := inputProps["operations"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(operationsField["type"]))
	operationsItems := operationsField["items"].(map[string]interface{})
	assert.Equal(t, "string", fmt.Sprint(operationsItems["type"]))

	targetsField := inputProps["targets"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(targetsField["type"]))
	targetsItems := targetsField["items"].(map[string]interface{})
	assert.Equal(t, "integer", fmt.Sprint(targetsItems["type"]))

	// Test output schema
	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	completedField := outputProps["completed"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(completedField["type"]))
	completedItems := completedField["items"].(map[string]interface{})
	assert.Equal(t, "boolean", fmt.Sprint(completedItems["type"]))

	errorsField := outputProps["errors"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(errorsField["type"]))
	errorsItems := errorsField["items"].(map[string]interface{})
	assert.Equal(t, "string", fmt.Sprint(errorsItems["type"]))

	percentagesField := outputProps["percentages"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(percentagesField["type"]))
	percentagesItems := percentagesField["items"].(map[string]interface{})
	assert.Equal(t, "number", fmt.Sprint(percentagesItems["type"]))
}

// TestPointerReturnType tests that functions returning pointer structs work correctly
func TestPointerReturnType(t *testing.T) {
	type SearchParams struct {
		Query string `json:"query" desc:"Search query" required:"true"`
	}

	type SearchResult struct {
		Found bool   `json:"found" desc:"Whether results were found"`
		Count int    `json:"count" desc:"Number of results"`
		Items []string `json:"items" desc:"Result items"`
	}

	// Function that returns a pointer to struct
	handler := func(ctx context.Context, params SearchParams) (*SearchResult, error) {
		return &SearchResult{
			Found: true,
			Count: 5,
			Items: []string{"item1", "item2"},
		}, nil
	}

	inputSchema, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	// Verify input schema
	assert.Equal(t, "SearchParams", inputSchema.Name)
	assert.Len(t, inputSchema.Fields, 1)

	// Verify output schema - should work with pointer return type
	assert.Equal(t, "SearchResult", outputSchema.Name)
	assert.Len(t, outputSchema.Fields, 3)

	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	foundField := outputProps["found"].(map[string]string)
	assert.Equal(t, "boolean", foundField["type"])

	countField := outputProps["count"].(map[string]string)
	assert.Equal(t, "integer", countField["type"])

	itemsField := outputProps["items"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(itemsField["type"]))
}

// TestMapFields tests map field support
func TestMapFields(t *testing.T) {
	type StatsResponse struct {
		Scores      map[string]int     `json:"scores" desc:"Score by category"`
		Metadata    map[string]string  `json:"metadata" desc:"Additional metadata"`
		Flags       map[string]bool    `json:"flags" desc:"Feature flags"`
		Percentages map[string]float64 `json:"percentages" desc:"Percentage values"`
	}

	handler := func(ctx context.Context, params SearchToolParams) (StatsResponse, error) {
		return StatsResponse{
			Scores: map[string]int{"math": 95, "english": 87},
		}, nil
	}

	_, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	// Check scores map (map[string]int)
	scoresField := outputProps["scores"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(scoresField["type"]))
	assert.Equal(t, "Score by category", scoresField["description"])
	scoresAdditional := scoresField["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "integer", scoresAdditional["type"])

	// Check metadata map (map[string]string)
	metadataField := outputProps["metadata"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(metadataField["type"]))
	metadataAdditional := metadataField["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "string", metadataAdditional["type"])

	// Check flags map (map[string]bool)
	flagsField := outputProps["flags"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(flagsField["type"]))
	flagsAdditional := flagsField["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "boolean", flagsAdditional["type"])

	// Check percentages map (map[string]float64)
	percentagesField := outputProps["percentages"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(percentagesField["type"]))
	percentagesAdditional := percentagesField["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "number", percentagesAdditional["type"])
}

// TestMapWithStructValues tests maps with struct values
func TestMapWithStructValues(t *testing.T) {
	type UserInfo struct {
		Name  string `json:"name" desc:"User name"`
		Age   int    `json:"age" desc:"User age"`
		Email string `json:"email" desc:"User email"`
	}

	type UsersResponse struct {
		Users map[string]UserInfo `json:"users" desc:"Users by ID"`
	}

	handler := func(ctx context.Context, params SearchToolParams) (UsersResponse, error) {
		return UsersResponse{
			Users: map[string]UserInfo{
				"user1": {Name: "Alice", Age: 30, Email: "alice@example.com"},
			},
		}, nil
	}

	_, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	// Check users map (map[string]UserInfo)
	usersField := outputProps["users"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(usersField["type"]))
	assert.Equal(t, "Users by ID", usersField["description"])

	usersAdditional := usersField["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "object", usersAdditional["type"])

	// Check nested struct properties
	userProps := usersAdditional["properties"].(map[string]interface{})
	assert.Contains(t, userProps, "name")
	assert.Contains(t, userProps, "age")
	assert.Contains(t, userProps, "email")

	nameField := userProps["name"].(map[string]string)
	assert.Equal(t, "string", nameField["type"])

	ageField := userProps["age"].(map[string]string)
	assert.Equal(t, "integer", ageField["type"])
}

// TestPointerParamsAndPointerReturn tests both pointer params and pointer return
func TestPointerParamsAndPointerReturn(t *testing.T) {
	type CompactSearchResponse struct {
		CompanyName string `json:"company_name" desc:"Company name"`
		CIK         string `json:"cik" desc:"Company CIK"`
	}

	type FullSearchResponse struct {
		CompanyName string `json:"company_name" desc:"Company name"`
		CIK         string `json:"cik" desc:"Company CIK"`
		Address     string `json:"address" desc:"Company address"`
	}

	type AllCompanySearchResponse struct {
		CompactSearchResponse *CompactSearchResponse `json:"compact_search_response,omitempty" desc:"Compact result"`
		FullSearchResponse    *FullSearchResponse    `json:"full_search_response,omitempty" desc:"Full result"`
	}

	type AllCompanySearchToolParams struct {
		CompanyName string `json:"company_name" desc:"Company name to search"`
		CompactMode bool   `json:"compact_mode" desc:"Use compact mode"`
	}

	// This is the real-world signature from the user's issue
	handler := func(ctx context.Context, params *AllCompanySearchToolParams) (*AllCompanySearchResponse, error) {
		return &AllCompanySearchResponse{
			CompactSearchResponse: &CompactSearchResponse{
				CompanyName: "Test Company",
				CIK:         "0001234567",
			},
		}, nil
	}

	inputSchema, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	// Verify input schema
	assert.Equal(t, "AllCompanySearchToolParams", inputSchema.Name)
	assert.Len(t, inputSchema.Fields, 2)

	// Verify output schema
	assert.Equal(t, "AllCompanySearchResponse", outputSchema.Name)
	assert.Len(t, outputSchema.Fields, 2)

	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	// Both pointer fields should be present
	assert.Contains(t, outputProps, "compact_search_response")
	assert.Contains(t, outputProps, "full_search_response")
}

// TestNestedMapInStruct tests maps inside nested structs
func TestNestedMapInStruct(t *testing.T) {
	type Summary struct {
		TotalCount      int            `json:"total_count" desc:"Total count"`
		CategoriesFound map[string]int `json:"categories_found" desc:"Categories found"`
	}

	type Metadata struct {
		SearchQuery string `json:"search_query" desc:"Search query"`
		ExecutedAt  string `json:"executed_at" desc:"Execution time"`
	}

	type SearchResponse struct {
		Summary  Summary  `json:"summary" desc:"Search summary"`
		Metadata Metadata `json:"metadata" desc:"Search metadata"`
	}

	handler := func(ctx context.Context, params SearchToolParams) (SearchResponse, error) {
		return SearchResponse{}, nil
	}

	_, outputSchema, err := NewSchemasFromFunc(handler)
	assert.NoError(t, err)

	outputMap := GetPropertiesMap(outputSchema)
	outputProps := outputMap["properties"].(map[string]interface{})

	// Check Summary nested object
	summaryField := outputProps["summary"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(summaryField["type"]))

	summaryProps := summaryField["properties"].(map[string]interface{})

	// Check CategoriesFound map inside Summary
	categoriesField := summaryProps["categories_found"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(categoriesField["type"]))
	categoriesAdditional := categoriesField["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "integer", categoriesAdditional["type"])
}

// TestSchemasFromFunc_PointerReturnWithMaps tests the exact scenario from the user's question
func TestSchemasFromFunc_PointerReturnWithMaps(t *testing.T) {
	// Minimal test struct with map fields
	type TestResponse struct {
		Name     string            `json:"name" desc:"Name field"`
		Tags     map[string]string `json:"tags" desc:"Tags map"`
		Counts   map[string]int    `json:"counts" desc:"Counts map"`
		Metadata map[string]interface{} `json:"metadata,omitempty" desc:"Metadata map"`
		Categories struct {
			Found map[string]int `json:"found" desc:"Categories found"`
		} `json:"categories" desc:"Categories"`
	}

	type TestParams struct {
		Query string `json:"query" desc:"Search query"`
	}

	// Function that returns pointer to struct (issue #1) with maps (issue #2)
	testFunction := func(ctx context.Context, params *TestParams) (*TestResponse, error) {
		return &TestResponse{
			Name:   "test",
			Tags:   map[string]string{"foo": "bar"},
			Counts: map[string]int{"total": 42},
		}, nil
	}

	schemaIn, schemaOut, err := SafeSchemasFromFunc(testFunction)

	if err != nil {
		t.Fatalf("Failed to generate schema: %v", err)
	}

	if schemaIn == nil {
		t.Fatal("Input schema should not be nil")
	}

	if schemaOut == nil {
		t.Fatal("Output schema should not be nil")
	}

	// Verify map fields are present in output schema
	properties := schemaOut["properties"].(map[string]interface{})

	// Check Name field (regular string)
	nameSchema := properties["name"].(map[string]string)
	assert.Equal(t, "string", nameSchema["type"])

	// Check Tags map[string]string
	tagsSchema := properties["tags"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(tagsSchema["type"]))
	assert.Equal(t, "Tags map", tagsSchema["description"])
	tagsAdditional := tagsSchema["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "string", tagsAdditional["type"])

	// Check Counts map[string]int
	countsSchema := properties["counts"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(countsSchema["type"]))
	assert.Equal(t, "Counts map", countsSchema["description"])
	countsAdditional := countsSchema["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "integer", countsAdditional["type"])

	// Check Metadata map[string]interface{} - should allow any value type
	metadataSchema := properties["metadata"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(metadataSchema["type"]))
	assert.Equal(t, "Metadata map", metadataSchema["description"])
	// For interface{} maps, additionalProperties should be true (allows any schema)
	assert.Equal(t, true, metadataSchema["additionalProperties"])

	// Check Categories nested struct with map inside
	categoriesSchema := properties["categories"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(categoriesSchema["type"]))
	categoriesProps := categoriesSchema["properties"].(map[string]interface{})

	// Check Found map[string]int inside Categories
	foundSchema := categoriesProps["found"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(foundSchema["type"]))
	assert.Equal(t, "Categories found", foundSchema["description"])
	foundAdditional := foundSchema["additionalProperties"].(map[string]interface{})
	assert.Equal(t, "integer", foundAdditional["type"])

	// Verify input schema
	inputProps := schemaIn["properties"].(map[string]interface{})
	querySchema := inputProps["query"].(map[string]string)
	assert.Equal(t, "string", querySchema["type"])
	assert.Equal(t, "Search query", querySchema["description"])
}

// TestMixedArrayTypes tests structs with both primitive arrays and object arrays
func TestMixedArrayTypes(t *testing.T) {
	type SubItem struct {
		Name  string `json:"name" desc:"Item name"`
		Value int    `json:"value" desc:"Item value"`
	}

	type MixedParams struct {
		Tags   []string  `json:"tags" desc:"Simple string tags"`
		Items  []SubItem `json:"items" desc:"Complex item objects"`
		Scores []float64 `json:"scores" desc:"Numeric scores"`
	}

	handler := func(ctx context.Context, params MixedParams) (string, error) {
		return "processed", nil
	}

	schema, err := SafeSchemaFromFunc(handler)
	assert.NoError(t, err)

	properties := schema["properties"].(map[string]interface{})

	// Test primitive array (tags)
	tagsField := properties["tags"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(tagsField["type"]))
	tagsItems := tagsField["items"].(map[string]interface{})
	assert.Equal(t, "string", fmt.Sprint(tagsItems["type"]))

	// Test object array (items)
	itemsField := properties["items"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(itemsField["type"]))
	itemsItems := itemsField["items"].(map[string]interface{})
	assert.Equal(t, "object", fmt.Sprint(itemsItems["type"]))
	itemsProps := itemsItems["properties"].(map[string]interface{})
	assert.Contains(t, itemsProps, "name")
	assert.Contains(t, itemsProps, "value")

	// Test primitive array (scores)
	scoresField := properties["scores"].(map[string]interface{})
	assert.Equal(t, "array", fmt.Sprint(scoresField["type"]))
	scoresItems := scoresField["items"].(map[string]interface{})
	assert.Equal(t, "number", fmt.Sprint(scoresItems["type"]))
}
