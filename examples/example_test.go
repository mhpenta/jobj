package examples

import (
	"fmt"
	"github.com/mhpenta/jobj"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadlinesResponse_CreateFields(t *testing.T) {
	h := NewHeadlineResponse()

	expectedFields := []*jobj.Field{
		jobj.Text("headline").
			Desc("The exact headline from the press release (in proper case)").Required(),
		jobj.Text("headline_without_company_name").Required().
			Desc("The headline from the press release modified to remove the company name (in proper case)"),
		jobj.Float("confidence").
			Desc("Confidence in the headlines extracted").Required(),
	}
	fmt.Println(h.GetSchemaString())

	correctSchema := `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "definitions": {
    "HeadlinesResponse": {
      "additionalProperties": false,
      "properties": {
        "confidence": {
          "description": "Confidence in the headlines extracted",
          "type": "number"
        },
        "headline": {
          "description": "The exact headline from the press release (in proper case)",
          "type": "string"
        },
        "headline_without_company_name": {
          "description": "The headline from the press release modified to remove the company name (in proper case)",
          "type": "string"
        }
      },
      "required": [
        "headline",
        "headline_without_company_name",
        "confidence"
      ],
      "type": "object"
    }
  },
  "$ref": "#/definitions/HeadlinesResponse"
}`

	assert.Equal(t, expectedFields, h.Fields)
	assert.Equal(t, correctSchema, h.GetSchemaString())
}
