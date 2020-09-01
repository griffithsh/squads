package data

import "github.com/griffithsh/squads/game"

func (a *Archive) PedestalAppearances() []int {
	return []int{
		1,
	}
}

func (a *Archive) GetPedestal(pedestalAppearance int) *game.Sprite {
	return &game.Sprite{
		Texture: "coasters.png",
		W:       48,
		H:       40,
	}
}
