package jobj

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"reflect"
	"strings"
)

type Schema struct {
	Name        string
	Description string
	UseXML      bool
	Fields      []*Field
}

func (r *Schema) GetDescription() string {
	return r.Description
}

func (r *Schema) GetFields() []*Field {
	return r.Fields
}

func (r *Schema) GetSchemaString() string {
	if r.UseXML {
		return r.GetXMLSchemaString()
	}

	schema := struct {
		Schema      string                 `json:"$schema"`
		Definitions map[string]interface{} `json:"definitions"`
		Reference   string                 `json:"$ref"`
	}{
		Schema: "http://json-schema.org/draft-07/schema#",
		Definitions: map[string]interface{}{
			r.Name: map[string]interface{}{
				"properties":           r.FieldsJson(),
				"type":                 "object",
				"required":             r.RequiredFields(),
				"additionalProperties": false,
			},
		},
		Reference: "#/definitions/" + r.Name,
	}
	schemaJson, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Println("Error marshalling JSON schema:", err)
		return ""
	}
	return string(schemaJson)
}

func (r *Schema) FieldsJson() map[string]interface{} {
	properties := make(map[string]interface{})
	for _, field := range r.Fields {
		if field.ValueAnyOf != nil {
			anyOf := make([]map[string]interface{}, 0, len(field.ValueAnyOf))
			for _, enum := range field.ValueAnyOf {
				anyOf = append(anyOf, map[string]interface{}{
					"const":       enum.Const,
					"description": enum.Description,
				})
			}

			fieldProps := map[string]interface{}{
				"anyOf": anyOf,
			}
			if field.ValueDescription != "" {
				fieldProps["description"] = field.ValueDescription
			}

			properties[field.ValueName] = fieldProps
			continue
		}

		if field.ValueType == "array" && field.SubFields != nil {
			arrayFieldProperties := make(map[string]interface{})
			requiredFields := []string{}

			// Collect required fields once
			for _, subField := range field.SubFields {
				if subField.ValueRequired {
					requiredFields = append(requiredFields, subField.ValueName)
				}
			}

			// Process all fields
			for _, subField := range field.SubFields {
				if subField.ValueAnyOf != nil {
					anyOf := make([]map[string]interface{}, 0, len(subField.ValueAnyOf))
					for _, enum := range subField.ValueAnyOf {
						anyOf = append(anyOf, map[string]interface{}{
							"const":       enum.Const,
							"description": enum.Description,
						})
					}

					fieldProps := map[string]interface{}{
						"anyOf": anyOf,
					}
					if subField.ValueDescription != "" {
						fieldProps["description"] = subField.ValueDescription
					}

					arrayFieldProperties[subField.ValueName] = fieldProps
					continue
				}

				if subField.ValueType == "object" && subField.SubFields != nil {
					objectFieldProperties := make(map[string]interface{})
					for _, objField := range subField.SubFields {
						objectFieldProperties[objField.ValueName] = map[string]string{
							"type":        objField.ValueType,
							"description": objField.ValueDescription,
						}
					}
					arrayFieldProperties[subField.ValueName] = map[string]interface{}{
						"type":        subField.ValueType,
						"description": subField.ValueDescription,
						"properties":  objectFieldProperties,
						"required":    subField.getRequiredFields(),
					}
					continue
				}

				arrayFieldProperties[subField.ValueName] = map[string]string{
					"type":        subField.ValueType,
					"description": subField.ValueDescription,
				}
			}

			properties[field.ValueName] = map[string]interface{}{
				"type":                 field.ValueType,
				"description":          field.ValueDescription,
				"additionalProperties": field.AdditionalProperties,
				"items": map[string]interface{}{
					"type":       "object",
					"properties": arrayFieldProperties,
					"required":   requiredFields,
				},
			}
			continue
		}

		if field.ValueType == "object" {
			objectFieldProperties := make(map[string]interface{})
			for _, subField := range field.SubFields {
				objectFieldProperties[subField.ValueName] = map[string]string{
					"type":        subField.ValueType,
					"description": subField.ValueDescription,
				}
			}
			properties[field.ValueName] = map[string]interface{}{
				"type":        field.ValueType,
				"description": field.ValueDescription,
				"properties":  objectFieldProperties,
				"required":    field.getRequiredFields(),
			}
			continue
		}

		properties[field.ValueName] = map[string]string{
			"type":        field.ValueType,
			"description": field.ValueDescription,
		}
	}
	return properties
}

