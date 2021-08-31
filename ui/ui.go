package ui

import (
	"fmt"
	"image"
	"io"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui/dynamic"
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
	for _, e := range uis.mgr.Get([]string{"UI"}) {
		uic := uis.mgr.Component(e, "UI").(*UI)
		// Figure out if the interaction we've got is positioned over a button
		// element, and if so, trigger its handler. The first button we find is
		// fine.

		/* We have access to the entire UI component, so we can figure out
		 * anything we need. The problem becomes how to share code that figures
		 * out dimensions and positioning between this and the ui visualiser.
		 * Can we start by just duplicating to help identify the common parts?
		 * */

		// uic.Doc.Type == ui.UIElement - first element must be UI!
		// uic.Doc.Attributes == [] - none are allowed!
		var f func(children []*Element, data interface{}, bounds image.Rectangle, align, valign string, scale float64) (image.Rectangle, error)
		f = func(children []*Element, data interface{}, bounds image.Rectangle, align, valign string, scale float64) (image.Rectangle, error) {
			maxColHeight := 0
			for _, child := range children {
				var err error
				switch child.Type {
				case PanelElement:
					w, h, err := child.DimensionsWith(data, bounds.Dx(), scale)
					if err != nil {
						return bounds, err
					}
					panelBounds := image.Rectangle{}
					switch align {
					default:
						fallthrough
					case "left":
						panelBounds.Min.X = bounds.Min.X
					case "right":
						panelBounds.Min.X = bounds.Max.X - w
					case "center":
						panelBounds.Min.X = bounds.Min.X + (bounds.Max.X-bounds.Min.X)/2 - w/2
					}
					switch valign {
					default:
						fallthrough
					case "top":
						panelBounds.Min.Y = bounds.Min.Y
					case "bottom":
						panelBounds.Min.Y = bounds.Max.Y - h
					case "middle":
						panelBounds.Min.Y = bounds.Min.Y + (bounds.Max.Y-bounds.Min.Y)/2 - h/2
					}

					panelBounds.Max = image.Point{
						X: panelBounds.Min.X + w,
						Y: panelBounds.Min.Y + h,
					}

					if bounds, err = f(child.Children, data, panelBounds, child.Attributes.Align(), child.Attributes.Valign(), scale); err != nil {
						return bounds, err
					}

				case PaddingElement:
					paddedBounds := bounds
					paddedBounds.Min.X += int(float64(child.Attributes.LeftPadding()) * scale)
					paddedBounds.Max.X -= int(float64(child.Attributes.RightPadding()) * scale)
					paddedBounds.Min.Y += int(float64(child.Attributes.TopPadding()) * scale)
					paddedBounds.Max.Y -= int(float64(child.Attributes.BottomPadding()) * scale)
					if bounds, err = f(child.Children, data, paddedBounds, child.Attributes.Align(), child.Attributes.Valign(), scale); err != nil {
						return bounds, err
					}

				case ColumnElement:
					// I think we need to know about siblings to do this correctly?
					// I don't think we can stomp bounds here?  Only the last Column of
					// adjacent siblings is block level.
					colBounds := bounds
					colBounds.Min.X += bounds.Dx() * child.Attributes.TwelfthsOffset() / 12
					w := bounds.Dx() * child.Attributes.Twelfths() / 12
					colBounds.Max.X = colBounds.Min.X + w
					takenBounds, err := f(child.Children, data, colBounds, child.Attributes.Align(), child.Attributes.Valign(), scale)
					if err != nil {
						return bounds, err
					}
					colHeight := takenBounds.Min.Y - bounds.Min.Y
					if colHeight > maxColHeight {
						maxColHeight = colHeight
					}

					// If the twelfths and the twelfths-offset total the full width of a
					// set of columns, then we know that this is the final column of a
					// group.
					if child.Attributes.Twelfths()+child.Attributes.TwelfthsOffset() == 12 {
						bounds.Min.Y += maxColHeight
						maxColHeight = 0
					}

				case TextElement:
					_, h, err := child.DimensionsWith(data, bounds.Dx(), scale)
					if err != nil {
						return bounds, fmt.Errorf("DimensionsWith: %v", err)
					}
					bounds.Min.Y += h

				case ButtonElement:
					// buttonHeight := int(ButtonHeight * scale)
					w, h, err := child.DimensionsWith(data, bounds.Dx(), scale)
					if err != nil {
						return bounds, err
					}
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

					// Is this interaction within this buttonDimensions?
					p := image.Point{int(ev.AbsoluteX), int(ev.AbsoluteY)}
					if p.In(buttonDimensions) {
						if err := dynamic.Call(child.Attributes["onclick"], data); err != nil {
							panic(fmt.Sprintf("dynamic call: %v", err))
						}
						// FIXME: somehow escape this recursion, signalling that
						// the click has been consumed.
					}

					bounds.Min.Y += h

				case ImageElement:
					if !child.Attributes.Intangible() {
						_, h, err := child.DimensionsWith(data, bounds.Dx(), scale)
						if err != nil {
							return bounds, fmt.Errorf("DimensionsWith: %v", err)
						}
						// NB this is descaled. why is this descaled?
						bounds.Min.Y += h //int(float64(h) / scale)
					}
				}
			}
			return bounds, nil
		}

		bounds := image.Rect(0, 0, uis.screenW, uis.screenH)
		scale := 2.0 // FIXME: this needs to be shared with the uiVisualiser!

		f(uic.Doc.Children, uic.Data, bounds, "center", "middle", scale)
	}
}

// heightOfText FIXME: should not accept scale, should accept a maxWidth int instead of a bounds Rectangle
func heightOfText(value string, size TextSize, bounds image.Rectangle, align TextLayout, scale float64) (height int) {
	text := NewText(value, size)

	// Spacer around each text instance.
	spacer := int(TextPadding * scale)

	// We know our bounds now, so we can split long lines.
	width := int(float64(bounds.Dx()) / scale)
	splitLines := SplitLines(text.Lines, width)

	y := float64(bounds.Min.Y + spacer)
	for i, line := range splitLines {
		x := float64(bounds.Min.X)
		if i != 0 {
			// If not the first line, add a line spacer.
			y += float64(LineSpacing(size)) * scale
		}

		// Different strategies based on width and word breaks...
		switch align {
		case TextLayoutRight:
			x = float64(bounds.Max.X) - float64(line.Width())*scale
		case TextLayoutCenter:
			x += float64(bounds.Dx()/2) - float64(line.Width()/2)*scale
		}

		tallest := 0
		wordSpace := SpaceWidth * scale
		if align == TextLayoutJustify && len(line) > 1 {
			extra := float64((float64(bounds.Dx()) - float64(line.Width())*scale) / float64(len(line)-1))
			wordSpace += extra
		}
		for _, word := range line {
			for i, char := range word.Characters {
				if char.Height > tallest {
					tallest = char.Height
				}
				x += float64(char.Width) * scale

				// Add spacing between letters for every letter except the last one.
				if i != len(word.Characters)-1 {
					x += float64(LetterSpacing) * scale
				}
			}
			x += wordSpace
		}

		y += float64(tallest) * scale
	}

	return spacer + int(y) - bounds.Min.Y
}

// heightOfText2 calculates how many unscaled pixels high a given text should
// be. It replaces (the deprecated) heightOfText.
func heightOfText2(value string, size TextSize, maxWidth int, align TextLayout) (height int) {
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
