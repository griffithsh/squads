package data

import "github.com/griffithsh/squads/game"

func (a *Archive) PedestalAppearances(sinister bool) []int {
	if sinister {
		return []int{2, 3}
	}
	return []int{1, 4}
}

func (a *Archive) GetPedestal(pedestalAppearance int) *game.Sprite {
	switch pedestalAppearance {
	case 2:
		return &game.Sprite{
			Texture: "combat/coasters.png",
			X:       48,
			W:       54,
			H:       42,
		}
	case 3:
		return &game.Sprite{
			Texture: "combat/coasters.png",
			X:       0,
			Y:       42,
			W:       39,
			H:       34,
		}
	case 4:
		return &game.Sprite{
			Texture: "combat/coasters.png",
			X:       0,
			Y:       76,
			W:       39,
			H:       34,
		}
	default:
		return &game.Sprite{
			Texture: "combat/coasters.png",
			W:       48,
			H:       40,
		}
	}
}
