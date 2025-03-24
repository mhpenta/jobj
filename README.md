# jobj

`jobj` is a Go package that provides a flexible way to create and manage JSON schema definitions. It offers an API for defining json data structures for the purpose of generating json schemas in context of requesting specific json objects from large language models. Works in a very limited capacity with XML schemas.

Does not attempt to implement the entire spec, only a subset for creating dynamic json schemas. 

## Features

- API for schema definition
- JSON to Go date type handling with `JsonDate`
- Support for nested objects and arrays
- Enum-like field constraints using `AnyOf`
- Required/optional field specification
- Generate JSON schemas from Go function signatures with the `funcschema` subpackage

## Usage And Examples

```bash
go get github.com/mhpenta/jobj
```

```go
package main

import "github.com/mhpenta/jobj"

type HeadlinesResponse struct {
    jobj.Response
}

func NewHeadlineResponse() *HeadlinesResponse {
    h := &HeadlinesResponse{}
    h.Name = "HeadlinesResponse"
    h.Description = "Response schema for press release headline extraction"
    h.Fields = []*jobj.Field{
        jobj.Text("headline").
            Desc("The exact headline from the press release").
            Required(),
        jobj.Float("confidence").
            Desc("Confidence score for the extraction").
            Required(),
    }
    return h
}
```

### Field Types

The package supports various field types:

- `Text(name string)` - String fields
- `Bool(name string)` - Boolean fields
- `Float(name string)` - Floating-point numbers
- `Int(name string)` - Integer fields
- `Date(name string)` - Date fields using the custom JsonDate type
- `Array(name string, fields []*Field)` - Array of objects
- `Object(name string, fields []*Field)` - Nested object structures
- `AnyOf(name string, enums []ConstDescription)` - Enumerated values

### Field Modifiers

Fields can be customized using chainable modifiers:

```go
jobj.Text("field").
    Desc("Field description").  // Add description
    Required().                 // Mark as required
    Optional().                // Mark as optional
    Type("custom_type").       // Set custom type
    SetValue("default")        // Set default value
```

### Schema Generation

Generate JSON schema:

```go
response := NewHeadlineResponse()
schemaJSON := response.GetSchemaString()
```

### Working with JsonDate

The package includes a custom `JsonDate` type for handling dates in YYYY-MM-DD format:

```go
type Document struct {
    PublishDate JsonDate `json:"publish_date"`
}
```

### The funcschema Subpackage

The `funcschema` subpackage allows you to automatically generate JSON schemas from Go function signatures. This is particularly useful for creating LLM tools that require parameter schemas.

```go
// Define your tool's parameter struct
type SearchParams struct {
    Query       string `json:"query" desc:"The search query string" required:"true"`
    ContentType string `json:"content_type,omitempty" desc:"Optional filter by content type"`
    Limit       int    `json:"limit,omitempty" desc:"Maximum number of results to return"`
}

// Create a function that uses the parameters
func (t *SearchTool) ExecuteSearch(ctx context.Context, params SearchParams) (*ToolResult, error) {
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
```

This approach makes it easy to maintain type safety while automatically generating parameter schemas for LLM tools, eliminating the need to manually define and keep schemas in sync with your code.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

BSD 3-Clause License

Copyright (c) 2025, github.com/mhpenta 

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.
Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.
Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.