//func (r *Schema) FieldsJson() map[string]interface{} {
//	properties := make(map[string]interface{})
//	for _, field := range r.Fields {
//
//		if field.ValueAnyOf != nil {
//			anyOf := make([]map[string]interface{}, 0, len(field.ValueAnyOf))
//			for _, enum := range field.ValueAnyOf {
//				anyOf = append(anyOf, map[string]interface{}{
//					"const":       enum.Const,
//					"description": enum.Description,
//				})
//			}
//
//			fieldProps := map[string]interface{}{
//				"anyOf": anyOf,
//			}
//
//			if field.ValueDescription != "" {
//				fieldProps["description"] = field.ValueDescription
//			}
//
//			properties[field.ValueName] = fieldProps
//			continue
//		}
//
//		if field.ValueType == "array" {
//			if field.SubFields != nil {
//				arrayFieldProperties := make(map[string]interface{})
//				requiredFields := []string{}
//
//				// First collect all required fields, regardless of type
//				for _, subField := range field.SubFields {
//					if subField.ValueRequired {
//						requiredFields = append(requiredFields, subField.ValueName)
//					}
//				}
//
//				// Then process the fields
//				for _, subField := range field.SubFields {
//					if subField.ValueAnyOf != nil {
//						anyOf := make([]map[string]interface{}, 0, len(subField.ValueAnyOf))
//						for _, enum := range subField.ValueAnyOf {
//							anyOf = append(anyOf, map[string]interface{}{
//								"const":       enum.Const,
//								"description": enum.Description,
//							})
//						}
//
//						fieldProps := map[string]interface{}{
//							"anyOf": anyOf,
//						}
//						if subField.ValueDescription != "" {
//							fieldProps["description"] = subField.ValueDescription
//						}
//
//						arrayFieldProperties[subField.ValueName] = fieldProps
//						continue
//					}
//
//					// Handle nested objects within array items
//					if subField.ValueType == "object" && subField.SubFields != nil {
//						objectFieldProperties := make(map[string]interface{})
//						for _, objField := range subField.SubFields {
//							objectFieldProperties[objField.ValueName] = map[string]string{
//								"type":        objField.ValueType,
//								"description": objField.ValueDescription,
//							}
//						}
//						arrayFieldProperties[subField.ValueName] = map[string]interface{}{
//							"type":        subField.ValueType,
//							"description": subField.ValueDescription,
//							"properties":  objectFieldProperties,
//							"required":    subField.getRequiredFields(),
//						}
//						continue
//					}
//
//					arrayFieldProperties[subField.ValueName] = map[string]string{
//						"type":        subField.ValueType,
//						"description": subField.ValueDescription,
//					}
//				}
//
//				properties[field.ValueName] = map[string]interface{}{
//					"type":                 field.ValueType,
//					"description":          field.ValueDescription,
//					"additionalProperties": field.AdditionalProperties,
//					"items": map[string]interface{}{
//						"type":       "object",
//						"properties": arrayFieldProperties,
//						"required":   requiredFields,
//					},
//				}
//				continue
//			}
//		}
//
//		if field.ValueType == "object" {
//			objectFieldProperties := make(map[string]interface{})
//			for _, subField := range field.SubFields {
//				objectFieldProperties[subField.ValueName] = map[string]string{
//					"type":        subField.ValueType,
//					"description": subField.ValueDescription,
//				}
//			}
//			properties[field.ValueName] = map[string]interface{}{
//				"type":        field.ValueType,
//				"description": field.ValueDescription,
//				"properties":  objectFieldProperties,
//				"required":    field.getRequiredFields(),
//			}
//			continue
//		}
//
//		properties[field.ValueName] = map[string]string{
//			"type":        field.ValueType,
//			"description": field.ValueDescription,
//		}
//	}
//	return properties
//}

func (r *Schema) RequiredFields() []string {
	var required []string
	for _, field := range r.Fields {
		if field.ValueType == "array" && field.ValueRequired {
			required = append(required, field.ValueName)
			continue
		}

		if field.ValueRequired {
			required = append(required, field.ValueName)
		}
	}

	if len(required) == 0 {
		return nil
	}
	return required
}

func (r *Schema) GetXMLSchemaString() string {
	type xsDocumentation struct {
		XMLName xml.Name `xml:"xs:documentation"`
		Content string   `xml:",chardata"`
	}

	type xsAnnotation struct {
		XMLName       xml.Name        `xml:"xs:annotation"`
		Documentation xsDocumentation `xml:"xs:documentation"`
	}

	type xsElement struct {
		XMLName    xml.Name `xml:"xs:element"`
		Name       string   `xml:"name,attr"`
		Type       string   `xml:"type,attr,omitempty"`
		MinOccur   string   `xml:"minOccurs,attr,omitempty"`
		MaxOccur   string   `xml:"maxOccurs,attr,omitempty"`
		Annotation xsAnnotation
	}

	type xsComplexType struct {
		XMLName  xml.Name
		Name     string      `xml:"name,attr"`
		Sequence []xsElement `xml:"sequence>element"`
	}

	type xsSchema struct {
		XMLName     xml.Name      `xml:"xs:schema"`
		Xmlns       string        `xml:"xmlns:xs,attr"`
		Element     xsElement     `xml:"element"`
		ComplexType xsComplexType `xml:"complexType"`
	}

	elements := make([]xsElement, 0, len(r.Fields))
	for _, field := range r.Fields {
		xsType := mapGoTypeToXSType(field.ValueType)
		elem := xsElement{
			Name: field.ValueName,
			Type: xsType,
			Annotation: xsAnnotation{
				Documentation: xsDocumentation{
					Content: field.ValueDescription,
				},
			},
		}
		if field.ValueRequired {
			elem.MinOccur = "1"
			elem.MaxOccur = "1"
		} else {
			elem.MinOccur = "0"
			elem.MaxOccur = "1"
		}
		elements = append(elements, elem)
	}

	schema := xsSchema{
		Xmlns: "http://www.w3.org/2001/XMLSchema",
		Element: xsElement{
			Name: r.Name,
			Type: r.Name + "Type",
		},
		ComplexType: xsComplexType{
			Name:     r.Name + "Type",
			Sequence: elements,
		},
	}

	output, err := xml.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Println("Error marshalling XML Schema:", err)
		return ""
	}

	// Add XML header and correct namespaces
	header := xml.Header
	return header + string(output)
}

