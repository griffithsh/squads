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
			if bounds, err = uv.drawChildren(screen, child.Children, panelBounds, child.Attributes.Align(), child.Attributes.Valign(), scale); err != nil {
				return bounds, err
			}

		case ui.RowElement:
			bounds, err = uv.drawChildren(screen, child.Children, bounds, child.Attributes.Align(), child.Attributes.Valign(), scale)
			if err != nil {
				return bounds, err
			}

		case ui.ColumnElement:
			// I think we need to know about siblings to do this correctly?
			// I don't think we can stomp bounds here?  Only the last Column of
			// adjacent siblings is block level.
			bounds, err = uv.drawChildren(screen, child.Children, bounds, child.Attributes.Align(), child.Attributes.Valign(), scale)
			if err != nil {
				return bounds, err
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
			bounds.Min.Y += int(float64(h) * scale)

		case ui.ButtonElement:
			const buttonHeight = 15
			label := child.Attributes["label"]
			width := child.Attributes.Width()
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
			// does the parent align left, right, or centre? Are we valigning it? Calculate a bounds from that.
			if err = uv.drawButton(screen, false, buttonDimensions, scale); err != nil {
				return bounds, err
			}
			uv.drawText(screen, label, ui.TextSizeNormal, buttonDimensions, ui.TextLayoutCenter, scale)
			bounds.Min.Y += buttonHeight

		case ui.ImageElement:
			// TODO: bigtime!
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

func split(lines []ui.Line, width int) []ui.Line {
	splitLines := []ui.Line{}

	for _, line := range lines {
		// If this line will fit within bounds, then it doesnt need splitting.
		if line.Width() <= width {
			splitLines = append(splitLines, line)
			continue
		}

		// Go through the line, cutting it into pieces that will fit.
		running := ui.Line{}
		for _, word := range line {
			// What if this word alone exceeds bounds?
			// FIXME:

			if running.Width()+ui.LetterSpacing+word.Width() > width {
				splitLines = append(splitLines, running)
				running = running[:0]
			}
			running = append(running, word)
		}
		splitLines = append(splitLines, running)
	}

	return splitLines
}

func (uv *uiVisualizer) drawText(screen *ebiten.Image, value string, size ui.TextSize, bounds image.Rectangle, align ui.TextLayout, scale float64) (height int, err error) {
	text := ui.NewText(value, size)

	// We know our bounds now, so we can split long lines.
	splitLines := split(text.Lines, bounds.Dx())

	img, err := uv.picForTexture(text.BitmapFontTexture)
	if err != nil {
		return 0, fmt.Errorf("picForTexture: %v", err)
	}

	// height is a return value.
	height = 0

	x := float64(bounds.Min.X)
	y := float64(bounds.Min.Y)
	for _, line := range splitLines {
		tallest := 0
		// different strategies based on width and word breaks...
		switch align {
		case ui.TextLayoutLeft:
			// flow from bounds.Min.X
		case ui.TextLayoutRight:
			// flow from bounds.Max.X - line.Width()
		case ui.TextLayoutCenter:
			// flow from ...?
		case ui.TextLayoutJustify:
			// flow from left, but divide extra space between words...
		}

		// FIXME: Clean up and share this implementation across layout styles.
		for _, word := range line {
			for _, char := range word.Characters {
				if char.Height > tallest {
					tallest = char.Height
				}
				img := img.SubImage(image.Rect(char.X, char.Y, char.X+char.Width, char.Y+char.Height)).(*ebiten.Image)
				op := ebiten.DrawImageOptions{}
				op.GeoM.Scale(scale, scale)
				op.GeoM.Translate(x, y)
				screen.DrawImage(img, &op)
				x += float64(char.Width+ui.LetterSpacing) * scale
			}
		}
		height += int(float64(tallest) * scale)
		y += float64(tallest+ui.LineSpacing(size)) * scale
	}

	if len(splitLines) > 0 {
		height += int(float64(ui.LineSpacing(size)*(len(splitLines)-1)) * scale)
	}
	return height, nil
}
