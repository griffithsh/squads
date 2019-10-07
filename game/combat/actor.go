package combat

import "github.com/griffithsh/squads/game"

// CurMax represents a value which has both a Current and Maximum value.
type CurMax struct {
	Cur int
	Max int
}

// Actor is a transient aggregation of the stats of a Character for the purposes
// of combat. Actors are created at the beginning of combat and are destroyed at
// the end of combat.
type Actor struct {
	Name      string
	Level     uint
	SmallIcon game.Sprite // (26x26)
	BigIcon   game.Sprite // (52x52)

	Size       game.CharacterSize
	Profession game.CharacterProfession
	Sex        game.CharacterSex

	PreparationThreshold CurMax
	ActionPoints         CurMax
	Health               CurMax

	Strength     int
	Dexterity    int
	Intelligence int
	Vitality     int
}

// Type of this Component.
func (*Actor) Type() string {
	return "Actor"
}
