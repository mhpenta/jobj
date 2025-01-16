// Package funcschema generates JSON schemas from function signatures for use with LLM tools.
//
// Example:
//
// The package expects functions with signature func(context.Context, any) (string, error)
// and generates schemas based on the second parameter's struct fields.
package funcschema
