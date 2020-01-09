package embark

import (
	"math/rand"
	"time"

	"github.com/griffithsh/squads/game"
)

type generator struct {
	r *rand.Rand
}

func newGenerator() generator {
	seed := time.Now().UnixNano()
	return generator{
		r: rand.New(rand.NewSource(seed)),
	}
}

func (g *generator) generateChar() *game.Character {
	sex := g.generateSex()

	small, big := g.generateIcons(sex)
	return &game.Character{
		Name:                 g.generateName(sex),
		Sex:                  sex,
		Profession:           game.Villager,
		PreparationThreshold: 685 + g.r.Intn(31),
		ActionPoints:         100,
		SmallIcon:            small,
		BigIcon:              big,
	}
}

func (g *generator) generateSex() game.CharacterSex {
	switch g.r.Int() % 2 {
	case 0:
		return game.Male
	default:
		return game.Female
	}
}
func (g *generator) generateName(sex game.CharacterSex) string {
	switch sex {
	case game.Male:
		i := g.r.Intn(len(maleNames))
		return maleNames[i]
	case game.Female:
		i := g.r.Intn(len(femaleNames))
		return femaleNames[i]
	}
	return "Samithee"
}

var maleNames = []string{
	"Arnold",
	"Bolus",
	"Cristian",
	"Devuan",
	"Edward",
	"Frederick",
	"Gerald",
	"Humperdink",
	"Ignatius",
	"Jamieson",
	"Kotlin",
	"Lucas",
	"Mattieson",
	"Nelson",
	"Ormond",
	"Oswalt",
	"Perogue",
	"Quincy",
	"Rederick",
	"Samithee",
	"Staunton",
	"Timjamen",
	"Thames",
	"Unicerve",
	"Variose",
	"Xactabol",
	"Yalladin",
	"Zomparion",
}

var femaleNames = []string{
	"Alyssa",
	"Balustrade",
	"Callisto",
	"Divernon",
	"Eloa",
	"Fankrastha",
	"Gordania",
	"Helloise",
	"Ismalloray",
	"Jannifern",
	"Katherita",
	"Legothory",
	"Maillorne",
	"Nostory",
	"Ollivene",
	"Pursivonian",
	"Qui",
	"Rimcy",
	"Sallivoce",
	"Theodora",
	"Undine",
	"Violet",
	"Winchester",
	"Xin",
	"Yellow",
	"Zenta",
}

func (g *generator) generateIcons(sex game.CharacterSex) (small game.Sprite, big game.Sprite) {
	switch g.r.Int() % 2 {
	case 0:
		small = game.Sprite{
			Texture: "hud.png",
			X:       0,
			Y:       76,
			W:       26,
			H:       26,
		}
		big = game.Sprite{
			Texture: "hud.png",
			X:       0,
			Y:       24,
			W:       52,
			H:       52,
		}
	default:
		small = game.Sprite{
			Texture: "portraits.png",
			X:       178,
			Y:       230,
			W:       26,
			H:       26,
		}
		big = game.Sprite{
			Texture: "portraits.png",
			X:       204,
			Y:       204,
			W:       52,
			H:       52,
		}
	}

	return small, big
}
