package ui

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/griffithsh/squads/ui/dynamic"
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
	ForElement
	RangeElement
)

func EvaluateIfExpression(expr string, data interface{}) bool {
	buf := bytes.NewBuffer([]byte{})
	if err := template.Must(template.New("text").Parse(fmt.Sprintf("{{ if %s }}true{{ end}}", expr))).Execute(buf, data); err != nil {
		panic(fmt.Errorf("execute: %v", err))
	}
	return buf.String() == "true"
}

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
	case IfElement:
		if !EvaluateIfExpression(el.Attributes["expr"], data) {
			return 0, 0, nil
		}
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

	case ForElement:
		// TODO:

	case RangeElement:
		field := el.Attributes["over"]
		childDatas, err := dynamic.Ranger(field, data)
		switch {
		case err != nil:
			return 0, 0, err
		case len(el.Children) == 0:
			return 0, 0, nil

		case el.Children[0].Type == ColumnElement:
			maxHeight := 0
			for i, column := range el.Children {
				maxWidth := mult(maxWidth, 1.0/12) * column.Attributes.Twelfths()
				_, h, err := column.DimensionsWith(childDatas[i], maxWidth)
				if err != nil {
					return 0, 0, err
				}
				if h > maxHeight {
					maxHeight = h
				}
			}
			return maxWidth, h, nil

		default:
			width, height := 0, 0
			for _, child := range el.Children {
				w, h, err := child.DimensionsWith(data, maxWidth)
				if err != nil {
					return 0, 0, err
				}
				if w > width {
					width = w
				}
				if h > height {
					height = h
				}
			}
			return width, height, nil
		}
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
	case "For":
		return ForElement
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

// permittedAttributes defines which attributes are allowed on a given element
// type. What's an XML schema anyway?
// First Element must be a "UI" element.
// A Panel has a width and a height attribute which must change based on scale
// the panel's L-T-R-B coordinates are determined by the alignment and padding of its parent,
// the scale

var permittedAttributes = map[ElementType][]string{
	// UIElement may not have width or height; it always takes up 100% of the screen.
	UIElement:      {"align", "valign"},
	PanelElement:   {"width", "height", "valign", "align", "outline"},
	PaddingElement: {"all", "vertical", "horizontal", "top", "bottom", "left", "right", "valign", "align"},
	ColumnElement:  {"twelfths", "align"},
	TextElement:    {"value", "size", "layout", "color", "width"},
	ButtonElement:  {"onclick", "label", "width"},
	ImageElement:   {"texture", "width", "height", "x", "y", "intangible"},
	IfElement:      {"expr"},
	ForElement:     {"index", "length"},
	RangeElement:   {"over"},
}

func validateAttributesForType(t ElementType, attrs map[string]string) error {
	permitted := permittedAttributes[t]
	for attr := range attrs {
		for _, yes := range permitted {
			if attr == yes {
				goto ok
			}
		}
		return fmt.Errorf("attribute %q is illegal on %s", attr, t.String())
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
