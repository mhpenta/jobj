package jobj

type DataType string

const (
	TypeObject  DataType = "object"
	TypeString  DataType = "string"
	TypeNumber  DataType = "number"
	TypeInteger DataType = "integer"
	TypeBoolean DataType = "boolean"
	TypeArray   DataType = "array"
)

type Field struct {
	ValueName            string
	ValueType            DataType
	ValueDescription     string
	Value                string
	ValueRequired        bool
	ValueAnyOf           []ConstDescription
	SubFields            []*Field
	AdditionalProperties bool // Default false for all, explicitly false for array
	ArrayItemType        DataType // For arrays of primitives (when SubFields is nil/empty)
}

type ConstDescription struct {
	Const       string
	Description string
}

func Text(name string) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     TypeString,
		ValueName:     name,
		ValueAnyOf:    nil,
	}
	return vb
}

func AnyOf(name string, enums []ConstDescription) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     "anyOf",
		ValueName:     name,
		ValueAnyOf:    enums,
	}
	return vb
}

func Bool(name string) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     TypeBoolean,
		ValueName:     name,
		ValueAnyOf:    nil,
	}
	return vb
}

func Float(name string) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     TypeNumber,
		ValueName:     name,
		ValueAnyOf:    nil,
	}
	return vb
}

func Date(name string) *Field {
	/*
	   JSON Schema spec represents dates as strings with a "format" attribute.
	   The proper schema would be: {"type": "string", "format": "date"}

	   We don't implement format here for simplicity, keeping the Field struct
	   minimal while covering the most common use cases for JSON Schema generation.

	   Future enhancement: Add a ValueFormat field to the Field struct and update
	   the schema generation to include format when specified:

	   type Field struct {
	       // existing fields...
	       ValueFormat string // For date, date-time, email, etc.
	   }

	   Then schema generation would include:
	   if field.ValueFormat != "" {
	       fieldProps["format"] = field.ValueFormat
	   }
	*/
	vb := &Field{
		ValueRequired: false,
		ValueType:     TypeString,
		ValueName:     name,
		ValueAnyOf:    nil,
	}
	return vb
}

func Array(name string, fields []*Field) *Field {
	vb := &Field{
		ValueRequired:        false,
		ValueType:            TypeArray,
		ValueName:            name,
		ValueAnyOf:           nil,
		SubFields:            fields,
		AdditionalProperties: false,
	}
	return vb
}

// ArrayOf creates an array field with primitive item types (e.g., []string, []int)
func ArrayOf(name string, itemType DataType) *Field {
	vb := &Field{
		ValueRequired:        false,
		ValueType:            TypeArray,
		ValueName:            name,
		ValueAnyOf:           nil,
		SubFields:            nil,
		ArrayItemType:        itemType,
		AdditionalProperties: false,
	}
	return vb
}

func Int(name string) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     TypeInteger,
		ValueName:     name,
		ValueAnyOf:    nil,
	}
	return vb
}

func Object(name string, fields []*Field) *Field {
	vb := &Field{
		ValueRequired:        false,
		ValueType:            TypeObject,
		ValueName:            name,
		ValueAnyOf:           nil,
		SubFields:            fields,
		AdditionalProperties: false,
	}
	return vb
}

func (vb *Field) Type(valueType DataType) *Field {
	vb.ValueType = valueType
	return vb
}

func (vb *Field) SetValue(value string) *Field {
	vb.Value = value
	return vb
}

func (vb *Field) Desc(valueDescription string) *Field {
	vb.ValueDescription = valueDescription
	return vb
}

func (vb *Field) Required() *Field {
	vb.ValueRequired = true
	return vb
}

func (vb *Field) Optional() *Field {
	vb.ValueRequired = false
	return vb
}
