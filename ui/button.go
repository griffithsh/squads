package ui

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
)

// Button creates a button
func Button(mgr *ecs.World, w, h int, l, t float64, layer int, absolute bool, text string, interact func(float64, float64)) ecs.Entity {
	ebg := ButtonBackground(mgr, w, h, l, t, layer, absolute)
	e := mgr.NewEntity()
	mgr.Dependency(ebg, e)

	mgr.AddComponent(e, &game.Font{
		Text: text,
		Size: "normal",
	})
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(w/2),
			Y: t + float64(h/2),
		},
		Layer:    layer + 1,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &Interactive{
		W:       float64(w),
		H:       float64(h),
		Trigger: interact,
	})

	return ebg
}

// ButtonBackground constructs a new background for a button. It visual only and contains no interaction.
func ButtonBackground(mgr *ecs.World, w, h int, l, t float64, layer int, absolute bool) ecs.Entity {
	const tileDimension int = 3
	// If w or h are not a multiple of tileDimension or less than 3 *
	// tileDimension, it is a programming error.
	if w%tileDimension != 0 || w < 3 || h%tileDimension != 0 || h < 3 {
		panic("Incorrect w/h for ui.Panel")
	}

	container := mgr.NewEntity()

	var e ecs.Entity

	// top left corner
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(tileDimension)/2,
			Y: t + float64(tileDimension)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 12, Y: 0,
		W: tileDimension, H: tileDimension,
	})

	// top centre
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(w)/2,
			Y: t + float64(tileDimension)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 15, Y: 0,
		W: tileDimension, H: tileDimension,
	})
	mgr.AddComponent(e, &game.SpriteRepeat{
		W: w - tileDimension*2,
		H: tileDimension,
	})

	// top right
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(w) - float64(tileDimension)/2,
			Y: t + float64(tileDimension)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 18, Y: 0,
		W: tileDimension, H: tileDimension,
	})

	// middle row, left side
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(tileDimension)/2,
			Y: t + float64(h)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 12, Y: 3,
		W: tileDimension, H: tileDimension,
	})
	mgr.AddComponent(e, &game.SpriteRepeat{
		W: tileDimension,
		H: h - tileDimension*2,
	})

	// middle row, center
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(w)/2,
			Y: t + float64(h)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 15, Y: 3,
		W: tileDimension, H: tileDimension,
	})
	mgr.AddComponent(e, &game.SpriteRepeat{
		W: w - tileDimension*2,
		H: h - tileDimension*2,
	})

	// middle row, right
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(w) - float64(tileDimension)/2,
			Y: t + float64(h)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 18, Y: 3,
		W: tileDimension, H: tileDimension,
	})
	mgr.AddComponent(e, &game.SpriteRepeat{
		W: tileDimension,
		H: h - tileDimension*2,
	})

	// bottom left
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(tileDimension)/2,
			Y: t + float64(h) - float64(tileDimension)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 12, Y: 6,
		W: tileDimension, H: tileDimension,
	})

	// bottom row, center
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(w)/2,
			Y: t + float64(h) - float64(tileDimension)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 15, Y: 6,
		W: tileDimension, H: tileDimension,
	})
	mgr.AddComponent(e, &game.SpriteRepeat{
		W: w - tileDimension*2,
		H: tileDimension,
	})

	// bottom row, center
	e = mgr.NewEntity()
	mgr.Dependency(container, e)
	mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: l + float64(w) - float64(tileDimension)/2,
			Y: t + float64(h) - float64(tileDimension)/2,
		},
		Layer:    layer,
		Absolute: absolute,
	})
	mgr.AddComponent(e, &game.Sprite{
		Texture: "ui.png",

		X: 18, Y: 6,
		W: tileDimension, H: tileDimension,
	})

	return container
}
