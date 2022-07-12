package ui

import (
	"fmt"
	"image"
	"io"
	"math"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui/dynamic"
)

type InteractiveRegion struct {
	Bounds  image.Rectangle
	Handler func()
}

type RenderInstruction interface {
}

type TextRenderInstruction struct {
	Text   string
	Size   TextSize
	Bounds image.Rectangle
	Layout TextLayout
}
type ImageRenderInstruction struct {
	Texture  string
	From     image.Rectangle
	AtX, AtY float64
}
type PanelRenderInstruction struct {
	Bounds image.Rectangle
}
type ButtonRenderInstruction struct {
	Active bool
	Bounds image.Rectangle
	Label  string
}

// UI is a Component that represents a UI. The stutter is unfortunate ...
type UI struct {
	Doc *Element

	Data interface{}

	interactives       []InteractiveRegion
	renderinstructions []RenderInstruction
}

func (c *UI) RenderInstructions() []RenderInstruction {
	return c.renderinstructions
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

	// Do we need to know the dimensions of all components at this point? Where
	// are the buttons? How can we bottom-align things if we don't know how tall
	// any elements are?
	// The alternative is to determine the positions and dimensions of all
	// elements at the time of Interaction.
	/*
		To know how big an element is, you'd need to know about data (for text
		elements), and sibling elements (for columns?), and the dimensions of child
		elements (for anything with children).  Should do that padding element type
		as well first. Maybe that would be the needed precedent for a column wrapper
		type element.
	*/
	return &UI{
		Doc:  el,
		Data: map[string]func(){},
	}
}

type UISystem struct {
	mgr              *ecs.World
	bus              *event.Bus
	screenW, screenH int
}

func NewUISystem(mgr *ecs.World, bus *event.Bus) *UISystem {
	uis := UISystem{
		mgr:     mgr,
		bus:     bus,
		screenW: 0, screenH: 0,
	}

	bus.Subscribe(UIInteract{}.Type(), func(t event.Typer) {
		ev := t.(*UIInteract)
		uis.Handle(ev)
	})
	bus.Subscribe(game.WindowSizeChanged{}.Type(), func(e event.Typer) {
		wsc := e.(*game.WindowSizeChanged)
		uis.screenW, uis.screenH = wsc.NewW, wsc.NewH

	})
	return &uis
}

