package jobj

import (
	"encoding/xml"
	"log"
)

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
		return "xs:float"
	default:
		return "xs:string"
	}
}
