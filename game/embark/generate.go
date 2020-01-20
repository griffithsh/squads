package embark

import (
	"fmt"
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
		Level:         1,
		Disambiguator: g.r.Float64(),

		Name:                 g.generateName(sex),
		Sex:                  sex,
		Profession:           game.Villager,
		PreparationThreshold: 685 + g.r.Intn(31),
		ActionPoints:         100,
		SmallIcon:            small,
		BigIcon:              big,

		StrengthPerLevel:     1.25 + g.r.Float64()*2.00,
		AgilityPerLevel:      1.25 + g.r.Float64()*2.00,
		IntelligencePerLevel: 0.75 + g.r.Float64()*1.25,
		VitalityPerLevel:     1.50 + g.r.Float64()*1.50,
	}
}

func (g *generator) generateSex() game.CharacterSex {
	// N.B. 33% female.
	switch g.r.Int() % 3 {
	case 0:
		return game.Female
	default:
		return game.Male
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
	return fmt.Sprintf("Unhandled Sex %s", sex)
}

var maleNames = []string{
	"Arneldo",
	"Atlus",
	"Axis",
	"Bentholemew",
	"Bolus",
	"Cristian",
	"Callo",
	"Devuan",
	"Edward",
	"Frederick",
	"Gerald",
	"Humperdink",
	"Ignatius",
	"Ignold",
	"Jamieson",
	"Jahnsenn",
	"Krastin",
	"Lucas",
	"Mattieson",
	"Nelson",
	"Nolan",
	"Ormond",
	"Oswalt",
	"Panseur",
	"Perogue",
	"Polter",
	"Punt",
	"Quincy",
	"Ramathese",
	"Rederick",
	"Roon",
	"Samithee",
	"Staunton",
	"Timjamen",
	"Thames",
	"Unicerve",
	"Variose",
	"Volturbulent",
	"Xactabol",
	"Yalladin",
	"Yossarian",
	"Zod",
	"Zomparion",
}

var femaleNames = []string{
	"Alyssa",
	"Balustrade",
	"Callisto",
	"Divernon",
	"Eloa",
	"Fankrastha",
	"Gao",
	"Gordania",
	"Hana",
	"Harmonia",
	"Helloise",
	"Ismalloray",
	"Jannifern",
	"Katherita",
	"Kalisto",
	"Kamio",
	"Ketlin",
	"Lanneth",
	"Legothory",
	"Maillorne",
	"Nostory",
	"Ollivene",
	"Pursivonian",
	"Qui",
	"Rimcy",
	"Sallivoce",
	"Sera",
	"Shanto",
	"Theodora",
	"Undine",
	"Victohia",
	"Violet",
	"Winchester",
	"Xin",
	"Yellow",
	"Zenta",
}

var maleIcons = []int{
	0, 1, 3, 4, 5, 6,
}

var femaleIcons = []int{
	0, 2, 3, 5, 7, 8,
}

func (g *generator) generateIcons(sex game.CharacterSex) (small game.Sprite, big game.Sprite) {
	i := 0
	switch sex {
	case game.Male:
		i = g.r.Intn(len(maleIcons))
		i = maleIcons[i]
	default:
		i = g.r.Intn(len(femaleIcons))
		i = femaleIcons[i]
	}
	return game.Sprite{
			Texture: "portraits-26.png",
			X:       i * 26,
			Y:       0,
			W:       26,
			H:       26,
		}, game.Sprite{
			Texture: "portraits-52.png",
			X:       i * 52,
			Y:       0,
			W:       52,
			H:       52,
		}
}
