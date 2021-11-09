package output

import (
	"bytes"
	"fmt"
	"image"
	"text/template"

	"github.com/griffithsh/squads/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type uiVisualizer struct {
	picForTexture func(filename string) (*ebiten.Image, error)
}

func newUIVisualizer(picForTexture func(filename string) (*ebiten.Image, error)) *uiVisualizer {
	uv := uiVisualizer{
		picForTexture: picForTexture,
	}
	return &uv
}

func (uv *uiVisualizer) Render(screen *ebiten.Image, uic *ui.UI) error {
	_, err := uv.drawChildren(screen, uic.Doc.Children, uic.Data, screen.Bounds(), uic.Doc.Attributes.Align(), uic.Doc.Attributes.Valign())

	return err
}

// drawChildren returns the remainder of the bounds, unused by the drawn children.
func (uv *uiVisualizer) drawChildren(screen *ebiten.Image, children []*ui.Element, data interface{}, bounds image.Rectangle, align, valign string) (image.Rectangle, error) {
	maxColHeight := 0
	maxWidth := bounds.Dx()
	widestChild := 0
	for _, child := range children {
		switch child.Type {
		case ui.PanelElement:
			w, h, err := child.DimensionsWith(data, maxWidth)
			if err != nil {
				return bounds, err
			}
			x, y := ui.AlignedXY(w, h, bounds, align, valign)

			panelBounds := image.Rect(x, y, x+w, y+h)
			if invis := child.Attributes["outline"]; invis != "false" {
				if err = uv.drawPanel(screen, panelBounds); err != nil {
					return bounds, err
				}
			}

			if bounds, err = uv.drawChildren(screen, child.Children, data, panelBounds, child.Attributes.Align(), child.Attributes.Valign()); err != nil {
				return bounds, err
			}
			if widestChild < bounds.Dx() {
				widestChild = bounds.Dx()
			}

		case ui.PaddingElement:
			paddedBounds := bounds
			paddedBounds.Min.X += child.Attributes.LeftPadding()
			paddedBounds.Min.Y += child.Attributes.TopPadding()
			paddedBounds.Max.X -= child.Attributes.RightPadding()
			paddedBounds.Max.Y -= child.Attributes.BottomPadding()

			w, h, err := child.DimensionsWith(data, paddedBounds.Dx())
			if err != nil {
				return bounds, err
			}
			x, y := ui.AlignedXY(w, h, paddedBounds, align, valign)

			childrenBounds := image.Rect(x, y, x+w, y+h)
			if bounds, err = uv.drawChildren(screen, child.Children, data, childrenBounds, child.Attributes.Align(), child.Attributes.Valign()); err != nil {
				return bounds, err
			}
			if widestChild < bounds.Dx() {
				widestChild = bounds.Dx()
			}

		case ui.ColumnElement:
			// I think we need to know about siblings to do this correctly?
			// I don't think we can stomp bounds here?  Only the last Column of
			// adjacent siblings is block level.
			colBounds := bounds
			colBounds.Min.X += bounds.Dx() * child.Attributes.TwelfthsOffset() / 12
			w := bounds.Dx() * child.Attributes.Twelfths() / 12
			colBounds.Max.X = colBounds.Min.X + w
			takenBounds, err := uv.drawChildren(screen, child.Children, data, colBounds, child.Attributes.Align(), child.Attributes.Valign())
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
				widestChild = bounds.Dx()
			}

		case ui.TextElement:
			label := child.Attributes["value"]
			sz := child.Attributes.FontSize()
			layout := child.Attributes.FontLayout()
			buf := bytes.NewBuffer([]byte{})
			if err := template.Must(template.New("text").Parse(label)).Execute(buf, data); err != nil {
				return bounds, fmt.Errorf("execute: %v, template: %q", err, label)
			}

			txtBounds := bounds
			if child.Attributes["width"] != "" {
				w, _, err := child.DimensionsWith(data, maxWidth)
				if err != nil {
					return bounds, err
				}
				txtBounds.Max.X = txtBounds.Min.X + w
			}
			h, err := uv.drawText(screen, buf.String(), sz, txtBounds, layout)
			if err != nil {
				return bounds, err
			}
			bounds.Min.Y += h
			if widestChild < txtBounds.Dx() {
				widestChild = txtBounds.Dx()
			}

		case ui.ButtonElement:
			label, err := ui.Resolve(child.Attributes["label"], data)
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
			if err = uv.drawButton(screen, false, buttonDimensions); err != nil {
				return bounds, err
			}

			// Push the text down a bit so it's vertically centered. 2 seems good?
			buttonDimensions.Min.Y += 2

			uv.drawText(screen, label, ui.TextSizeNormal, buttonDimensions, ui.TextLayoutCenter)
			bounds.Min.Y += h
			if widestChild < buttonDimensions.Dx() {
				widestChild = buttonDimensions.Dx()
			}

		case ui.ImageElement:
			texture, err := ui.Resolve(child.Attributes["texture"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve texture: %v", err)
			}
			img, err := uv.picForTexture(texture)
			if err != nil {
				return bounds, fmt.Errorf("picForTexture: %v", err)
			}
			width, err := ui.ResolveInt(child.Attributes["width"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve width: %v", err)
			}
			height, err := ui.ResolveInt(child.Attributes["height"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve height: %v", err)
			}
			x, err := ui.ResolveInt(child.Attributes["x"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve x: %v", err)
			}
			y, err := ui.ResolveInt(child.Attributes["y"], data)
			if err != nil {
				return bounds, fmt.Errorf("resolve y: %v", err)
			}
			img = img.SubImage(image.Rect(x, y, x+width, y+height)).(*ebiten.Image)
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(bounds.Min.X), float64(bounds.Min.Y))
			screen.DrawImage(img, &op)

			if !child.Attributes.Intangible() {
				bounds.Min.Y += height
				if widestChild < width {
					widestChild = width
				}
			}
		case ui.IfElement:
			expr := child.Attributes["expr"]
			if ui.EvaluateIfExpression(expr, data) {
				w, h, err := child.DimensionsWith(data, maxWidth)
				if err != nil {
					return bounds, err
				}
				x, y := ui.AlignedXY(w, h, bounds, align, valign)

				childrenBounds := image.Rect(x, y, x+w, y+h)
				if bounds, err = uv.drawChildren(screen, child.Children, data, childrenBounds, child.Attributes.Align(), child.Attributes.Valign()); err != nil {
					return bounds, err
				}
				if widestChild < bounds.Dx() {
					widestChild = bounds.Dx()
				}
			}
		}
	}
	bounds.Max.X = bounds.Min.X + widestChild
	return bounds, nil
}

// drawPanel renders a nine-slice "panel" that covers the provided Rectangle,
// and with a pixel granularity controlled by scale. The provided scale does not
// affect the portion of the screen covered by the panel.
func (uv *uiVisualizer) drawPanel(screen *ebiten.Image, r image.Rectangle) error {
	img, err := uv.picForTexture("ui.png")
	if err != nil {
		return fmt.Errorf("picForTexture: %v", err)
	}
	imgs := [9]*ebiten.Image{
		img.SubImage(image.Rect(0, 0, 4, 4)).(*ebiten.Image),
		img.SubImage(image.Rect(4, 0, 8, 4)).(*ebiten.Image),
		img.SubImage(image.Rect(8, 0, 12, 4)).(*ebiten.Image),
		img.SubImage(image.Rect(0, 4, 4, 8)).(*ebiten.Image),
		img.SubImage(image.Rect(4, 4, 8, 8)).(*ebiten.Image),
		img.SubImage(image.Rect(8, 4, 12, 8)).(*ebiten.Image),
		img.SubImage(image.Rect(0, 8, 4, 12)).(*ebiten.Image),
		img.SubImage(image.Rect(4, 8, 8, 12)).(*ebiten.Image),
		img.SubImage(image.Rect(8, 8, 12, 12)).(*ebiten.Image),
	}

	return nineSlice(screen, r, imgs, 4, 4)
}

func nineSlice(screen *ebiten.Image, r image.Rectangle, img [9]*ebiten.Image, pad, stride int) error {

	// centersWide, middlesTall is how many copies of the center and middle
	// pieces of the 9-slice are required given this current scale.
	centersWide := int(float64(r.Dx()-(pad+pad)) / (float64(stride)))
	middlesTall := int(float64(r.Dy()-(pad+pad)) / (float64(stride)))

	// top row
	tl := img[0]
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	screen.DrawImage(tl, &op)

	tc := img[1]
	for i := 0; i < centersWide; i++ {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(r.Min.X+pad+i*stride), float64(r.Min.Y))
		screen.DrawImage(tc, &op)
	}

	tr := img[2]
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.Min.X+pad+centersWide*stride), float64(r.Min.Y))
	screen.DrawImage(tr, &op)

	// middle rows
	ml := img[3]
	mc := img[4]
	mr := img[5]
	for j := 0; j < middlesTall; j++ {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y+pad+j*stride))
		screen.DrawImage(ml, &op)

		// center middle
		for i := 0; i < centersWide; i++ {
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(r.Min.X+pad+i*stride), float64(r.Min.Y+pad+j*stride))
			screen.DrawImage(mc, &op)
		}

		op = ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(r.Min.X+pad+centersWide*stride), float64(r.Min.Y+pad+j*stride))
		screen.DrawImage(mr, &op)

	}

	// bottom row
	bl := img[6]
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y+pad+middlesTall*stride))
	screen.DrawImage(bl, &op)

	bc := img[7]
	for i := 0; i < centersWide; i++ {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(r.Min.X+pad+i*stride), float64(r.Min.Y+pad+middlesTall*stride))
		screen.DrawImage(bc, &op)
	}

	br := img[8]
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.Min.X+pad+centersWide*stride), float64(r.Min.Y+pad+middlesTall*stride))
	screen.DrawImage(br, &op)
	return nil
}

