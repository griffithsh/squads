package ui

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
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
	RangeElement
)

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

	case RangeElement:
		field := el.Attributes["over"]
		childDatas, err := dynamic.Ranger(field, data)
		switch {
		case err != nil:
			// It is not possible to range over this field.
			return 0, 0, err

		case len(el.Children) == 0:
			// There are no children, so the range takes no space.
			return 0, 0, nil

		case el.Children[0].Type == ColumnElement:
			// The children of the Range are Columns.
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
