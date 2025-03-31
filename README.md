# jobj

`jobj` is a Go package that provides a flexible way to create and manage JSON schema definitions. It offers an API for defining JSON data structures for generating JSON schemas specifically in the context of requesting structured JSON objects from large language models (LLMs). It also provides limited support for XML schemas

## Why `jobj`?

If you want to create a json schema from a struct, you can use packages like [jsonschema](https://github.com/invopop/jsonschema) that use reflection and field tags. However, using tags creates ergonomic challenges for the programmer, specifically if you ever require a long description on a field. `jobj` is designed to get around this ergonomic issue and to create json object tooling for LLM development, which requires only a subset of the json schema. 

## JSON Schema Specification Coverage

This package implements a focused subset of the [JSON Schema Draft-07](https://json-schema.org/specification-links.html#draft-7) specification, prioritizing the elements most useful for LLM interactions:

### Implemented
- Core schema structure with `$schema`, `definitions`, and `$ref`
- Object type definitions with properties and typing
- Required/optional field specification
- Property descriptions
- Basic property types: string, number, integer, boolean
- Nested objects and arrays
- Enum-like constraints using `anyOf`
- Default value specification
- `additionalProperties` field control

### Not Implemented
- Format validation (except for custom `JsonDateTime` type)
- Regular expression pattern validation
- Numeric constraints (minimum, maximum, etc.)
- String constraints (minLength, maxLength, etc.)
- Array constraints (minItems, maxItems, etc.)
- Schema composition (allOf, oneOf, not)
- External schema references
- Conditional schemas (if-then-else)

## Features

- Clean API for schema definition with chainable methods
- Custom JSON date/time handling with `JsonDateTime` type supporting multiple formats
- Support for nested objects and arrays
- Enum-like field constraints using `AnyOf`
- Required/optional field specification
- XML schema generation (limited subset)
- Generate JSON schemas from Go function signatures via the `funcschema` subpackage
- Validation against Go structs
- Struct tag parsing for automated schema generation

## Usage And Examples

```bash
go get github.com/mhpenta/jobj
```

### Basic Schema Definition

```go
package main

import "github.com/mhpenta/jobj"

type HeadlinesResponse struct {
    jobj.Schema
}

func NewHeadlineResponse() *HeadlinesResponse {
	h := &HeadlinesResponse{}
	h.Name = "HeadlinesResponse"
	h.Description = "Response schema for press release headline extraction"
	h.Fields = []*jobj.Field{
		jobj.Text("headline").
			Desc("The exact headline from the press release (in proper case)").
			Required(),
		jobj.Text("headline_without_company_name").
			Desc("The headline from the press release modified to remove the company name (in proper case)").
			Required(),
		jobj.Float("confidence").
			Desc("Confidence in the headlines extracted").
			Required(),
	}
	return h
}

// Generate JSON schema
schema := NewHeadlineResponse()
schemaJSON := schema.GetSchemaString()
fmt.Println(schemaJSON)
```

Prints:
```json
{
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
}

```


### Field Types

The package supports various field types:

- `Text(name string)` - String fields
- `Bool(name string)` - Boolean fields
- `Float(name string)` - Floating-point numbers
- `Int(name string)` - Integer fields
- `Date(name string)` - Date fields using the custom JsonDateTime type
- `Array(name string, fields []*Field)` - Array of objects
- `Object(name string, fields []*Field)` - Nested object structures
- `AnyOf(name string, enums []ConstDescription)` - Enumerated values

### Field Modifiers

Fields can be customized using chainable modifiers:

```go
jobj.Text("field").
    Desc("Field description").  // Add description
    Required().                 // Mark as required (adds field name to schema's "required" array)
    Optional().                // Mark as optional (removes field from "required" array)
    Type("custom_type").       // Set custom type
    SetValue("default")        // Set default value
```

### Working with JsonDateTime

The package includes a custom `JsonDateTime` type for handling dates:

```go
type Document struct {
    PublishDate JsonDateTime `json:"publish_date"`
}
```

The `JsonDateTime` type offers robust date parsing that handles multiple formats, including:
- ISO 8601/RFC3339 (2006-01-02T15:04:05Z)
- RFC3339 with timezone (2006-01-02T15:04:05-07:00)
- Simple date format (2006-01-02)
- Various other common time formats

### The funcschema Subpackage

The `funcschema` subpackage allows you to automatically generate JSON schemas from Go function signatures. This is particularly useful for creating LLM tools that require parameter schemas.

```go
// Define your tool's parameter struct
type SearchToolParams struct {
    ID    int    `desc:"ID of item to search" required:"true" `
    Query string `desc:"Query to search for, e.g., xyz" required:"true"`
}

// Create a function that uses the parameters
func (t *SearchTool) ExecuteSearch(ctx context.Context, params *SearchToolParams) (*ToolResult, error) {
    // Implementation...
}

// Generate a JSON schema from the function
func (t *SearchTool) Parameters() map[string]interface{} {
    schema, err := funcschema.SafeSchemaFromFunc(t.ExecuteSearch)
    if err != nil {
        // Handle error
    }
    return schema
}

// Called by the agent
func (t *SearchTool) Execute(ctx context.Context, params json.RawMessage) (*tools.ToolResult, error) {
    var paramsStruct *SearchToolParams
    if err := json.Unmarshal(params, paramsStruct); err != nil {
        // Handle error
    }
    
    searchResult, err := t.ExecuteSearch(ctx, paramsStruct)
    if err != nil {
        // Handle error
    }
	return searchResult, nil
}
```

The `funcschema` subpackage's schema for tool use looks like: 

```json
{
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
}
```

The `funcschema` subpackage offers several options:
- `SchemaFromStruct[T]()` - Generate schema directly from a struct type
- `NewSchemaFromFuncV2()` - Type-safe schema generation with generics
- `NewSchemaFromFunc()` - Non-generic version for compatibility
- `GetPropertiesMap()` - Convert schema to a properties map for LLM tool definitions

### XML Schema Support

While primarily focused on JSON Schema, `jobj` provides limited XML Schema generation via the `GetXMLSchemaString()` method:

```go
schema := NewHeadlineResponse()
schema.UseXML = true
xmlSchema := schema.GetXMLSchemaString()
```

XML support is limited to basic type mapping and does not implement the full XML Schema specification.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

Copyright (c) 2025, github.com/mhpenta

All rights reserved.