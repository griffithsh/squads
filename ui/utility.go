package ui

import (
	"bytes"
	"fmt"
	"image"
	"strconv"
	"text/template"
)

// Resolve the value of text by executing it as a template with the provided
// data source.
// Values like "45" are returned directly, while value like {{ .Age }} are
// resolved by pulling the field (if it exists) with the name "Age" from data.
func Resolve(text string, data interface{}) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := template.Must(template.New("").Parse(text)).Execute(buf, data); err != nil {
		return "", fmt.Errorf("execute: %v, template: %q", err, text)
	}
	return buf.String(), nil
}

// ResolveInt functions like Resolve, but returns an integral value.
func ResolveInt(text string, data interface{}) (int, error) {
	str, err := Resolve(text, data)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

func mult(num int, scale float64) int {
	return int(float64(num) * scale)
}

// AlignedXY calculates the top-left x,y coordinates for an element with the
// provided dimensions and alignments within a parent bounds.
func AlignedXY(w, h int, parent image.Rectangle, align, valign string) (x, y int) {
	switch align {
	default:
		fallthrough
	case "left":
		x = parent.Min.X
	case "right":
		x = parent.Max.X - w
	case "center":
		x = parent.Min.X + (parent.Max.X-parent.Min.X)/2 - w/2
	}
	switch valign {
	default:
		fallthrough
	case "top":
		y = parent.Min.Y
	case "bottom":
		y = parent.Max.Y - h
	case "middle":
		y = parent.Min.Y + (parent.Max.Y-parent.Min.Y)/2 - h/2
	}

	return x, y
}

func EvaluateIfExpression(expr string, data interface{}) bool {
	buf := bytes.NewBuffer([]byte{})
	if err := template.Must(template.New("text").Parse(fmt.Sprintf("{{ if %s }}true{{ end}}", expr))).Execute(buf, data); err != nil {
		panic(fmt.Errorf("execute: %v", err))
	}
	return buf.String() == "true"
}
