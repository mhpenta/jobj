package jobj

import (
	"testing"
)

// TestResponseWithMatchingFields tests a valid case where schema fields match struct fields
func TestResponseWithMatchingFields(t *testing.T) {
	// Create a response struct with matching schema fields
	resp := NewValidTestResponse()

	// Validate the schema against the struct
	err := resp.Validate(resp)
	if err != nil {
		t.Errorf("Validate failed for valid struct: %v", err)
	}
}

// TestResponseWithMissingFieldInSchema tests when a struct field is missing from the schema
func TestResponseWithMissingFieldInSchema(t *testing.T) {
	// Create response with a missing field in schema
	resp := NewMissingSchemaFieldTestResponse()

	// Validation should fail because the required struct field is missing from schema
	err := resp.Validate(resp)
	if err == nil {
		t.Error("Validate should have failed for missing schema field")
	}
}

// TestResponseWithMissingFieldInStruct tests when a schema field is missing from the struct
func TestResponseWithMissingFieldInStruct(t *testing.T) {
	// Create response with a missing field in struct
	resp := NewMissingStructFieldTestResponse()

	// Validation should fail because schema field doesn't exist in struct
	err := resp.Validate(resp)
	if err == nil {
		t.Error("Validate should have failed for missing struct field")
	}
}

// TestResponseWithTypeMismatch tests type compatibility checking
func TestTypeMismatchValidation(t *testing.T) {
	// Create response with a type mismatch
	resp := NewTypeMismatchTestResponse()

	// Validation should fail because of type mismatch
	err := resp.Validate(resp)
	if err == nil {
		t.Error("Validate should have failed for type mismatch")
	}
}

// TestResponseWithOmitemptyMismatch tests required field vs omitempty tag checking
func TestOmitemptyMismatchValidation(t *testing.T) {
	// Create response with a mismatch between schema required field and struct omitempty tag
	resp := NewOmitemptyMismatchTestResponse()

	// Validation should fail because required field has omitempty tag
	err := resp.Validate(resp)
	if err == nil {
		t.Error("Validate should have failed for omitempty mismatch")
	}
}

// TestResponseWithEmbeddedFields tests validation with embedded fields
func TestEmbeddedFieldsValidation(t *testing.T) {
	// Create response with embedded fields
	resp := NewEmbeddedFieldsTestResponse()

	// Validation should succeed for properly configured embedded fields
	err := resp.Validate(resp)
	if err != nil {
		t.Errorf("Validate failed for valid embedded fields: %v", err)
	}
}

// Test structures and constructors

// ValidTestResponse represents a valid response struct
type ValidTestResponse struct {
	Schema
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Count       int     `json:"count"`
	IsValid     bool    `json:"is_valid"`
	Score       float64 `json:"score"`
}

func NewValidTestResponse() *ValidTestResponse {
	r := &ValidTestResponse{}
	r.Name = "ValidTestResponse"
	r.Description = "A test response object"
	r.Fields = []*Field{
		Text("name").Desc("The name field").Required(),
		Text("description").Desc("The description field").Required(),
		Int("count").Desc("The count field").Required(),
		Bool("is_valid").Desc("The validity field").Required(),
		Float("score").Desc("The score field").Required(),
	}
	return r
}

// MissingSchemaFieldTestResponse has a field in struct but not in schema
type MissingSchemaFieldTestResponse struct {
	Schema
	Name  string `json:"name"`
	Count int    `json:"count"` // This field is not in the schema
}

func NewMissingSchemaFieldTestResponse() *MissingSchemaFieldTestResponse {
	r := &MissingSchemaFieldTestResponse{}
	r.Name = "MissingSchemaFieldTestResponse"
	r.Description = "A test response with missing schema field"
	r.Fields = []*Field{
		Text("name").Desc("The name field").Required(),
		// count field is missing from schema but exists in struct
	}
	return r
}

// MissingStructFieldTestResponse has a field in schema but not in struct
type MissingStructFieldTestResponse struct {
	Schema
	Name string `json:"name"`
	// count field is missing from struct but exists in schema
}

func NewMissingStructFieldTestResponse() *MissingStructFieldTestResponse {
	r := &MissingStructFieldTestResponse{}
	r.Name = "MissingStructFieldTestResponse"
	r.Description = "A test response with missing struct field"
	r.Fields = []*Field{
		Text("name").Desc("The name field").Required(),
		Int("count").Desc("The count field").Required(), // This field doesn't exist in struct
	}
	return r
}

// TypeMismatchTestResponse has a type mismatch between schema and struct
type TypeMismatchTestResponse struct {
	Schema
	Name  string `json:"name"`
	Count string `json:"count"` // This is a string but schema defines int
}

func NewTypeMismatchTestResponse() *TypeMismatchTestResponse {
	r := &TypeMismatchTestResponse{}
	r.Name = "TypeMismatchTestResponse"
	r.Description = "A test response with type mismatch"
	r.Fields = []*Field{
		Text("name").Desc("The name field").Required(),
		Int("count").Desc("The count field").Required(), // Defined as int but struct has string
	}
	return r
}

// OmitemptyMismatchTestResponse tests omitempty tag handling
type OmitemptyMismatchTestResponse struct {
	Schema
	Name  string `json:"name"`
	Count int    `json:"count,omitempty"` // Has omitempty but schema defines as required
}

func NewOmitemptyMismatchTestResponse() *OmitemptyMismatchTestResponse {
	r := &OmitemptyMismatchTestResponse{}
	r.Name = "OmitemptyMismatchTestResponse"
	r.Description = "A test response with omitempty mismatch"
	r.Fields = []*Field{
		Text("name").Desc("The name field").Required(),
		Int("count").Desc("The count field").Required(), // Defined as required but struct has omitempty
	}
	return r
}

// Embedded types for testing embedded fields
type EmbeddedData struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// EmbeddedFieldsTestResponse tests validation with embedded structs
type EmbeddedFieldsTestResponse struct {
	Schema
	Name string       `json:"name"`
	Data EmbeddedData `json:"data"`
}

func NewEmbeddedFieldsTestResponse() *EmbeddedFieldsTestResponse {
	r := &EmbeddedFieldsTestResponse{}
	r.Name = "EmbeddedFieldsTestResponse"
	r.Description = "A test response with embedded fields"

	dataFields := []*Field{
		Text("value").Desc("The value field").Required(),
		Int("count").Desc("The count field").Required(),
	}

	r.Fields = []*Field{
		Text("name").Desc("The name field").Required(),
		Object("data", dataFields).Desc("The embedded data object").Required(),
	}
	return r
}
