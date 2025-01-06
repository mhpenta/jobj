# jobj

`jobj` is a Go package that provides a flexible and type-safe way to create and manage JSON and XML schema definitions. It offers a fluent API for defining complex data structures with field validation, custom types, and schema generation capabilities.

## Features

- Fluent API for schema definition
- Support for both JSON and XML schema generation
- Custom date type handling with `JsonDate`
- Field validation and type checking
- Support for nested objects and arrays
- Enum-like field constraints using `AnyOf`
- Required/optional field specification
- Schema validation against Go structs

## Installation

```bash
go get github.com/mhpenta/jobj
```

## Usage

### Basic Example

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

### Working with JsonDate

The package includes a custom `JsonDate` type for handling dates in YYYY-MM-DD format:

```go
type Document struct {
    PublishDate JsonDate `json:"publish_date"`
}
```

### Schema Generation

Generate JSON schema:

```go
response := NewHeadlineResponse()
schemaJSON := response.GetSchemaString()
```

Generate XML schema:

```go
response := NewHeadlineResponse()
response.UseXML = true
schemaXML := response.GetSchemaString()
```

### Schema Validation

Validate that your schema matches a struct:

```go
type MyStruct struct {
    Headline    string  `json:"headline"`
    Confidence  float64 `json:"confidence"`
}

response := NewHeadlineResponse()
isValid := response.Validate(&MyStruct{})
```

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

## Note

This package is designed for creating schema definitions that can be used with JSON and XML data structures. It's particularly useful for defining API responses and data validation requirements.