package ui

import (
	"fmt"
	"image"
	"io"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
)

// UI is a Component that represents a UI. The stutter is unfortunate ...
type UI struct {
	Doc *Element

	Data interface{}
}

// Type of this Component.
func (*UI) Type() string {
	return "UI"
}

// NewUI construct a new UI Component from a declared XML template. You're
// responsible for assigning Data to the Component before rendering it.
func NewUI(r io.Reader) *UI {
	el, err := parse(r)
	if err != nil {
		panic(fmt.Sprintf("parse UI template: %v", err))
	}
	return &UI{
		Doc:  el,
		Data: map[string]func(){},
	}
}