func (uv *uiVisualizer) drawButton(screen *ebiten.Image, active bool, r image.Rectangle) error {
	img, err := uv.picForTexture("ui.png")
	if err != nil {
		return fmt.Errorf("picForTexture: %v", err)
	}
	imgs := [9]*ebiten.Image{
		img.SubImage(image.Rect(12, 0, 15, 3)).(*ebiten.Image),
		img.SubImage(image.Rect(15, 0, 18, 3)).(*ebiten.Image),
		img.SubImage(image.Rect(18, 0, 21, 3)).(*ebiten.Image),
		img.SubImage(image.Rect(12, 3, 15, 6)).(*ebiten.Image),
		img.SubImage(image.Rect(15, 3, 18, 6)).(*ebiten.Image),
		img.SubImage(image.Rect(18, 3, 21, 6)).(*ebiten.Image),
		img.SubImage(image.Rect(12, 6, 15, 9)).(*ebiten.Image),
		img.SubImage(image.Rect(15, 6, 18, 9)).(*ebiten.Image),
		img.SubImage(image.Rect(18, 6, 21, 9)).(*ebiten.Image),
	}
	if active {
		imgs = [9]*ebiten.Image{
			img.SubImage(image.Rect(21, 0, 24, 3)).(*ebiten.Image),
			img.SubImage(image.Rect(24, 0, 27, 3)).(*ebiten.Image),
			img.SubImage(image.Rect(27, 0, 30, 3)).(*ebiten.Image),
			img.SubImage(image.Rect(21, 3, 24, 6)).(*ebiten.Image),
			img.SubImage(image.Rect(24, 3, 27, 6)).(*ebiten.Image),
			img.SubImage(image.Rect(27, 3, 30, 6)).(*ebiten.Image),
			img.SubImage(image.Rect(21, 6, 24, 9)).(*ebiten.Image),
			img.SubImage(image.Rect(24, 6, 27, 9)).(*ebiten.Image),
			img.SubImage(image.Rect(27, 6, 30, 9)).(*ebiten.Image),
		}
	}

	if err := nineSlice(screen, r, imgs, 3, 3); err != nil {
		return fmt.Errorf("nineSlice: %v", err)
	}

	return nil
}

