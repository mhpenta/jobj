package jobj

type Field struct {
	ValueName            string
	ValueType            string
	ValueDescription     string
	Value                string
	ValueRequired        bool
	ValueAnyOf           []ConstDescription
	SubFields            []*Field
	AdditionalProperties bool // Default false for all, explicitly false for array
}

type ConstDescription struct {
	Const       string
	Description string
}

func Text(name string) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     "string",
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
		ValueType:     "boolean",
		ValueName:     name,
		ValueAnyOf:    nil,
	}
	return vb
}

func Float(name string) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     "number",
		ValueName:     name,
		ValueAnyOf:    nil,
	}
	return vb
}

func Date(name string) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     "string",
		//ValueFormat: "date"
		ValueName:  name,
		ValueAnyOf: nil,
	}
	return vb
}

func Array(name string, fields []*Field) *Field {
	vb := &Field{
		ValueRequired:        false,
		ValueType:            "array",
		ValueName:            name,
		ValueAnyOf:           nil,
		SubFields:            fields,
		AdditionalProperties: false,
	}
	return vb
}

func Int(name string) *Field {
	vb := &Field{
		ValueRequired: false,
		ValueType:     "integer",
		ValueName:     name,
		ValueAnyOf:    nil,
	}
	return vb
}

func Object(name string, fields []*Field) *Field {
	vb := &Field{
		ValueRequired:        false,
		ValueType:            "object",
		ValueName:            name,
		ValueAnyOf:           nil,
		SubFields:            fields,
		AdditionalProperties: false,
	}
	return vb
}

func (vb *Field) Type(valueType string) *Field {
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
