package funcschema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserInfo struct {
	Name   string `desc:"User's full name" required:"true"`
	Age    int    `desc:"User's age in years"`
	Active bool   `desc:"Whether the user is active"`
}

func TestSchemaFromStruct(t *testing.T) {
	schema, err := SchemaFromStruct[UserInfo]()
	if err != nil {
		t.Fatalf("Failed to generate schema: %v", err)
	}

	expected := `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "definitions": {
    "UserInfo": {
      "additionalProperties": false,
      "properties": {
        "Active": {
          "description": "Whether the user is active",
          "type": "boolean"
        },
        "Age": {
          "description": "User's age in years",
          "type": "integer"
        },
        "Name": {
          "description": "User's full name",
          "type": "string"
        }
      },
      "required": [
        "Name"
      ],
      "type": "object"
    }
  },
  "$ref": "#/definitions/UserInfo"
}`

	// Verify schema output matches expected
	assert.Equal(t, expected, schema.GetSchemaString())

	// Test error case with non-struct type
	_, err = SchemaFromStruct[int]()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected struct type")
}
