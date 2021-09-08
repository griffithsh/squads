package output

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sort"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/res"
	"github.com/griffithsh/squads/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

// Visualizer is a System that draws visible things to the screen with the ebiten game engine.
type Visualizer struct {
	textures      map[string]*ebiten.Image
	imageProvider ImageProvider

	worldCanvas    *ebiten.Image
	absoluteCanvas *ebiten.Image
	uiCanvas       *ebiten.Image

	// randFloats are used to define whether a tile of an Entity with an Alpha
	// component should be shown or not.
	randFloats []float64

	uv *uiVisualizer
}

// ImageProvider is the contract that the Visualizer needs to retrieve images to
// use for Sprites.
type ImageProvider interface {
	GetImage(name string) (val image.Image, ok bool)
}

// NewVisualizer creates a new Visualizer.
func NewVisualizer(images ImageProvider) *Visualizer {
	// 4099 is the first prime number after 4096. Prime numbers have no
	// factors, so they won't repeat a pattern on power of two sized textures.
	// 4096 is 64*64. 64 is 512/8. 512x512 is a reasonably large texture size,
	// and 8 is the width/height of one of the tiles.
	pre := make([]float64, 4099)
	for i := 0; i < len(pre); i++ {
		pre[i] = rand.Float64()
	}

	v := Visualizer{
		textures:      map[string]*ebiten.Image{},
		imageProvider: images,
		worldCanvas:   ebiten.NewImage(1, 1),
		uiCanvas:      ebiten.NewImage(1, 1),
		randFloats:    pre,
	}
	v.uv = newUIVisualizer(v.picForTexture)
	return &v
}

type entity struct {
	e      ecs.Entity
	s      *game.Sprite
	p      *game.Position
	offset *game.RenderOffset
	scale  *game.Scale
	repeat *game.SpriteRepeat
	alpha  *game.Alpha
}

// drawImageOptions creates an new ebiten.DrawImageOptions for this entity.
func (e entity) drawImageOptions(x, y, w, h, xOff, yOff float64) *ebiten.DrawImageOptions {
	op := ebiten.DrawImageOptions{}

	// ebiten uses top-left corner coordinates, so we need to translate
	// from center-based coordinates by subtracting half the width/height.
	op.GeoM.Translate(-float64(e.s.W/2), -float64(e.s.H/2))

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
func (r *Visualizer) getEntities(mgr *ecs.World) []entity {
	raw := mgr.Get([]string{"Sprite", "Position"})

	entities := make([]entity, 0, len(raw))
	for _, e := range raw {
		// Filter out any Hidden Entities.
		if mgr.Component(e, "Hidden") != nil {
			continue
		}

		ent := entity{
			e: e,
			p: mgr.Component(e, "Position").(*game.Position),
		}

		// Don't return entities with alpha set to zero.
		alpha, ok := mgr.Component(e, "Alpha").(*game.Alpha)
		if ok {
			if alpha.Value == 0 {
				continue
			}
			ent.alpha = alpha
		}

		ent.s = mgr.Component(e, "Sprite").(*game.Sprite)
		// Don't attempt to render sprites without a Texture.
		if ent.s.Texture == "" {
			continue
		}

		// TODO: Don't bother including anything that's outside the visible
		// area. Ebiten does not perform frustrum culling.

		if offset, ok := mgr.Component(e, "RenderOffset").(*game.RenderOffset); ok {
			ent.offset = offset
		}
		if scale, ok := mgr.Component(e, "Scale").(*game.Scale); ok {
			ent.scale = scale
		}
		if repeat, ok := mgr.Component(e, "SpriteRepeat").(*game.SpriteRepeat); ok {
			ent.repeat = repeat
		}
		if ok {
			ent.alpha = alpha
		}
		entities = append(entities, ent)
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

		ia := 1.0
		if entities[i].alpha != nil {
			ia = entities[i].alpha.Value
		}
		ja := 1.0
		if entities[j].alpha != nil {
			ja = entities[j].alpha.Value
		}
		if ia != ja {
			return ia < ja
		}
		return entities[i].e < entities[j].e
	})

	return entities
}