func (sys *UISystem) Update() error {
	uiScale := 2.0 // FIXME!
	screen := image.Rect(0, 0, int(float64(sys.screenW)/uiScale), int(float64(sys.screenH)/uiScale))
	for _, e := range sys.mgr.Get([]string{"UI"}) {
		uic := sys.mgr.Component(e, "UI").(*UI)

		uic.interactives = uic.interactives[:0]
		uic.renderinstructions = uic.renderinstructions[:0]

		// NB, uic.Doc.Type == UIElement

		err := calculateChildren(uic, uic.Doc.Children, uic.Data, screen, uic.Doc.Attributes.Align(), uic.Doc.Attributes.Valign(), 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uis *UISystem) Handle(ev *UIInteract) {
	uiScale := 2.0 // FIXME this needs to come from somewhere ...

	interactPoint := image.Point{int(ev.AbsoluteX / uiScale), int(ev.AbsoluteY / uiScale)}
	for _, e := range uis.mgr.Get([]string{"UI"}) {
		uic := uis.mgr.Component(e, "UI").(*UI)

		for _, interactive := range uic.interactives {
			if interactPoint.In(interactive.Bounds) {
				interactive.Handler()
				return
			}
		}
	}

	// Unhandled by any UI interactive region, so pass the event on.
	uis.bus.Publish(&Interact{X: ev.X, Y: ev.Y, AbsoluteX: ev.AbsoluteX, AbsoluteY: ev.AbsoluteY})
}

// realiseChildren replaces IfElements and RangeElements with their children as
// the data dictates.  Also returns the datas that should be applied per child -
// remember RangeElements are bound to alternative data sources.
func realiseChildren(children []*Element, data interface{}) ([]*Element, []interface{}) {
	result := make([]*Element, 0, len(children))
	datas := make([]interface{}, 0, len(children))
	for _, child := range children {
		switch child.Type {
		case IfElement:
			if EvaluateIfExpression(child.Attributes["expr"], data) {
				c, d := realiseChildren(child.Children, data)
				result = append(result, c...)
				datas = append(datas, d...)
			}
		case RangeElement:
			field := child.Attributes["over"]
			childDatas, err := dynamic.Ranger(field, data)
			if err != nil {
				panic(fmt.Sprintf("dynamic.Ranger: %v", err))
			}
			for _, childData := range childDatas {
				c, d := realiseChildren(child.Children, childData)
				result = append(result, c...)
				datas = append(datas, d...)
			}

		default:
			result = append(result, child)
			datas = append(datas, data)
		}
	}
	return result, datas
}

// widthOf returns the configured width of an element. Does not recurse.
// Generally infers width from the width attribute.
func widthOf(child *Element, data interface{}, available image.Rectangle) int {

	availableWidth := available.Dx()

	attrWidth, errWidthAttr := ResolveInt(child.Attributes["width"], data)

	w := availableWidth
	switch child.Type {
	case TextElement:
		// FIXME: if a text is short enough to not take up the available width,
		// then it should be align-able. Need a widthOfText function though ...
		if errWidthAttr == nil {
			w = attrWidth
		}

	case ButtonElement:
		if errWidthAttr == nil {
			w = attrWidth
		}

	default:
		if errWidthAttr == nil {
			w = attrWidth
		}
	}
	return w
}

func alignmentOf(child *Element, data interface{}, available image.Rectangle, align, valign string) (x, y int) {
	// A shortcut appears...
	if align == "left" && valign == "top" {
		return available.Min.X, available.Min.Y
	}

	availableWidth := available.Dx()
	availableHeight := available.Dy()

	attrWidth, errWidthAttr := ResolveInt(child.Attributes["width"], data)
	attrHeight, errHeightAttr := ResolveInt(child.Attributes["height"], data)

	w, h := availableWidth, availableHeight
	switch child.Type {
	case TextElement:
		// FIXME: if a text is short enough to not take up the available width,
		// then it should be align-able. Need a widthOfText function though ...
		if errWidthAttr == nil {
			w = attrWidth
		}
		txt, err := Resolve(child.Attributes["value"], data)
		if err != nil {
			panic(fmt.Sprintf("resolve TextElement value: %v", err))
		}
		h = heightOfText(txt, child.Attributes.FontSize(), w, child.Attributes.FontLayout())

	case ButtonElement:
		if errWidthAttr == nil {
			w = attrWidth
		}
		h = ButtonHeight

	default:
		if errWidthAttr == nil {
			w = attrWidth
		}
		if errHeightAttr == nil {
			h = attrHeight
		}
	}

	switch align {
	default:
		fallthrough
	case "left":
		x = available.Min.X
	case "right":
		x = available.Max.X - w
	case "center":
		x = available.Min.X + (available.Max.X-available.Min.X)/2 - w/2
	}
	switch valign {
	default:
		fallthrough
	case "top":
		y = available.Min.Y
	case "bottom":
		y = available.Max.Y - h
	case "middle":
		y = available.Min.Y + (available.Max.Y-available.Min.Y)/2 - h/2
	}

	return x, y
}

// calculateChildren dispatches to either calculateColumnChildren or
// calculateNonColumnChildren depending on if the first child is a ColumnElement
// or not.
func calculateChildren(root *UI, children []*Element, data interface{}, available image.Rectangle, align, valign string, depth int) error {
	realChildren, datas := realiseChildren(children, data)

	maxHeight := 0
	sumHeights := 0

	twelfthOffset := 0

	for i, child := range realChildren {
		data := datas[i]
		height := heightOf(child, data, available.Dx())
		width := widthOf(child, data, available)
		x, y := alignmentOf(child, data, available, align, valign)
		bounds := image.Rect(x, y, x+width, y+height)
		switch child.Type {

		// If and Range should be absorbed by realiseChildren().
		default:
			fallthrough
		case IfElement:
			fallthrough
		case RangeElement:
			panic(fmt.Sprintf("element type %q is unacceptable in this context", child.Type))

		case PanelElement:
			if invis := child.Attributes["outline"]; invis != "false" {
				root.renderinstructions = append(root.renderinstructions, PanelRenderInstruction{
					Bounds: bounds,
				})
			}

			err := calculateChildren(root, child.Children, data, bounds, child.Attributes.Align(), child.Attributes.Valign(), depth+1)
			if err != nil {
				return fmt.Errorf("<%s>: %v", child.Type, err)
			}

		case PaddingElement:
			paddedBounds := bounds
			paddedBounds.Min.X += child.Attributes.LeftPadding()
			paddedBounds.Min.Y += child.Attributes.TopPadding()
			paddedBounds.Max.X -= child.Attributes.RightPadding()
			paddedBounds.Max.Y -= child.Attributes.BottomPadding()

			err := calculateChildren(root, child.Children, data, paddedBounds, child.Attributes.Align(), child.Attributes.Valign(), depth+1)
			if err != nil {
				return fmt.Errorf("<%s>: %v", child.Type, err)
			}

		case ColumnElement:
			twelfths := child.Attributes.Twelfths()
			columnBounds := bounds
			columnBounds.Min.X += int(math.Round(float64(width) * float64(twelfthOffset) / 12))
			w := int(math.Round(float64(width) * float64(twelfths) / 12))
			columnBounds.Max.X = columnBounds.Min.X + w

			err := calculateChildren(root, child.Children, data, columnBounds, child.Attributes.Align(), child.Attributes.Valign(), depth+1)
			if err != nil {
				return fmt.Errorf("<%s>: %v", child.Type, err)
			}

			twelfthOffset += twelfths

		case ImageElement:
			texture, err := Resolve(child.Attributes["texture"], data)
			if err != nil {
				return fmt.Errorf("resolve texture: %v", err)
			} else if texture == "" {
				return fmt.Errorf("resolve texture from %q output empty string", child.Attributes["texture"])
			}
			fromX, err := ResolveInt(child.Attributes["x"], data)
			if err != nil {
				return fmt.Errorf("resolve x: %v", err)
			}
			fromY, err := ResolveInt(child.Attributes["y"], data)
			if err != nil {
				return fmt.Errorf("resolve y: %v", err)
			}
			// NB, we don't need to grab the width attribute again. Intangible only nixes height.
			fromHeight, err := ResolveInt(child.Attributes["height"], data)
			if err != nil {
				return fmt.Errorf("resolve height: %v", err)
			}

			root.renderinstructions = append(root.renderinstructions, ImageRenderInstruction{
				Texture: texture,
				From:    image.Rect(fromX, fromY, fromX+width, fromY+fromHeight),
				AtX:     float64(x),
				AtY:     float64(y),
			})

			if onclick := child.Attributes["onclick"]; onclick != "" {
				id, err := Resolve(child.Attributes["id"], data)
				if err != nil {
					return fmt.Errorf("Resolve %s: %v", child.Attributes["id"], err)
				}
				root.interactives = append(root.interactives, InteractiveRegion{
					Bounds: image.Rect(x, y, x+width, y+height),
					Handler: func() {
						if err := dynamic.Call(onclick, data, id); err != nil {
							panic(fmt.Sprintf("dynamic call: %v", err))
						}
					},
				})
			}

		case TextElement:
			label := child.Attributes["value"]
			sz := child.Attributes.FontSize()
			layout := child.Attributes.FontLayout()
			label, err := Resolve(label, data)
			if err != nil {
				return fmt.Errorf("Resolve %s: %v", child.Attributes["value"], err)
			}

			txtBounds := available
			txtBounds.Max.X = txtBounds.Min.X + width

			root.renderinstructions = append(root.renderinstructions, TextRenderInstruction{
				Text:   label,
				Size:   sz,
				Bounds: txtBounds,
				Layout: layout,
			})

		case ButtonElement:
			label, err := Resolve(child.Attributes["label"], data)
			if err != nil {
				return fmt.Errorf("resolve button label: %v", err)
			}

			buttonDimensions := image.Rect(x, y, x+width, y+height)

			id, err := Resolve(child.Attributes["id"], data)
			if err != nil {
				return fmt.Errorf("Resolve %s: %v", child.Attributes["id"], err)
			}
			root.interactives = append(root.interactives, InteractiveRegion{
				Bounds: buttonDimensions,
				Handler: func() {
					if err := dynamic.Call(child.Attributes["onclick"], data, id); err != nil {
						panic(fmt.Sprintf("dynamic call: %v", err))
					}
				},
			})
			root.renderinstructions = append(root.renderinstructions, ButtonRenderInstruction{
				Active: false,
				Bounds: buttonDimensions,
				Label:  label,
			})
		}

		if child.Type == ColumnElement {
			if height > maxHeight {
				maxHeight = height
			}
		} else {
			available.Min.Y += height
			sumHeights += height
		}
	}
	// if children were columns, return max height
	// else return sum of heights
	available.Min.X += maxHeight + sumHeights

	return nil
}
