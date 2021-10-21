package ui

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"
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

func (m AttributeMap) LeftPadding() int {
	if left, ok := m["left"]; ok {
		if padding, err := strconv.Atoi(left); err == nil {
			return padding
		}
	}
	if horz, ok := m["horizontal"]; ok {
		if padding, err := strconv.Atoi(horz); err == nil {
			return padding
		}
	}
	if all, ok := m["all"]; ok {
		if padding, err := strconv.Atoi(all); err == nil {
			return padding
		}
	}

	return 0
}
func (m AttributeMap) RightPadding() int {
	if right, ok := m["right"]; ok {
		if padding, err := strconv.Atoi(right); err == nil {
			return padding
		}
	}
	if horz, ok := m["horizontal"]; ok {
		if padding, err := strconv.Atoi(horz); err == nil {
			return padding
		}
	}
	if all, ok := m["all"]; ok {
		if padding, err := strconv.Atoi(all); err == nil {
			return padding
		}
	}

	return 0
}
func (m AttributeMap) TopPadding() int {
	if top, ok := m["top"]; ok {
		if padding, err := strconv.Atoi(top); err == nil {
			return padding
		}
	}
	if vert, ok := m["vertical"]; ok {
		if padding, err := strconv.Atoi(vert); err == nil {
			return padding
		}
	}
	if all, ok := m["all"]; ok {
		if padding, err := strconv.Atoi(all); err == nil {
			return padding
		}
	}

	return 0
}
func (m AttributeMap) BottomPadding() int {
	if bottom, ok := m["bottom"]; ok {
		if padding, err := strconv.Atoi(bottom); err == nil {
			return padding
		}
	}
	if vert, ok := m["vertical"]; ok {
		if padding, err := strconv.Atoi(vert); err == nil {
			return padding
		}
	}
	if all, ok := m["all"]; ok {
		if padding, err := strconv.Atoi(all); err == nil {
			return padding
		}
	}

	return 0
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

func (m AttributeMap) Intangible() bool {
	str, ok := m["intangible"]
	if !ok {
		return false
	}
	return str == "true"
}

type Element struct {
	Type       ElementType
	Attributes AttributeMap

	Children []*Element
}

// widestChild of an element, taking into account that the widestChild might be
// the maxWidth if there are any Columns.
func (el *Element) widestChild(data interface{}, maxWidth int) (int, error) {
	maxChild := 0
	for _, child := range el.Children {
		if child.Type == ColumnElement {
			// We can riskily shortcut here so long as we can assume that a set
			// of columns always take up 100% of the width.
			return maxWidth, nil
		}
		w, _, err := child.DimensionsWith(data, maxWidth)
		if err != nil {
			return 0, fmt.Errorf("dimensions of %s: %v", child.Type, err)
		}

		if w > maxChild {
			maxChild = w
		}
	}
	return maxChild, nil
}

// DimensionsWith calculates the dimensions of an Element with the given data
// and available dimensions. Text elements vary in size due to the data that
// could be inserted into them.
func (el *Element) DimensionsWith(data interface{}, maxWidth int) (w, h int, err error) {
	switch el.Type {
	case PanelElement:
		width, err := ResolveInt(el.Attributes["width"], data)
		if err != nil {
			// Since there's no configured width, use the width of the widest child.
			width = maxWidth
			maxChild, err := el.widestChild(data, maxWidth)
			if err != nil {
				return 0, 0, fmt.Errorf("widestChild of %s: %v", el.Type, err)
			}
			if maxChild < width {
				width = maxChild
			}
		}
		height, err := ResolveInt(el.Attributes["height"], data)
		if err != nil {
			// Since there was no height configured for this, let's sum the
			// children's heights.
			height = 0
			for _, child := range el.Children {
				_, h, err := child.DimensionsWith(data, maxWidth)
				if err != nil {
					return 0, 0, fmt.Errorf("dimensions of %s: %v", el.Type, err)
				}
				height += h
			}
		}
		return width, height, nil

	case PaddingElement:
		width, err := ResolveInt(el.Attributes["width"], data)
		if err != nil {
			// Since there's no configured width, use the width of the widest child.
			width = maxWidth
			maxChild, err := el.widestChild(data, maxWidth)
			if err != nil {
				return 0, 0, fmt.Errorf("widestChild of %s: %v", el.Type, err)
			}
			if maxChild < width {
				width = maxChild
			}
		}
		height, err := ResolveInt(el.Attributes["height"], data)
		if err != nil {
			// Since there was no height configured for this, let's sum the
			// children's heights with the vertical padding values.
			height = el.Attributes.TopPadding() + el.Attributes.BottomPadding()
			for _, child := range el.Children {
				_, h, err := child.DimensionsWith(data, maxWidth)
				if err != nil {
					return 0, 0, fmt.Errorf("dimensions of %s: %v", child.Type, err)
				}
				height += h
			}
		}
		return width, height, nil
	case ColumnElement:
		width := mult(maxWidth, 1.0/12) * el.Attributes.Twelfths()
		height := 0
		for _, child := range el.Children {
			_, h, err := child.DimensionsWith(data, maxWidth)
			if err != nil {
				return 0, 0, fmt.Errorf("dimensions of %s: %v", el.Type, err)
			}
			height += h
		}
		return width, height, nil

	case TextElement:
		label := el.Attributes["value"]
		sz := el.Attributes.FontSize()
		layout := el.Attributes.FontLayout()
		buf := bytes.NewBuffer([]byte{})
		if err := template.Must(template.New("text").Parse(label)).Execute(buf, data); err != nil {
			return 0, 0, fmt.Errorf("execute: %v, template: %q", err, label)
		}
		width := maxWidth
		height := 0
		if el.Attributes["width"] != "" {
			width, _ = ResolveInt(el.Attributes["width"], data)
		}
		height = heightOfText(buf.String(), sz, width, layout)
		return width, height, nil

	case ButtonElement:
		width, err := ResolveInt(el.Attributes["width"], data)
		if err != nil {
			// FIXME: It might be worth calculating an appropriate width for the
			// button given its text label instead of relying on a default.
			width = DefaultButtonWidth
		}
		return width, ButtonHeight, nil

	case ImageElement:
		width, err := ResolveInt(el.Attributes["width"], data)
		if err != nil {
			return 0, 0, fmt.Errorf("ResolveInt width: %v", err)
		}
		height, err := ResolveInt(el.Attributes["height"], data)
		if err != nil {
			return 0, 0, fmt.Errorf("ResolveInt height: %v", err)
		}
		return width, height, nil
	}
	panic(fmt.Sprintf("unhandled element type: %q", el.Type))
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
	PaddingElement: {"all", "vertical", "horizontal", "top", "bottom", "left", "right"},
	ColumnElement:  {"twelfths", "align"},
	TextElement:    {"value", "size", "layout", "color", "width"},
	ButtonElement:  {"onclick", "label", "width"},
	ImageElement:   {"texture", "width", "height", "x", "y", "intangible"},
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
