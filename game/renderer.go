package game

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/res"
	"github.com/hajimehoshi/ebiten"
)

// Renderer is a System that draws world-positioned Sprites to the screen.
type Renderer struct {
	textures      map[string]*ebiten.Image
	imageProvider ImageProvider

	worldCanvas *ebiten.Image
	uiCanvas    *ebiten.Image
}

// ImageProvider is the contract that the Renderer needs to retrieve images to
// use for Sprites.
type ImageProvider interface {
	GetImage(name string) (val image.Image, ok bool)
}

// NewRenderer creates a new Renderer.
func NewRenderer(images ImageProvider) *Renderer {
	c, _ := ebiten.NewImage(1, 1, ebiten.FilterNearest)
	return &Renderer{
		textures:      map[string]*ebiten.Image{},
		imageProvider: images,
		worldCanvas:   c,
	}
}

type entity struct {
	s      *Sprite
	p      *Position
	offset *RenderOffset
	scale  *Scale
	repeat *SpriteRepeat
	alpha  *Alpha
}

// drawImageOptions creates an new ebiten.DrawImageOptions for this entity.
func (e entity) drawImageOptions(x, y, w, h, xOff, yOff float64) *ebiten.DrawImageOptions {
	op := ebiten.DrawImageOptions{}

	// ebiten uses top-left corner coordinates, so we need to translate
	// from center-based coordinates by subtracting half the width/height.
	op.GeoM.Translate(-0.5*float64(e.s.W), -0.5*float64(e.s.H))

	// Some Entities might have an intrinsic scale.
	if e.scale != nil {
		op.GeoM.Scale(e.scale.X, e.scale.Y)
	}

	// If we are dealing in world-coordinates, then translate for the values
	// of what the camera is focused on.
	if !e.p.Absolute {
		op.GeoM.Translate(-x, -y)
	}

	// Translate for the location of the Entity
	op.GeoM.Translate(e.p.Center.X, e.p.Center.Y)

	// Apply passed offset.
	op.GeoM.Translate(xOff, yOff)

	// Some Sprites may have an offset configured.
	op.GeoM.Translate(float64(e.s.OffsetX), float64(e.s.OffsetY))

	// There could also be an offset configured at the Entity level.
	if e.offset != nil {
		op.GeoM.Translate(float64(e.offset.X), float64(e.offset.Y))
	}

	if !e.p.Absolute {
		// We need to correct for the dimensions of the screen, or the
		// focus will appear in the top-left corner of the screen.
		op.GeoM.Translate(w/2, h/2)
	}
	return &op
}

// getEntities returns a sorted list of entities that have renderable
// components.
// FIXME: getEntities should be refactored so that non-Sprite renderable
// components can also be returned.
func (r *Renderer) getEntities(mgr *ecs.World) []entity {
	raw := mgr.Get([]string{"Sprite", "Position"})

	entities := make([]entity, 0, len(raw))
	for _, e := range raw {
		// Filter out any Hidden Entities.
		if mgr.Component(e, "Hidden") != nil {
			continue
		}
		// Don't attempt to render sprites without a Texture.
		sprite := mgr.Component(e, "Sprite").(*Sprite)
		if sprite.Texture == "" {
			continue
		}

		// TODO: Don't bother including anything that's outside the visible
		// area. Ebiten does not perform frustrum culling.

		entities = append(entities, entity{
			s: sprite,
			p: mgr.Component(e, "Position").(*Position),
		})
		if offset, ok := mgr.Component(e, "RenderOffset").(*RenderOffset); ok {
			entities[len(entities)-1].offset = offset
		}
		if scale, ok := mgr.Component(e, "Scale").(*Scale); ok {
			entities[len(entities)-1].scale = scale
		}
		if repeat, ok := mgr.Component(e, "SpriteRepeat").(*SpriteRepeat); ok {
			entities[len(entities)-1].repeat = repeat
		}
		if alpha, ok := mgr.Component(e, "Alpha").(*Alpha); ok {
			entities[len(entities)-1].alpha = alpha
		}
	}

	// sort by position layer, position.Y - sprite.Y/2
	sort.Slice(entities, func(i, j int) bool {
		if entities[i].p.Layer != entities[j].p.Layer {
			return entities[i].p.Layer < entities[j].p.Layer
		}

		iExtent := entities[i].p.Center.Y + float64(entities[i].s.H/2)
		jExtent := entities[j].p.Center.Y + float64(entities[j].s.H/2)

		iExtent += float64(entities[i].s.OffsetY)
		jExtent += float64(entities[j].s.OffsetY)

		ix := int(entities[i].p.Center.X)
		jx := int(entities[j].p.Center.X)

		if entities[i].offset != nil {
			iExtent += float64(entities[i].offset.Y)
			ix += entities[i].offset.X
		}
		if entities[j].offset != nil {
			jExtent += float64(entities[j].offset.Y)
			jx += entities[j].offset.X
		}
		if iExtent != jExtent {
			return iExtent < jExtent
		}

		if ix != jx {
			return ix < jx
		}

		// TODO: also sort by color? See https://github.com/hajimehoshi/ebiten/wiki/Performance-Tips

		if entities[i].s.Texture != entities[j].s.Texture {
			return entities[i].s.Texture < entities[j].s.Texture
		}

		return entities[i].alpha.Value < entities[j].alpha.Value
	})

	return entities
}

