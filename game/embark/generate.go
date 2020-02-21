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
		InherantPreparation:  -50 + g.r.Intn(101),
		InherantActionPoints: int(g.r.NormFloat64()*1.4 + 8),
		SmallIcon:            small,
		BigIcon:              big,

		StrengthPerLevel:     1.25 + g.r.Float64()*2.00,
		AgilityPerLevel:      1.25 + g.r.Float64()*2.00,
		IntelligencePerLevel: 0.75 + g.r.Float64()*1.25,
		VitalityPerLevel:     1.50 + g.r.Float64()*1.50,
		Masteries:            g.generateMasteries(),
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
	"Bacon",
	"Bentholemew",
	"Bolus",
	"Cristian",
	"Callo",
	"Devuan",
	"Donkey",
	"Dungaree",
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
	"Pantry",
	"Perogue",
	"Polter",
	"Punt",
	"Quincy",
	"Ramathese",
	"Rederick",
	"Roon",
	"Samithee",
	"Satchel",
	"Staunton",
	"Timjamen",
	"Thames",
	"Unicerve",
	"Variose",
	"Volturbulent",
	"Xactabol",
	"Xerxes",
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
	"Euphemia",
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
	"Satchel",
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

func (g *generator) generateMasteries() map[game.Mastery]int {
	result := map[game.Mastery]int{
		// Randomly distributed 0-2
		game.ShortRangeMeleeMastery: 0,
		game.LongRangeMeleeMastery:  0,
		game.RangedCombatMastery:    0,

		// 33% chance of 1
		game.CraftsmanshipMastery: 0,

		// Randomly distributed 0-2
		game.FireMastery:      0,
		game.WaterMastery:     0,
		game.EarthMastery:     0,
		game.AirMastery:       0,
		game.LightningMastery: 0,

		// 25% chance of 1
		game.DarkMastery:  0,
		game.LightMastery: 0,
	}

	for i := 0; i < g.r.Intn(3); i++ {
		result[game.Mastery(g.r.Intn(3))]++
	}

	switch g.r.Intn(3) {
	case 2:
		result[game.CraftsmanshipMastery]++
	}

	for i := 0; i < g.r.Intn(3); i++ {
		result[game.Mastery(g.r.Intn(5)+4)]++
	}

	switch g.r.Intn(4) {
	case 3:
		result[game.Mastery(g.r.Intn(2)+9)]++
	}

	return result
}
