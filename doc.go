// Package jobj implements a subset of JSON Schema Draft-07 for use with large language models.
//
// Example:
//
//	type Response struct {
//	    jobj.Schema
//	}
//
//	func NewResponse() *Response {
//	    r := &Response{}
//	    r.Name = "Response"
//	    r.Fields = []*jobj.Field{
//	        jobj.Text("message").Required(),
//	        jobj.Int("count"),
//	    }
//	    return r
//	}
//
// The package supports basic types, arrays, objects and enums. External schema
// references are not supported. XML support is limited to basic type mapping.
package jobj