func (r *Renderer) picForTexture(filename string) (*ebiten.Image, error) {
	if pic, ok := r.textures[filename]; ok {
		return pic, nil
	}

	inline, ok := r.imageProvider.GetImage(filename)
	if !ok {
		inline, ok = res.Images[filename]
		if !ok {
			return nil, fmt.Errorf("%s missing", filename)
		}
	}
	img, err := ebiten.NewImageFromImage(inline, ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("load texture: %v", err)
	}
	r.textures[filename] = img
	return img, nil
}

// ensureCanvases makes sure the canvases used for offscreen rendering each
// frame are the correct size based on the current screen dimensions and zoom.
func (r *Renderer) ensureCanvases(w, h, z float64) {
	b := r.worldCanvas.Bounds()
	if b.Max.X-b.Min.X != int(w/z) || b.Max.Y-b.Min.Y != int(h/z) {
		r.worldCanvas, _ = ebiten.NewImage(int(w/z), int(h/z), ebiten.FilterNearest)

		// NB: UI elements are not zoomed by z, which is the world zoom here.
		r.uiCanvas, _ = ebiten.NewImage(int(w), int(h), ebiten.FilterNearest)
	}
}

// target returns the correct render target to use for this entity.
func (r *Renderer) target(e entity) *ebiten.Image {
	result := r.worldCanvas
	if e.p.Absolute {
		result = r.uiCanvas
	}
	return result
}

// renderEntity renders a single entity to the appropriate render target in the Renderer.
func (r *Renderer) renderEntity(e entity, focusX, focusY, zoom, screenW, screenH float64) error {
	img, err := r.picForTexture(e.s.Texture)
	if err != nil {
		return fmt.Errorf("get texture: %v", err)
	}

	img, ok := img.SubImage(image.Rectangle{image.Point{e.s.X, e.s.Y}, image.Point{e.s.X + e.s.W, e.s.Y + e.s.H}}).(*ebiten.Image)
	if !ok {
		return fmt.Errorf("SubImage %s: invalid type cast", e.s.Texture)
	}
	target := r.target(e)
	if e.repeat != nil {
		tileW, tileH := img.Size()

		for y := 0; y < e.repeat.H/tileH+1; y++ {
			if y == e.repeat.H/tileH {
				// last
				remainder := e.repeat.H % tileH
				if remainder == 0 {
					continue
				}
				// FIXME: this is unfinished and not usable for cicada
				// layers with non-evenly divisible bounds.
				// TODO: cut off the bottom of the img for the last row,
				// then continue into the for x ...
				// TODO: set tileH to be what's in remainder ...
			}

			for x := 0; x < e.repeat.W/tileW+1; x++ {
				if x == e.repeat.W/tileW {
					// last
					remainder := e.repeat.W % tileW
					if remainder == 0 {
						continue
					}
					// FIXME: this is unfinished and not usable for cicada
					// layers with non-evenly divisible bounds.
					// TODO: cut off the right of the img for the last
					// column, then continue onto the DrawImage call ...
					// TODO: set tileW to be what's in the remainder ...
				}
				offX := float64(x*tileW) - float64(e.repeat.W)/2 + float64(tileW)/2
				offY := float64(y*tileH) - float64(e.repeat.H)/2 + float64(tileH)/2

				op := e.drawImageOptions(focusX, focusY, screenW/zoom, screenH/zoom, offX, offY)
				target.DrawImage(img, op)
			}
		}
		if err != nil {
			return fmt.Errorf("SpriteRepeat: %v", err)
		}
		return nil
	}

	op := e.drawImageOptions(focusX, focusY, screenW/zoom, screenH/zoom, 0, 0)
	target.DrawImage(img, op)
	return nil
}

// Render all sprites in the world to the screen. We need to know where in the
// world we are focused, as well as how zoomed in we are, and the dimensions of
// the screen.
func (r *Renderer) Render(screen *ebiten.Image, mgr *ecs.World, focusX, focusY, zoom, screenW, screenH float64) error {
	r.ensureCanvases(screenW, screenH, zoom)

	entities := r.getEntities(mgr)

	r.worldCanvas.Fill(color.NRGBA{40, 34, 31, 0xff})
	r.uiCanvas.Clear()

	for _, e := range entities {
		err := r.renderEntity(e, focusX, focusY, zoom, screenW, screenH)
		if err != nil {
			return err
		}
	}

	screen.Clear()

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(zoom, zoom)
	screen.DrawImage(r.worldCanvas, &op)

	op = ebiten.DrawImageOptions{}
	screen.DrawImage(r.uiCanvas, &op)

	return nil
}
