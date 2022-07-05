package skill

type Classification int

//go:generate go run github.com/dmarkham/enumer -type=Classification -json
const (
	Skill Classification = iota
	Attack
	Spell
	Attunement
)
