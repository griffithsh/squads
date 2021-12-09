package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"text/template"

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
	screenW, screenH int
}

func NewUISystem(mgr *ecs.World, bus *event.Bus) *UISystem {
	uis := UISystem{
		mgr:     mgr,
		screenW: 0, screenH: 0,
	}

	bus.Subscribe(Interact{}.Type(), func(t event.Typer) {
		ev := t.(*Interact)
		uis.Handle(ev)
	})
	bus.Subscribe(game.WindowSizeChanged{}.Type(), func(e event.Typer) {
		wsc := e.(*game.WindowSizeChanged)
		uis.screenW, uis.screenH = wsc.NewW, wsc.NewH

	})
	return &uis
}

func (uis *UISystem) Handle(ev *Interact) {
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
}

// calculateChildren dispatches to either calculateColumnChildren or
// calculateNonColumnChildren depending on if the first child is a ColumnElement
// or not.
func (sys *UISystem) calculateChildren(root *UI, children []*Element, data interface{}, bounds image.Rectangle, align, valign string, depth int) (image.Rectangle, error) {
	// If the first child is a column, make the big assumption that all children are columns.
	if len(children) > 0 && children[0].Type == ColumnElement {
		// Make a copy of data for every child.
		datas := make([]interface{}, len(children))
		for i := range datas {
			datas[i] = data
		}

		return sys.calculateColumnChildren(root, children, datas, bounds, depth+1)
	}
	return sys.calculateNonColumnChildren(root, children, data, bounds, align, valign, depth+1)
}

func (sys *UISystem) calculateColumnChildren(root *UI, columns []*Element, datas []interface{}, bounds image.Rectangle, depth int) (image.Rectangle, error) {
	if len(columns) != len(datas) {
		return bounds, fmt.Errorf("mismatch between columns(%d) and datas(%d)", len(columns), len(datas))
	}
	maxColHeight := 0
	twelfthOffset := 0
	for i, column := range columns {
		if column.Type != ColumnElement {
			return bounds, fmt.Errorf("non-column type %T", column.Type)
		}
		colBounds := bounds
		colBounds.Min.X += bounds.Dx() * twelfthOffset / 12
		w := bounds.Dx() * column.Attributes.Twelfths() / 12
		colBounds.Max.X = colBounds.Min.X + w
		takenBounds, err := sys.calculateChildren(root, column.Children, datas[i], colBounds, column.Attributes.Align(), column.Attributes.Valign(), depth)
		if err != nil {
			return bounds, err
		}
		colHeight := takenBounds.Min.Y - bounds.Min.Y
		if colHeight > maxColHeight {
			maxColHeight = colHeight
		}
		twelfthOffset += column.Attributes.Twelfths()
	}

	// max column height
	bounds.Min.Y += maxColHeight
	return bounds, nil
}

