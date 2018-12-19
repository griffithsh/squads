package game

import (
	"fmt"
	"image"
	"os"
	"sort"

	"github.com/faiface/pixel"
	"github.com/griffithsh/squads/ecs"
)

// Renderer is a System that draws positioned Sprites to the screen.
type Renderer struct {
	textures map[string]pixel.Picture
}

// NewRenderer creates a new Renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		textures: map[string]pixel.Picture{},
	}
}

type entity struct {
	s *Sprite
	p *Position
}

func (r *Renderer) getEntities(mgr *ecs.World) []entity {
	raw := mgr.Get([]string{"Sprite", "Position"})
	entities := make([]entity, len(raw))
	for i, e := range raw {
		entities[i] = entity{
			s: mgr.Component(e, "Sprite").(*Sprite),
			p: mgr.Component(e, "Position").(*Position),
		}
	}

	// sort by position layer, position Z, position.Y - sprite.Y/2
	sort.Slice(entities, func(i, j int) bool {
		if entities[i].p.Layer != entities[j].p.Layer {
			return entities[i].p.Layer < entities[j].p.Layer
		}

		iExtent := entities[i].p.Center.Y + float64(entities[i].s.H)
		jExtent := entities[j].p.Center.Y + float64(entities[j].s.H)
		if iExtent != jExtent {
			return iExtent < jExtent
		}

		return entities[i].s.Texture < entities[j].s.Texture
	})

	return entities
}

func (r *Renderer) picForTexture(tex string) (pixel.Picture, error) {
	if pic, ok := r.textures[tex]; ok {
		return pic, nil
	}

	pic, err := loadPicture(tex)
	if err != nil {
		return nil, fmt.Errorf("loadPicture: %v", err)
	}
	r.textures[tex] = pic
	return pic, nil
}

// Render all sprites in the world to the Target.
func (r *Renderer) Render(win pixel.Target, cam pixel.Matrix, mgr *ecs.World) error {
	var batch *pixel.Batch

	entities := r.getEntities(mgr)

	var lastTexture string

	var pic pixel.Picture
	for _, e := range entities {
		if e.s.Texture != lastTexture {
			p, err := r.picForTexture(e.s.Texture)
			if err != nil {
				return fmt.Errorf("picForTexture: %v", err)
			}

			if batch != nil {
				batch.Draw(win)
			}
			batch = pixel.NewBatch(&pixel.TrianglesData{}, p)
			batch.SetMatrix(cam)
			pic = p
		}

		height := pic.Bounds().H()
		sprite := pixel.NewSprite(pic, pixel.Rect{
			Min: pixel.Vec{
				X: float64(e.s.X),
				Y: height - float64(e.s.Y),
			},
			Max: pixel.Vec{
				X: float64(e.s.X + e.s.W),
				Y: height - float64(e.s.H+e.s.Y),
			},
		})

		// faiface/pixel inverts the Y coordinate
		y := -e.p.Center.Y

		move := pixel.Vec{X: e.p.Center.X, Y: y}

		if e.s.Color != nil {
			sprite.DrawColorMask(batch, pixel.IM.Moved(move), e.s.Color)
		} else {
			sprite.Draw(batch, pixel.IM.Moved(move))
		}

		lastTexture = e.s.Texture
	}
	if batch != nil {
		batch.Draw(win)
	}

	return nil
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
