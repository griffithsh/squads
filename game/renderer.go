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
	textures map[string]*ebiten.Image
}

// NewRenderer creates a new Renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		textures: map[string]*ebiten.Image{},
	}
}

type entity struct {
	s      *Sprite
	p      *Position
	offset *RenderOffset
	scale  *Scale
	repeat *SpriteRepeat
}

// drawImageOptions creates an newebite.DrawImageOptions for this entity.
func (e entity) drawImageOptions(x, y, zoom, w, h, xOff, yOff float64) *ebiten.DrawImageOptions {
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
		// Scale the rendered entities based on the zoom value from the camera.
		// NB: This needs to happen after the other translations!
		op.GeoM.Scale(zoom, zoom)

		// We also need to correct for the dimensions of the screen, or the
		// focus will appear in the top-left corner of the screen. This comes
		// after the scaling, because the screen size does not change based on
		// the zoom.
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

		return entities[i].s.Texture < entities[j].s.Texture
	})

	return entities
}

func (r *Renderer) picForTexture(filename string) (*ebiten.Image, error) {
	if pic, ok := r.textures[filename]; ok {
		return pic, nil
	}

	inline, ok := res.Images[filename]
	if !ok {
		return nil, fmt.Errorf("%s missing", filename)
	}
	img, err := ebiten.NewImageFromImage(inline, ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("load texture: %v", err)
	}
	r.textures[filename] = img
	return img, nil
}

// Render all sprites in the world to the Target. We need to know where in the
// world we are focused, as well as how zoomed in we are, and the dimensions of
// the screen.
func (r *Renderer) Render(screen *ebiten.Image, mgr *ecs.World, focusX, focusY, zoom, screenW, screenH float64) error {
	entities := r.getEntities(mgr)

	screen.Fill(color.NRGBA{40, 34, 31, 0xff})
	for _, e := range entities {
		img, err := r.picForTexture(e.s.Texture)
		if err != nil {
			return fmt.Errorf("get texture: %v", err)
		}

		img, ok := img.SubImage(image.Rectangle{image.Point{e.s.X, e.s.Y}, image.Point{e.s.X + e.s.W, e.s.Y + e.s.H}}).(*ebiten.Image)
		if !ok {
			return fmt.Errorf("SubImage %s: invalid type cast", e.s.Texture)
		}

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

					op := e.drawImageOptions(focusX, focusY, zoom, screenW, screenH, offX, offY)
					screen.DrawImage(img, op)
				}
			}
			if err != nil {
				return fmt.Errorf("SpriteRepeat: %v", err)
			}
			continue
		}

		op := e.drawImageOptions(focusX, focusY, zoom, screenW, screenH, 0, 0)
		if err := screen.DrawImage(img, op); err != nil {
			return fmt.Errorf("DrawImage: %v", err)
		}
	}

	return nil
}