func (uv *uiVisualizer) drawText(screen *ebiten.Image, value string, size ui.TextSize, bounds image.Rectangle, align ui.TextLayout) (height int, err error) {
	text := ui.NewText(value, size)

	// Spacer around each text instance.
	spacer := ui.TextPadding

	// We know our bounds now, so we can split long lines.
	width := bounds.Dx()
	splitLines := ui.SplitLines(text.Lines, width)

	img, err := uv.picForTexture(text.BitmapFontTexture)
	if err != nil {
		return 0, fmt.Errorf("picForTexture: %v", err)
	}

	y := float64(bounds.Min.Y + spacer)
	for i, line := range splitLines {
		x := float64(bounds.Min.X)
		if i != 0 {
			// If not the first line, add a line spacer.
			y += float64(ui.LineSpacing(size))
		}

		// Different strategies based on width and word breaks...
		switch align {
		case ui.TextLayoutRight:
			x = float64(bounds.Max.X) - float64(line.Width())
		case ui.TextLayoutCenter:
			x += float64(bounds.Dx()/2) - float64(line.Width()/2)
		}

		tallest := 0
		wordSpace := float64(ui.SpaceWidth)
		if align == ui.TextLayoutJustify && len(line) > 1 {
			extra := float64((float64(bounds.Dx()) - float64(line.Width())) / float64(len(line)-1))
			wordSpace += extra
		}
		for _, word := range line {
			for i, char := range word.Characters {
				if char.Height > tallest {
					tallest = char.Height
				}
				img := img.SubImage(image.Rect(char.X, char.Y, char.X+char.Width, char.Y+char.Height)).(*ebiten.Image)
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(x, y)
				screen.DrawImage(img, &op)
				x += float64(char.Width)

				// Add spacing between letters for every letter except the last one.
				if i != len(word.Characters)-1 {
					x += float64(ui.LetterSpacing)
				}
			}
			x += wordSpace
		}

		y += float64(tallest)
	}

	return spacer + int(y) - bounds.Min.Y, nil
}