func (sys *UISystem) calculateNonColumnChildren(root *UI, children []*Element, data interface{}, bounds image.Rectangle, align, valign string, depth int) (image.Rectangle, error) {
	maxWidth := bounds.Dx()
	widestChild := 0

	for _, child := range children {
		switch child.Type {
		case PanelElement:
			w, h, err := child.DimensionsWith(data, maxWidth)
			if err != nil {
				return bounds, err
			}
			x, y := AlignedXY(w, h, bounds, align, valign)

			panelBounds := image.Rect(x, y, x+w, y+h)
			if invis := child.Attributes["outline"]; invis != "false" {
				root.renderinstructions = append(root.renderinstructions, PanelRenderInstruction{
					Bounds: panelBounds,
				})
			}

			if bounds, err = sys.calculateChildren(root, child.Children, data, panelBounds, child.Attributes.Align(), child.Attributes.Valign(), depth); err != nil {
				return bounds, err
			}
			if widestChild < bounds.Dx() {
				widestChild = bounds.Dx()
			}

		case PaddingElement:
			paddedBounds := bounds
			paddedBounds.Min.X += child.Attributes.LeftPadding()
			paddedBounds.Min.Y += child.Attributes.TopPadding()
			paddedBounds.Max.X -= child.Attributes.RightPadding()
			paddedBounds.Max.Y -= child.Attributes.BottomPadding()

			w, h, err := child.DimensionsWith(data, paddedBounds.Dx())
			if err != nil {
				return bounds, err
			}
			x, y := AlignedXY(w, h, paddedBounds, align, valign)

			childrenBounds := image.Rect(x, y, x+w, y+h)
			if _, err = sys.calculateChildren(root, child.Children, data, childrenBounds, child.Attributes.Align(), child.Attributes.Valign(), depth); err != nil {
				return bounds, err
			}
			if widestChild < bounds.Dx() {
				widestChild = bounds.Dx()
			}
			bounds.Min.Y += h

		case TextElement:
			label := child.Attributes["value"]
			sz := child.Attributes.FontSize()
			layout := child.Attributes.FontLayout()
			buf := bytes.NewBuffer([]byte{})
			if err := template.Must(template.New("text").Parse(label)).Execute(buf, data); err != nil {
				return bounds, fmt.Errorf("execute: %v, template: %q", err, label)
			}

			txtBounds := bounds
			maxTextWidth := maxWidth
			if child.Attributes["width"] != "" {
				maxTextWidth = child.Attributes.Width()
			}
			w, h, err := child.DimensionsWith(data, maxTextWidth)
			if err != nil {
				return bounds, err
			}
			txtBounds.Max.X = txtBounds.Min.X + w
			root.renderinstructions = append(root.renderinstructions, TextRenderInstruction{
				Text:   buf.String(),
				Size:   sz,
				Bounds: txtBounds,
				Layout: layout,
			})
			bounds.Min.Y += h
			if widestChild < txtBounds.Dx() {
				widestChild = txtBounds.Dx()
			}

		case ButtonElement:
			label, err := Resolve(child.Attributes["label"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve button label: %v", err)
			}
			w, h, err := child.DimensionsWith(data, maxWidth)
			if err != nil {
				return bounds, err
			}
			// Does the parent align left, right, or centre? Are we valigning
			// it? Calculate buttonDimensions from that.
			l := bounds.Min.X
			switch align {
			case "right":
				l = bounds.Max.X - w
			case "center":
				l = bounds.Min.X + (bounds.Max.X-bounds.Min.X)/2 - w/2
			default: // left
			}
			t := bounds.Min.Y
			switch valign {
			case "bottom":
				t = bounds.Max.Y - h
			case "middle":
				t = bounds.Min.Y + (bounds.Max.Y-bounds.Min.Y)/2 - h/2
			default: // top
			}
			buttonDimensions := image.Rect(l, t, l+w, t+h)
			root.interactives = append(root.interactives, InteractiveRegion{
				Bounds: buttonDimensions,
				Handler: func() {
					if err := dynamic.Call(child.Attributes["onclick"], data); err != nil {
						panic(fmt.Sprintf("dynamic call: %v", err))
					}
				},
			})
			root.renderinstructions = append(root.renderinstructions, ButtonRenderInstruction{
				Active: false,
				Bounds: buttonDimensions,
				Label:  label,
			})
			bounds.Min.Y += h
			if widestChild < buttonDimensions.Dx() {
				widestChild = buttonDimensions.Dx()
			}

		case ImageElement:
			texture, err := Resolve(child.Attributes["texture"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve texture: %v", err)
			}
			width, err := ResolveInt(child.Attributes["width"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve width: %v", err)
			}
			height, err := ResolveInt(child.Attributes["height"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve height: %v", err)
			}
			x, err := ResolveInt(child.Attributes["x"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve x: %v", err)
			}
			y, err := ResolveInt(child.Attributes["y"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve y: %v", err)
			}
			root.renderinstructions = append(root.renderinstructions, ImageRenderInstruction{
				Texture: texture,
				From:    image.Rect(x, y, x+width, y+height),
				AtX:     float64(bounds.Min.X),
				AtY:     float64(bounds.Min.Y),
			})

			if !child.Attributes.Intangible() {
				bounds.Min.Y += height
				if widestChild < width {
					widestChild = width
				}
			}

		case IfElement:
			expr := child.Attributes["expr"]
			if EvaluateIfExpression(expr, data) {
				w, h, err := child.DimensionsWith(data, maxWidth)
				if err != nil {
					return bounds, err
				}
				x, y := AlignedXY(w, h, bounds, align, valign)

				childrenBounds := image.Rect(x, y, x+w, y+h)
				if bounds, err = sys.calculateChildren(root, child.Children, data, childrenBounds, child.Attributes.Align(), child.Attributes.Valign(), depth); err != nil {
					return bounds, err
				}
				if widestChild < bounds.Dx() {
					widestChild = bounds.Dx()
				}
			}

		case RangeElement:
			field := child.Attributes["over"]
			childDatas, err := dynamic.Ranger(field, data)
			if err != nil {
				str, _ := json.Marshal(data)
				return bounds, fmt.Errorf("range over %q in %v: %v", field, string(str), err)
			}
			switch {
			case len(child.Children) == 0:
				// It's a no-op if there are no children!

			case child.Children[0].Type == ColumnElement:
				vchildren := make([]*Element, len(childDatas))
				for i := 0; i < len(childDatas); i++ {
					vchildren[i] = child.Children[0]
				}
				bounds, err := sys.calculateColumnChildren(root, vchildren, childDatas, bounds, depth)
				if err != nil {
					return bounds, err
				}

			default: // Sibling context is irrelevant when not dealing with columns.
				for _, item := range childDatas {
					w, h, err := child.DimensionsWith(data, maxWidth)
					if err != nil {
						return bounds, err
					}
					x, y := AlignedXY(w, h, bounds, align, valign)

					childrenBounds := image.Rect(x, y, x+w, y+h)
					if bounds, err = sys.calculateChildren(root, child.Children, item, childrenBounds, child.Attributes.Align(), child.Attributes.Valign(), depth); err != nil {
						return bounds, err
					}
					if widestChild < bounds.Dx() {
						widestChild = bounds.Dx()
					}
				}

			}

		}
	}
	bounds.Max.X = bounds.Min.X + widestChild
	return bounds, nil
}

func (sys *UISystem) Update() error {
	uiScale := 2.0 // FIXME!
	screen := image.Rect(0, 0, int(float64(sys.screenW)/uiScale), int(float64(sys.screenH)/uiScale))
	for _, e := range sys.mgr.Get([]string{"UI"}) {
		uic := sys.mgr.Component(e, "UI").(*UI)

		uic.interactives = uic.interactives[:0]
		uic.renderinstructions = uic.renderinstructions[:0]

		_, err := sys.calculateChildren(uic, uic.Doc.Children, uic.Data, screen, uic.Doc.Attributes.Align(), uic.Doc.Attributes.Valign(), 0)
		if err != nil {
			return err
		}
	}
	return nil
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
