package test

import (
	"context"
	"fmt"
	"github.com/mhpenta/jobj/funcschema"
	"log/slog"
	"testing"
)

type SearchTool struct{}

type SearchToolParams struct {
	ID    int
	Query string
}

func (f *SearchTool) Parameters() map[string]interface{} {
	schema, err := funcschema.NewSchemaFromFunc(f.SearchForData)
	if err != nil {
		slog.Error("Failed to parse function schema", "error", err)
		return nil
	}

	return funcschema.GetPropertiesMap(schema) // Or whatever subset might be useful
}

func (f *SearchTool) SearchForData(ctx context.Context, params *SearchToolParams) (string, error) {
	return fmt.Sprintf("Searching for %s", params.Query), nil
}

func TestSearchTool_Parameters(t *testing.T) {
	searchTool := &SearchTool{}
	paramsMap := searchTool.Parameters()

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
