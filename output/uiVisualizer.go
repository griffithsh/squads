package output

import (
	"fmt"
	"image"

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

func (uv *uiVisualizer) Render(screen *ebiten.Image, doc *ui.Element, data map[string]func(), scale float64) error {
	// FIXME: remove this random code that's only here to test stuff.
	// if err := uv.drawPanel(screen, image.Rect(62, 64, 962, 704), scale); err != nil {
	// 	return fmt.Errorf("drawPanel: %v", err)
	// }
	// if err := uv.drawButton(screen, true, image.Rect(70, 72, 190, 102), scale); err != nil {
	// 	return fmt.Errorf("drawButton: %v", err)
	// }

	// FIXME: figure out the boundaries for the children.
	boundaries := image.Rect(0, 0, 1024, 768)

	_, err := uv.drawChildren(screen, doc.Children, boundaries, "center", "middle", scale)

	return err
}

// there's clearly an entry point, and then something that can recurse through child elements ...
// But does there need to be a "new" that captures some of the basics?

// func (uv *uiVisualizer) draw(screen *ebiten.Image, src *ui.UI, position *game.Position, scale *game.Scale) error {
// 	// TODO
// 	// ui.Data
// 	xScale, yScale := 1.0, 1.0
// 	if scale != nil {
// 		xScale = scale.X
// 		yScale = scale.Y
// 	}
// 	w, err := strconv.Atoi(src.Doc.Attributes["width"])
// 	if err != nil {
// 		return fmt.Errorf("parse document width: %v", err)
// 	}
// 	p, ok := src.Doc.Attributes["padding"]
// 	if !ok {
// 		p = "0"
// 	}
// 	padding, err := strconv.Atoi(p)
// 	if err != nil {
// 		return fmt.Errorf("parse document padding: %v", err)
// 	}
// 	childWidth := w - (padding * 2)
// 	if src.Doc.Type == ui.PanelElement {

// 	}
// 	return nil
// }

// drawChildren returns the remainder of the bounds, unused by the drawn children.
func (uv *uiVisualizer) drawChildren(screen *ebiten.Image, children []*ui.Element, bounds image.Rectangle, align, valign string, scale float64) (image.Rectangle, error) {
	maxColHeight := 0
	for _, child := range children {
		var err error
		switch child.Type {
		case ui.PanelElement:
			w := child.Attributes.Width()
			h := child.Attributes.Height()
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

			if err = uv.drawPanel(screen, panelBounds, scale); err != nil {
				return bounds, err
			}

			padding := int(float64(child.Attributes.Padding()) * scale)
			panelBounds.Min.X += padding
			panelBounds.Min.Y += padding
			panelBounds.Max.X -= padding
			panelBounds.Max.Y -= padding
			if bounds, err = uv.drawChildren(screen, child.Children, panelBounds, child.Attributes.Align(), child.Attributes.Valign(), scale); err != nil {
				return bounds, err
			}

		case ui.RowElement:
			// I keep coming back to RowElement and wondering what its purpose was...?
			bounds, err = uv.drawChildren(screen, child.Children, bounds, child.Attributes.Align(), child.Attributes.Valign(), scale)
			if err != nil {
				return bounds, err
			}

		case ui.ColumnElement:
			// I think we need to know about siblings to do this correctly?
			// I don't think we can stomp bounds here?  Only the last Column of
			// adjacent siblings is block level.
			colBounds := bounds
			colBounds.Min.X += bounds.Dx() * child.Attributes.TwelfthsOffset() / 12
			w := bounds.Dx() * child.Attributes.Twelfths() / 12
			colBounds.Max.X = colBounds.Min.X + w
			takenBounds, err := uv.drawChildren(screen, child.Children, colBounds, child.Attributes.Align(), child.Attributes.Valign(), scale)
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

		case ui.TextElement:
			label := child.Attributes["value"]
			sz := child.Attributes.FontSize()
			layout := child.Attributes.FontLayout()

			txtBounds := bounds
			if child.Attributes["width"] != "" {
				txtBounds.Max.X = txtBounds.Min.X + child.Attributes.Width()
			}
			h, err := uv.drawText(screen, label, sz, txtBounds, layout, scale)
			if err != nil {
				return bounds, err
			}
			bounds.Min.Y += h

		case ui.ButtonElement:
			buttonHeight := int(15 * scale)
			label := child.Attributes["label"]
			width := child.Attributes.Width()
			// Does the parent align left, right, or centre? Are we valigning
			// it? Calculate buttonDimensions from that.
			l := bounds.Min.X
			switch align {
			case "right":
				l = bounds.Max.X - width
			case "center":
				l = bounds.Min.X + (bounds.Max.X-bounds.Min.X)/2 - width/2
			default: // left
			}
			t := bounds.Min.Y
			switch valign {
			case "bottom":
				t = bounds.Max.Y - buttonHeight
			case "middle":
				t = bounds.Min.Y + (bounds.Max.Y-bounds.Min.Y)/2 - buttonHeight/2
			default: // top
			}
			buttonDimensions := image.Rect(l, t, l+width, t+buttonHeight)
			if err = uv.drawButton(screen, false, buttonDimensions, scale); err != nil {
				return bounds, err
			}

			// Push the text down a bit so it's vertically centered.
			buttonDimensions.Min.Y += int(scale * 2)

			uv.drawText(screen, label, ui.TextSizeNormal, buttonDimensions, ui.TextLayoutCenter, scale)
			bounds.Min.Y += buttonHeight

		case ui.ImageElement:
			texture := child.Attributes["texture"]
			img, err := uv.picForTexture(texture)
			if err != nil {
				return bounds, fmt.Errorf("picForTexture: %v", err)
			}
			width := child.Attributes.Width()
			height := child.Attributes.Height()
			x := child.Attributes.X()
			y := child.Attributes.Y()
			img = img.SubImage(image.Rect(x, y, x+width, y+height)).(*ebiten.Image)
			op := ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(bounds.Min.X), float64(bounds.Min.Y))
			screen.DrawImage(img, &op)

			bounds.Min.Y += int(float64(height) * scale)
		}
	}
	return bounds, nil
}

