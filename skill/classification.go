package skill

type Classification int

//go:generate stringer -type=Classification
const (
	Skill Classification = iota
	Attack
	Spell
	Attunement
)

func ClassificationFromString(s string) *Classification {
	for i := 0; i <= int(Attunement); i++ {
		c := Classification(i)

		if c.String() == s {
			return &c
		}
	}
	return nil
}
