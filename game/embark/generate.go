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

		Disambiguator: g.r.Float64(),
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
	return "Samithee"
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
	"Theodora",
	"Undine",
	"Victohia",
	"Violet",
	"Winchester",
	"Xin",
	"Yellow",
	"Zenta",
}

func (g *generator) generateIcons(sex game.CharacterSex) (small game.Sprite, big game.Sprite) {
	i := g.r.Int() % 9

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
