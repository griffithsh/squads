package data

import "github.com/griffithsh/squads/game"

// Profession returns the details of a profession.
func (a *Archive) Profession(profession string) *game.ProfessionDetails {
	switch profession {
	case "Wolf":
		return &game.ProfessionDetails{
			ActionPoints: 40,
			Preparation:  400,
		}
	case "Skeleton":
		return &game.ProfessionDetails{
			ActionPoints: 40,
			Preparation:  900,
		}
	default:
		return &game.ProfessionDetails{
			ActionPoints: 60,
			Preparation:  200,
		}
	}
}