func (r *Schema) GetXMLSchemaString2() string {
	type xsDocumentation struct {
		XMLName xml.Name `xml:"xs:documentation"`
		Content string   `xml:",chardata"`
	}

	type xsAnnotation struct {
		XMLName       xml.Name        `xml:"xs:annotation"`
		Documentation xsDocumentation `xml:"xs:documentation"`
	}

	type xsElement struct {
		XMLName    xml.Name
		Name       string       `xml:"name,attr"`
		Type       string       `xml:"type,attr,omitempty"`
		MinOccur   string       `xml:"minOccurs,attr,omitempty"`
		MaxOccur   string       `xml:"maxOccurs,attr,omitempty"`
		Annotation xsAnnotation `xml:"annotation,omitempty"`
	}

	type xsComplexType struct {
		XMLName  xml.Name
		Name     string      `xml:"name,attr"`
		Sequence []xsElement `xml:"sequence>element"`
	}

	type xsSchema struct {
		XMLName     xml.Name      `xml:"xs:schema"`
		Xmlns       string        `xml:"xmlns:xs,attr"`
		Element     xsElement     `xml:"element"`
		ComplexType xsComplexType `xml:"complexType"`
	}

	elements := make([]xsElement, 0, len(r.Fields))
	for _, field := range r.Fields {
		xsType := mapGoTypeToXSType(field.ValueType)
		elem := xsElement{
			Name: field.ValueName,
			Type: xsType,
			Annotation: xsAnnotation{
				Documentation: xsDocumentation{
					Content: field.ValueDescription,
				},
			},
		}
		if field.ValueRequired {
			elem.MinOccur = "1"
			elem.MaxOccur = "1"
		} else {
			elem.MinOccur = "0"
			elem.MaxOccur = "1"
		}
		elements = append(elements, elem)
	}

	schema := xsSchema{
		Xmlns: "http://www.w3.org/2001/XMLSchema",
		Element: xsElement{
			Name: r.Name,
			Type: r.Name + "Type",
		},
		ComplexType: xsComplexType{
			Name:     r.Name + "Type",
			Sequence: elements,
		},
	}

	output, err := xml.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Println("Error marshalling XML Schema:", err)
		return ""
	}

	// Add XML header and correct namespaces
	header := xml.Header
	return header + string(output)
}

func mapGoTypeToXSType(goType string) string {
	switch goType {
	case "string":
		return "xs:string"
	case "boolean":
		return "xs:boolean"
	case "number":
		return "xs:float" // or xs:decimal, xs:int based on what is more appropriate
	default:
		return "xs:string" // default or unknown types
	}
}

// Validate checks if the response fields are present in the container struct, which must be passed into it since the
// embedded Schema struct does not have access to the container struct. Validation proves that the schema requested
// is able to be Marshalled by the container struct. Mostly used in testing. Returns true if all fields are present,
func (r *Schema) Validate(container interface{}) bool {
	containerValue := reflect.ValueOf(container)
	containerType := containerValue.Type()
	if containerType.Kind() != reflect.Struct && !(containerType.Kind() == reflect.Ptr && containerType.Elem().Kind() == reflect.Struct) {
		return false
	}
	if containerType.Kind() == reflect.Ptr {
		containerValue = containerValue.Elem()
		containerType = containerType.Elem()
	}
	fieldNames := make(map[string]struct{}, containerType.NumField())
	for i := 0; i < containerType.NumField(); i++ {
		field := containerType.Field(i)
		if field.Name != "Schema" {
			fieldNames[field.Tag.Get("json")] = struct{}{}
		}
	}
	for _, field := range r.Fields {
		fieldName := strings.ToLower(field.ValueName)
		if _, exists := fieldNames[fieldName]; !exists {
			return false
		}
	}
	return true
}

func (f *Field) getRequiredFields() []string {
	var required []string
	for _, subField := range f.SubFields {
		if subField.ValueRequired {
			required = append(required, subField.ValueName)
		}
	}
	return required
}
