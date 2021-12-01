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

func (uv *uiVisualizer) Render(screen *ebiten.Image, uic *ui.UI) error {
	for _, instruction := range uic.RenderInstructions() {
		switch t := instruction.(type) {
		case ui.ImageRenderInstruction:
			img, err := uv.picForTexture(t.Texture)
			if err != nil {
				return fmt.Errorf("picForTexture: %v", err)
			}
			img = img.SubImage(t.From).(*ebiten.Image)
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(t.AtX, t.AtY)
			screen.DrawImage(img, &op)

		case ui.PanelRenderInstruction:
			if err := uv.drawPanel(screen, t.Bounds); err != nil {
				return err
			}

		case ui.ButtonRenderInstruction:
			if err := uv.drawButton(screen, false, t.Bounds); err != nil {
				return err
			}
			txtBounds := t.Bounds
			txtBounds.Min.Y += 2
			if _, err := uv.drawText(screen, t.Label, ui.TextSizeNormal, txtBounds, ui.TextLayoutCenter); err != nil {
				return err
			}

		case ui.TextRenderInstruction:
			if _, err := uv.drawText(screen, t.Text, t.Size, t.Bounds, t.Layout); err != nil {
				return err
			}
		}
	}

	return nil
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