func (r *Visualizer) picForTexture(filename string) (*ebiten.Image, error) {
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
	img := ebiten.NewImageFromImage(inline)
	r.textures[filename] = img
	return img, nil
}

// ensureCanvases makes sure the canvases used for offscreen rendering each
// frame are the correct size based on the current screen dimensions and zoom.
// FIXME: this would work better with some sort of event subscription model I think?
func (r *Visualizer) ensureCanvases(w, h, z, uiScale float64) {
	b := r.worldCanvas.Bounds()
	// Have bounds changed?
	if b.Max.X-b.Min.X != int(w/z) || b.Max.Y-b.Min.Y != int(h/z) {
		r.worldCanvas = ebiten.NewImage(int(w/z), int(h/z))

		// NB: Absolute elements are not zoomed by z, which is the world zoom here.
		r.absoluteCanvas = ebiten.NewImage(int(w), int(h))
	}

	b = r.uiCanvas.Bounds()
	if b.Max.X-b.Min.X != int(w/uiScale) || b.Max.Y-b.Min.Y != int(h/uiScale) {
		r.uiCanvas = ebiten.NewImage(int(w/uiScale), int(h/uiScale))
	}

}

// target returns the correct render target to use for this entity.
func (r *Visualizer) target(e entity) *ebiten.Image {
	result := r.worldCanvas
	if e.p.Absolute {
		result = r.absoluteCanvas
	}
	return result
}

func (r *Visualizer) renderRepeatedEntity(e entity, img *ebiten.Image, target *ebiten.Image, focusX, focusY, zoom, screenW, screenH float64) {
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
}

// renderEntity renders a single entity to the appropriate render target in the Visualizer.
func (r *Visualizer) renderEntity(e entity, focusX, focusY, zoom, screenW, screenH float64) error {
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
		r.renderRepeatedEntity(e, img, target, focusX, focusY, zoom, screenW, screenH)
		return nil
	}

	// FIXME: if an entity is repeated and has an alpha the alpha will ignored!
	if e.alpha != nil {
		const tileDim = 8
		w, h := img.Size()
		for y := 0; y < h; y += tileDim {
			for x := 0; x < w; x += tileDim {
				rng := r.randFloats[(x/tileDim+y/tileDim*64)%len(r.randFloats)]
				show := (rng < e.alpha.Value)
				if !show {
					continue
				}
				tile, _ := img.SubImage(image.Rectangle{image.Point{e.s.X + x, e.s.Y + y}, image.Point{e.s.X + x + tileDim, e.s.Y + y + tileDim}}).(*ebiten.Image)

				offX, offY := float64(x), float64(y)

				op := e.drawImageOptions(focusX, focusY, screenW/zoom, screenH/zoom, offX, offY)
				target.DrawImage(tile, op)
			}
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
func (r *Visualizer) Render(screen *ebiten.Image, mgr *ecs.World, focusX, focusY, zoom, screenW, screenH float64) error {
	uiScale := 2.0 // FIXME: should this be tied to zoom? Or separate?
	r.ensureCanvases(screenW, screenH, zoom, uiScale)

	entities := r.getEntities(mgr)

	r.worldCanvas.Fill(color.NRGBA{40, 34, 31, 0xff})
	r.absoluteCanvas.Clear()
	r.uiCanvas.Clear()

	for _, e := range entities {
		err := r.renderEntity(e, focusX, focusY, zoom, screenW, screenH)
		if err != nil {
			return err
		}
	}

	for _, e := range mgr.Get([]string{"UI"}) {
		uic := mgr.Component(e, "UI").(*ui.UI)
		if err := r.uv.Render(r.uiCanvas, uic); err != nil {
			return err
		}
	}

	screen.Clear()

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(zoom, zoom)
	screen.DrawImage(r.worldCanvas, &op)

	op = ebiten.DrawImageOptions{}
	screen.DrawImage(r.absoluteCanvas, &op)

	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(uiScale, uiScale)
	screen.DrawImage(r.uiCanvas, &op)

	return nil
}
