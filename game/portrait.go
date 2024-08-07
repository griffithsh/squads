package game

// BigPortrait is the dimensions of the big portraits in the game.
var BigPortrait = 52

// SmallPortrait is the dimensions of the small portraits in the game.
var SmallPortrait = BigPortrait / 2

// PortraitBGBig stores a list of portrait backgrounds for large icons.
var PortraitBGBig = []Sprite{
	Sprite{
		Texture: "portrait-backgrounds.png", X: 0, Y: 0, W: BigPortrait, H: BigPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: BigPortrait, Y: 0, W: BigPortrait, H: BigPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: BigPortrait * 2, Y: 0, W: BigPortrait, H: BigPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: BigPortrait * 3, Y: 0, W: BigPortrait, H: BigPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: 0, Y: BigPortrait, W: BigPortrait, H: BigPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: BigPortrait, Y: BigPortrait, W: BigPortrait, H: BigPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: BigPortrait * 2, Y: BigPortrait, W: BigPortrait, H: BigPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: BigPortrait * 3, Y: BigPortrait, W: BigPortrait, H: BigPortrait,
	},
}

// PortraitBGSmall stores a list of portrait backgrounds for small icons.
var PortraitBGSmall = []Sprite{
	Sprite{
		Texture: "portrait-backgrounds.png", X: 0, Y: 208, W: SmallPortrait, H: SmallPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: SmallPortrait, Y: 208, W: SmallPortrait, H: SmallPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: SmallPortrait * 2, Y: 208, W: SmallPortrait, H: SmallPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: SmallPortrait * 3, Y: 208, W: SmallPortrait, H: SmallPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: 208, Y: 0, W: SmallPortrait, H: SmallPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: 208, Y: SmallPortrait, W: SmallPortrait, H: SmallPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: 208, Y: SmallPortrait * 2, W: SmallPortrait, H: SmallPortrait,
	},
	Sprite{
		Texture: "portrait-backgrounds.png", X: 208, Y: SmallPortrait * 3, W: SmallPortrait, H: SmallPortrait,
	},
}

// PortraitFrameBig stores a list of portrait frames for large icons.
var PortraitFrameBig = []Sprite{
	Sprite{
		Texture: "portrait-frames.png", X: 0, Y: 0, W: BigPortrait, H: BigPortrait,
	},
}

// PortraitFrameSmall stores a list of portrait frames for small icons.
var PortraitFrameSmall = []Sprite{
	Sprite{
		Texture: "portrait-frames.png", X: 0, Y: 208, W: SmallPortrait, H: SmallPortrait,
	},
}
