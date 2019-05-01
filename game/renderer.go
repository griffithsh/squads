package game

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/griffithsh/squads/ecs"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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
	offset *SpriteOffset
}

// getEntities returns a sorted list of entities that have renderable
// components.
// FIXME: getEntities should be refactored so that non-Sprite renderable
// components can also be returned.
func (r *Renderer) getEntities(mgr *ecs.World) []entity {
	raw := mgr.Get([]string{"Sprite", "Position"})
	entities := make([]entity, len(raw))
	for i, e := range raw {
		entities[i] = entity{
			s: mgr.Component(e, "Sprite").(*Sprite),
			p: mgr.Component(e, "Position").(*Position),
		}
		if offset, ok := mgr.Component(e, "SpriteOffset").(*SpriteOffset); ok {
			entities[i].offset = offset
		}
	}

	// sort by position layer, position.Y - sprite.Y/2
	sort.Slice(entities, func(i, j int) bool {
		if entities[i].p.Layer != entities[j].p.Layer {
			return entities[i].p.Layer < entities[j].p.Layer
		}

		iExtent := entities[i].p.Center.Y + float64(entities[i].s.H/2)
		jExtent := entities[j].p.Center.Y + float64(entities[j].s.H/2)
		if entities[i].offset != nil {
			iExtent += float64(entities[i].offset.Y)
		}
		if entities[j].offset != nil {
			jExtent += float64(entities[j].offset.Y)
		}
		if iExtent != jExtent {
			return iExtent < jExtent
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

	img, _, err := ebitenutil.NewImageFromFile(filename, ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("load texture: %v", err)
	}
	r.textures[filename] = img
	return img, nil
}

// Render all sprites in the world to the Target. We need to know where in the
// world we are focused, as well as how zoomed in we are.
func (r *Renderer) Render(screen *ebiten.Image, x, y, zoom, w, h float64, mgr *ecs.World) error {
	entities := r.getEntities(mgr)

	screen.Fill(color.NRGBA{40, 34, 31, 0xff})
	for _, e := range entities {
		op := &ebiten.DrawImageOptions{}

		// Translate for the focus values from the camera
		op.GeoM.Translate(-x, -y)

		// Translate for the location of the Entity
		op.GeoM.Translate(e.p.Center.X, e.p.Center.Y)

		// ebiten uses top-left corner coordinates, so we need to translate
		// from center-based coordinates by subtracting half the width/height.
		op.GeoM.Translate(-0.5*float64(e.s.W), -0.5*float64(e.s.H))

		// Some sprites may need to be drawn with an offset.
		if e.offset != nil {
			op.GeoM.Translate(float64(e.offset.X), float64(e.offset.Y))
		}

		// Scale the rendered entities based on the zoom value
		// NB: This needs to happen after the other translations!
		op.GeoM.Scale(zoom, zoom)

		// We also need to correct for the dimensions of the screen, or the
		// focus will appear in the top-left corner of the screen. This comes
		// after the scaling, because the screen size does not change based on
		// the zoom.
		op.GeoM.Translate(w/2, h/2)

		img, err := r.picForTexture(e.s.Texture)
		if err != nil {
			return fmt.Errorf("get texture: %v", err)
		}

		img, ok := img.SubImage(image.Rectangle{image.Point{e.s.X, e.s.Y}, image.Point{e.s.X + e.s.W, e.s.Y + e.s.H}}).(*ebiten.Image)
		if !ok {
			return fmt.Errorf("SubImage %s: invalid type cast", e.s.Texture)
		}

		if err := screen.DrawImage(img, op); err != nil {
			return fmt.Errorf("DrawImage: %v", err)
		}
	}

	return nil
}
