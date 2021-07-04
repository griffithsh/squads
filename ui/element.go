package ui

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
)

//go:generate stringer -type=ElementType
type ElementType int

const (
	UIElement ElementType = iota
	PanelElement
	PaddingElement
	RowElement
	ColumnElement
	TextElement
	ButtonElement
	ImageElement
)

type AttributeMap map[string]string

func (m AttributeMap) Width() int {
	str, ok := m["width"]
	if !ok {
		// then auto?
		return 0
	}

	w, err := strconv.Atoi(str)
	if err != nil {
		// invalid value?
		return 0
	}
	return w
}
func (m AttributeMap) Height() int {
	str, ok := m["height"]
	if !ok {
		// then auto?
		return 0
	}

	h, err := strconv.Atoi(str)
	if err != nil {
		// invalid value?
		return 0
	}
	return h
}

func (m AttributeMap) X() int {
	str, ok := m["x"]
	if !ok {
		// then auto?
		return 0
	}

	x, err := strconv.Atoi(str)
	if err != nil {
		// invalid value?
		return 0
	}
	return x
}

func (m AttributeMap) Y() int {
	str, ok := m["y"]
	if !ok {
		// then auto?
		return 0
	}

	y, err := strconv.Atoi(str)
	if err != nil {
		// invalid value?
		return 0
	}
	return y
}

func (m AttributeMap) Padding() int {
	str, ok := m["all"]
	if !ok {
		return 0
	}

	padding, err := strconv.Atoi(str)
	if err != nil {
		// invalid value?
		return 0
	}
	return padding
}

func (m AttributeMap) Twelfths() int {
	str, ok := m["twelfths"]
	if !ok {
		return 0
	}

	twelfths, err := strconv.Atoi(str)
	if err != nil {
		// invalid value?
		return 0
	}
	return twelfths
}

func (m AttributeMap) TwelfthsOffset() int {
	str, ok := m["twelfths-offset"]
	if !ok {
		return 0
	}

	off, err := strconv.Atoi(str)
	if err != nil {
		// invalid value?
		return 0
	}
	return off
}

func (m AttributeMap) Align() string {
	v := m["align"]
	switch v {
	case "left":
		fallthrough
	case "right":
		fallthrough
	case "center":
		return v

	default:
		return "left"
	}
}

func (m AttributeMap) Valign() string {
	v := m["valign"]
	switch v {
	case "top":
		fallthrough
	case "bottom":
		fallthrough
	case "middle":
		return v

	default:
		return "top"
	}
}

func (m AttributeMap) FontSize() TextSize {
	switch m["size"] {
	default:
		fallthrough
	case "normal":
		return TextSizeNormal
	case "small":
		return TextSizeSmall
	}
}

func (m AttributeMap) FontLayout() TextLayout {
	switch m["layout"] {
	default:
		fallthrough
	case "left":
		return TextLayoutLeft
	case "right":
		return TextLayoutRight
	case "justify":
		return TextLayoutJustify
	case "center":
		return TextLayoutCenter
	}
}

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

	kids, err := getContents(dec)
	if err != nil {
		return nil, fmt.Errorf("getContents: %v", err)
	}
	e := Element{
		Type:       getType(start),
		Attributes: getAttributes(start),
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

	// Normalise Column twelfths and inject twelfth-offsets ...
	var recurse func([]*Element)
	recurse = func(elements []*Element) {
		collected := []*Element{}

		for _, el := range elements {
			if el.Type == ColumnElement {
				collected = append(collected, el)
			} else if len(collected) > 0 {
				normaliseTwelfths(collected)
				collected = collected[:0]
			}
			recurse(el.Children)
		}
		if len(collected) > 0 {
			normaliseTwelfths(collected)
		}
	}
	recurse(e.Children)

	return &e, nil
}

func normaliseTwelfths(cols []*Element) {
	// first figure out the exact twelfth values for each element
	unspecified := []*Element{}
	unclaimed := 12
	for _, el := range cols {
		twelfths := el.Attributes.Twelfths()
		if twelfths == 0 {
			unspecified = append(unspecified, el)
		}
		unclaimed -= twelfths
	}

	origUnclaimed := unclaimed
	for i, el := range unspecified {
		twelfths := 0
		if i == 0 {
			// Add modulo to the first column - it's as good as anywhere.
			twelfths += origUnclaimed % len(unspecified)
		}

		// now assign an equal division of the remainder
		twelfths += origUnclaimed / len(unspecified)

		unclaimed -= twelfths
		el.Attributes["twelfths"] = strconv.Itoa(twelfths)
	}

	if unclaimed != 0 {
		// panic!
		panic(fmt.Sprintf("unclaimed nonzero: %d", unclaimed))
	}

	// then figure out the start-offsets
	for _, el := range cols {
		twelfth := el.Attributes.Twelfths()
		el.Attributes["twelfths-offset"] = strconv.Itoa(unclaimed)
		unclaimed += twelfth
	}
}

func getType(start xml.StartElement) ElementType {
	switch start.Name.Local {
	case "UI":
		return UIElement
	case "Panel":
		return PanelElement
	case "Padding":
		return PaddingElement
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

func getAttributes(start xml.StartElement) map[string]string {
	result := map[string]string{}
	for _, attr := range start.Attr {
		result[attr.Name.Local] = attr.Value
	}
	return result
}

// permittedAttributes defines which attributes are allowed on a given element
// type. What's an XML schema anyway?
// First Element must be a "UI" element.
// A Panel has a width and a height attribute which must change based on scale
// the panel's L-T-R-B coordinates are determined by the alignment and padding of its parent,
// the scale

var permittedAttributes = map[ElementType][]string{
	UIElement:      {},
	PanelElement:   {"width", "height"},
	PaddingElement: {"all"},
	RowElement:     {"align"},
	ColumnElement:  {"twelfths", "align"},
	TextElement:    {"value", "size", "layout", "color", "width"},
	ButtonElement:  {"onclick", "label", "width"},
	ImageElement:   {"texture", "width", "height", "x", "y"},
}

func validateAttributesForType(t ElementType, attrs map[string]string) error {
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
