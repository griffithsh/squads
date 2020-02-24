package skill

// Some skills should be staticly compiled into the executable, so that they can
// have special logic. They're not really part of game data.
const (
	BasicMovement  = "static-movement"
	UseConsumable  = "use-consumable"
	FleeFromCombat = "flee-from-combat"
	EndTurn        = "end-turn"
)
