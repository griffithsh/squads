package ui

import (
	"encoding/xml"
	"fmt"
	"io"
)

//go:generate stringer -type=ElementType
type ElementType int

const (
	PanelElement ElementType = iota
	RowElement
	ColumnElement
	TextElement
	ButtonElement
	ImageElement
)

type Element struct {
	Type       ElementType
	Attributes map[string]interface{}

	Children []*Element
}

func parse(r io.Reader) (*Element, error) {
	dec := xml.NewDecoder(r)
	token, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("get first token: %v", err)
	}
	start, ok := token.(xml.StartElement)
	if !ok {
		return nil, fmt.Errorf("first token not xml.StartElement: %T", token)
	}
	kids, err := getContents(dec)
	if err != nil {
		return nil, fmt.Errorf("getContents: %v", err)
	}
	e := Element{
		Type:       getType(start),
		Attributes: getAttributes(start),
		Children:   kids,
	}
	if token, err = dec.Token(); token != nil || err != io.EOF {
		return nil, fmt.Errorf("unexpected second top-level element: %v, should be io.EOF, got %v", token, err)
	}
	return &e, nil
}

func getType(start xml.StartElement) ElementType {
	switch start.Name.Local {
	case "Panel":
		return PanelElement
	case "Row":
		return RowElement
	case "Column":
		return ColumnElement
	case "Text":
		return TextElement
	case "Button":
		return ButtonElement
	case "Image":
		return ImageElement
	default:
		panic(fmt.Sprintf("unknown element %v", start.Name))
	}
}

func getAttributes(start xml.StartElement) map[string]interface{} {
	result := map[string]interface{}{}
	for _, attr := range start.Attr {
		result[attr.Name.Local] = attr.Value
	}
	return result
}

var permittedAttributes = map[ElementType][]string{
	PanelElement:  {"width"},
	RowElement:    {},
	ColumnElement: {"twelths"},
	TextElement:   {"value"},
	ButtonElement: {"onclick", "label"},
	ImageElement:  {"texture", "width", "height", "x", "y"},
}

func validateAttributesForType(t ElementType, attrs map[string]interface{}) error {
	permitted := permittedAttributes[t]
	for attr := range attrs {
		for _, yes := range permitted {
			if attr == yes {
				goto ok
			}
		}
		return fmt.Errorf("invalid attribute %s for type %s", attr, t.String())
	ok:
	}
	return nil
}

// getContents recursively retrieves the contents of an element. The passed
// Decoder is assumed to have just called Token, and received an
// xml.StartElement.
func getContents(dec *xml.Decoder) ([]*Element, error) {
	var elements []*Element

	for {
		token, err := dec.Token()
		if err != nil {
			// this is invalid, because if we have just opened an element, then
			// we need to get at least an end element.
			return nil, fmt.Errorf("error while element unclosed: %v", err)
		}
		switch t := token.(type) {
		case xml.EndElement:
			// ok! we can return here, because there's nothing more to do for this element
			return elements, nil
		case xml.StartElement:
			// assign attrs: t.Attr
			// recurse to extract children
			kids, err := getContents(dec)
			if err != nil {
				return nil, err
			}
			elementType := getType(t)
			attrs := getAttributes(t)
			if err := validateAttributesForType(elementType, attrs); err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			e := Element{
				Type:       elementType,
				Attributes: attrs,
				Children:   kids,
			}
			elements = append(elements, &e)
		}
	}
}
