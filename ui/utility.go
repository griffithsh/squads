package ui

import (
	"bytes"
	"fmt"
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

func EvaluateIfExpression(expr string, data interface{}) bool {
	buf := bytes.NewBuffer([]byte{})
	if err := template.Must(template.New("text").Parse(fmt.Sprintf("{{ if %s }}true{{ end}}", expr))).Execute(buf, data); err != nil {
		panic(fmt.Errorf("execute: %v", err))
	}
	return buf.String() == "true"
}
