package ui

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/griffithsh/squads/ui/dynamic"
)

// heightOfDescendants of an Element. Passing a slice of datas means that this
// function is appropriate to call on elements with column and non-column
// children.
func heightOfDescendants(element *Element, datas []interface{}, maxWidth int) int {
	if len(element.Children) == 0 {
		return 0
	}

	tallest := 0
	sumHeights := 0

	// When we pass mismatched numbers of children versus datas it might mean
	// that we've inside a RangeElement, and the RangeElement has multiple
	// chldren. The factor here is how far through the datas we should be as a
	// ratio.
	factor := len(element.Children) / len(datas)
	for i, child := range element.Children {
		data := datas[i/factor]

		height := heightOf(child, data, maxWidth)
		if height > tallest {
			tallest = height
		}
		sumHeights += height
	}

	if element.Children[0].Type == ColumnElement {
		return tallest
	}
	return sumHeights
}

// heightOf any element. This is how much the element claims from the available
// layout, and not the visible height. An intangible <Image /> therefore has
// zero height.
func heightOf(element *Element, data interface{}, maxWidth int) int {
	switch element.Type {
	case ButtonElement:
		return ButtonHeight

	case ImageElement:
		if element.Attributes.Intangible() {
			return 0
		}
		height, err := ResolveInt(element.Attributes["height"], data)
		if err != nil {
			return 0
		}
		return height

	case TextElement:
		buf := bytes.NewBuffer([]byte{})
		if err := template.Must(template.New("text").Parse(element.Attributes["value"])).Execute(buf, data); err != nil {
			return 0
		}
		return heightOfText(buf.String(), element.Attributes.FontSize(), maxWidth, element.Attributes.FontLayout())

	case UIElement:
		return heightOfDescendants(element, []interface{}{data}, maxWidth)

	case PanelElement:
		if _, ok := element.Attributes["height"]; ok {
			return element.Attributes.Height()
		}
		return heightOfDescendants(element, []interface{}{data}, maxWidth)

	case PaddingElement:
		height := heightOfDescendants(element, []interface{}{data}, maxWidth)
		return height + element.Attributes.TopPadding() + element.Attributes.BottomPadding()

	case IfElement:
		if !EvaluateIfExpression(element.Attributes["expr"], data) {
			return 0
		}
		return heightOfDescendants(element, []interface{}{data}, maxWidth)

	case ColumnElement:
		return heightOfDescendants(element, []interface{}{data}, maxWidth)

	case RangeElement:
		field := element.Attributes["over"]
		childDatas, err := dynamic.Ranger(field, data)
		if err != nil {
			panic(fmt.Sprintf("could not range over %q of %v: %v", field, data, err))
		}
		pseudo := *element
		pseudo.Children = make([]*Element, 0, len(childDatas)*len(element.Children))
		for i := 0; i < len(childDatas); i++ {
			pseudo.Children = append(pseudo.Children, element.Children...)
		}
		return heightOfDescendants(&pseudo, childDatas, maxWidth)

	default:
		panic(fmt.Sprintf("heightOf does not handle %q", element.Type))
	}
}

// heightOfText calculates how many pixels high a given text should be.
func heightOfText(value string, size TextSize, maxWidth int, align TextLayout) (height int) {
	text := NewText(value, size)

	// Spacer around each text instance.
	spacer := TextPadding

	// We know our max width, so we can split long lines.
	width := maxWidth
	splitLines := SplitLines(text.Lines, width)

	y := spacer
	for i, line := range splitLines {
		if i != 0 {
			// If not the first line, add a line spacer.
			y += LineSpacing(size)
		}

		tallest := 0
		for _, word := range line {
			for _, char := range word.Characters {
				if char.Height > tallest {
					tallest = char.Height
				}
			}
		}

		y += tallest
	}

	return spacer + y
}
