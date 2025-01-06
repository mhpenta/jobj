package examples

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTranscriptCorrectionsResponse(t *testing.T) {
	schema := NewTranscriptCorrectionsResponse()

	fmt.Println(schema.GetSchemaString())

	correct := `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "definitions": {
    "TranscriptCorrectionsResponse": {
      "additionalProperties": false,
      "properties": {
        "corrections": {
          "additionalProperties": false,
          "description": "",
          "items": {
            "properties": {
              "correction": {
                "description": "Details of the correction to be made to the transcript",
                "properties": {
                  "correction_type": {
                    "anyOf": [
                      {
                        "const": "merge",
                        "description": "merge two paragraphs into one with a single speaker"
                      },
                      {
                        "const": "speaker_correction",
                        "description": "correct the speaker number"
                      }
                    ],
                    "description": "Type of correction to be made"
                  },
                  "new_speaker_number": {
                    "description": "New speaker number to be assigned to the paragraph",
                    "type": "integer"
                  },
                  "paragraph_for_new_speaker_number": {
                    "description": "Paragraph number to be assigned to the new speaker number",
                    "type": "integer"
                  }
                },
                "required": [
                  "correction_type",
                  "new_speaker_number",
                  "paragraph_for_new_speaker_number"
                ],
                "type": "object"
              }
            },
            "required": [
              "correction"
            ],
            "type": "object"
          },
          "type": "array"
        }
      },
      "required": [
        "corrections"
      ],
      "type": "object"
    }
  },
  "$ref": "#/definitions/TranscriptCorrectionsResponse"
}`

	assert.Equal(t, "TranscriptCorrectionsResponse", schema.Name)
	assert.Equal(t, schema.GetSchemaString(), correct)
}