// drawPanel renders a nine-slice "panel" that covers the provided Rectangle,
// and with a pixel granularity controlled by scale. The provided scale does not
// affect the portion of the screen covered by the panel.
func (uv *uiVisualizer) drawPanel(screen *ebiten.Image, r image.Rectangle, scale float64) error {
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

	return nineSlice(screen, r, scale, imgs, 4, 4)
}

func nineSlice(screen *ebiten.Image, r image.Rectangle, scale float64, img [9]*ebiten.Image, pad, stride int) error {
	scaleX := scale
	scaleY := scaleX

	// centersWide, middlesTall is how many copies of the center and middle
	// pieces of the 9-slice are required given this current scale.
	centersWide := int(float64(r.Dx()-int(float64(pad+pad)*scaleX)) / (float64(stride) * scaleX))
	middlesTall := int(float64(r.Dy()-int(float64(pad+pad)*scaleY)) / (float64(stride) * scaleY))

	// top row
	tl := img[0]
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	screen.DrawImage(tl, &op)

	tc := img[1]
	for i := 0; i < centersWide; i++ {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(r.Min.X)+float64(pad+i*stride)*scaleX, float64(r.Min.Y))
		screen.DrawImage(tc, &op)
	}

	tr := img[2]
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(r.Min.X)+float64(pad+centersWide*stride)*scaleX, float64(r.Min.Y))
	screen.DrawImage(tr, &op)

	// middle rows
	ml := img[3]
	mc := img[4]
	mr := img[5]
	for j := 0; j < middlesTall; j++ {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y)+float64(pad+j*stride)*scaleY)
		screen.DrawImage(ml, &op)

		// center middle
		for i := 0; i < centersWide; i++ {
			op := ebiten.DrawImageOptions{}
			op.GeoM.Scale(scaleX, scaleY)
			op.GeoM.Translate(float64(r.Min.X)+float64(pad+i*stride)*scaleX, float64(r.Min.Y)+float64(pad+j*stride)*scaleY)
			screen.DrawImage(mc, &op)
		}

		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(r.Min.X)+float64(pad+centersWide*stride)*scaleX, float64(r.Min.Y)+float64(pad+j*stride)*scaleY)
		screen.DrawImage(mr, &op)

	}

	// bottom row
	bl := img[6]
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y)+float64(pad+middlesTall*stride)*scaleY)
	screen.DrawImage(bl, &op)

	bc := img[7]
	for i := 0; i < centersWide; i++ {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(r.Min.X)+float64(pad+i*stride)*scaleX, float64(r.Min.Y)+float64(pad+middlesTall*stride)*scaleY)
		screen.DrawImage(bc, &op)
	}

	br := img[8]
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(r.Min.X)+float64(pad+centersWide*stride)*scaleX, float64(r.Min.Y)+float64(pad+middlesTall*stride)*scaleY)
	screen.DrawImage(br, &op)
	return nil
}

func (uv *uiVisualizer) drawButton(screen *ebiten.Image, active bool, r image.Rectangle, scale float64) error {
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

	if err := nineSlice(screen, r, scale, imgs, 3, 3); err != nil {
		return fmt.Errorf("nineSlice: %v", err)
	}

	return nil
}

func (uv *uiVisualizer) drawText(screen *ebiten.Image, value string, size ui.TextSize, bounds image.Rectangle, align ui.TextLayout, scale float64) (height int, err error) {
	text := ui.NewText(value, size)

	// Spacer around each text instance.
	spacer := int(1 * scale)

	// We know our bounds now, so we can split long lines.
	width := int(float64(bounds.Dx()) / scale)
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
			y += float64(ui.LineSpacing(size)) * scale
		}

		// Different strategies based on width and word breaks...
		switch align {
		case ui.TextLayoutRight:
			x = float64(bounds.Max.X) - float64(line.Width())*scale
		case ui.TextLayoutCenter:
			x += float64(bounds.Dx()/2) - float64(line.Width()/2)*scale
		}

		tallest := 0
		wordSpace := ui.SpaceWidth * scale
		if align == ui.TextLayoutJustify && len(line) > 1 {
			extra := float64((float64(bounds.Dx()) - float64(line.Width())*scale) / float64(len(line)-1))
			wordSpace += extra
		}
		for _, word := range line {
			for i, char := range word.Characters {
				if char.Height > tallest {
					tallest = char.Height
				}
				img := img.SubImage(image.Rect(char.X, char.Y, char.X+char.Width, char.Y+char.Height)).(*ebiten.Image)
				op := ebiten.DrawImageOptions{}
				op.GeoM.Scale(scale, scale)
				op.GeoM.Translate(x, y)
				screen.DrawImage(img, &op)
				x += float64(char.Width) * scale

				// Add spacing between letters for every letter except the last one.
				if i != len(word.Characters)-1 {
					x += float64(ui.LetterSpacing) * scale
				}
			}
			x += wordSpace
		}

		y += float64(tallest) * scale
	}

	return spacer + int(y) - bounds.Min.Y, nil
}
