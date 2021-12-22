package ui

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

//go:generate stringer -type=ElementType

type ElementType int

const (
	UIElement ElementType = iota
	PanelElement
	PaddingElement
	ColumnElement
	TextElement
	ButtonElement
	ImageElement
	IfElement
	RangeElement
)

type Element struct {
	Type       ElementType
	Attributes AttributeMap

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
	} else if getType(start) != UIElement {
		return nil, fmt.Errorf("first element not <UI />: %s", getType(start))
	}

	attrs := getAttributes(start)
	if err := validateAttributesForType(UIElement, attrs); err != nil {
		return nil, fmt.Errorf("validate attributes: %v", err)
	}

	kids, err := getContents(dec)
	if err != nil {
		return nil, fmt.Errorf("getContents: %v", err)
	}
	e := Element{
		Type:       getType(start),
		Attributes: attrs,
		Children:   kids,
	}

	// Consume any leftovers after the first top-level element.
	for {
		token, err := dec.Token()
		if token == nil && err == io.EOF {
			break
		}

		switch t := token.(type) {
		case xml.CharData:
			trailing := string(bytes.Trim(t, " \r\n\t"))
			if trailing != "" {
				return nil, fmt.Errorf("trailing characters after closing element: %s", trailing)
			}
		case xml.StartElement:
			return nil, fmt.Errorf("unexpected second top-level element: %T", t)
		}
	}

	return &e, nil
}

func getType(start xml.StartElement) ElementType {
	switch start.Name.Local {
	case "UI":
		return UIElement
	case "Panel":
		return PanelElement
	case "Padding":
		return PaddingElement
	case "Column":
		return ColumnElement
	case "Text":
		return TextElement
	case "Button":
		return ButtonElement
	case "Image":
		return ImageElement
	case "If":
		return IfElement
	case "Range":
		return RangeElement
	default:
		panic(fmt.Sprintf("unknown element %v", start.Name))
	}
}

func getAttributes(start xml.StartElement) map[string]string {
	result := map[string]string{}
	for _, attr := range start.Attr {
		result[attr.Name.Local] = attr.Value
	}
	return result
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

		case xml.CharData:
			raw := strings.TrimSpace(string(t))
			if raw == "" {
				break
			}

			e := Element{
				Type: TextElement,
				Attributes: map[string]string{
					"value": raw,
				},
				Children: []*Element{},
			}
			elements = append(elements, &e)
		}
	}
}
