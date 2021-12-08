package ui

import (
	"fmt"
	"strconv"
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
	RangeElement:   {"over"},
}

func validateAttributesForType(t ElementType, attrs AttributeMap) error {
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
