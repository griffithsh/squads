package embark

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/item"
	"github.com/griffithsh/squads/mathx"
	"github.com/griffithsh/squads/skill"
)

type generator struct {
	r       *rand.Rand
	archive Archive
}

func newGenerator(archive Archive) generator {
	seed := time.Now().UnixNano()
	return generator{
		r:       rand.New(rand.NewSource(seed)),
		archive: archive,
	}
}

func (g *generator) generateChar() *game.Character {
	sex := g.generateSex()

	hairs := g.archive.HairVariations()
	skins := g.archive.SkinVariations()

	hair := hairs[g.r.Intn(len(hairs))]
	skin := skins[g.r.Intn(len(skins))]

	return &game.Character{
		Level:         1,
		Disambiguator: g.r.Float64(),

		Name:                 g.generateName(sex),
		Sex:                  sex,
		Profession:           "Villager",
		InherantPreparation:  -50 + g.r.Intn(101),
		InherantActionPoints: int(g.r.NormFloat64()*1.4 + 8),
		Hair:                 hair,
		Skin:                 skin,
		PortraitBG:           g.r.Intn(mathx.MinI(len(game.PortraitBGBig), len(game.PortraitBGSmall))),
		PortraitFrame:        g.r.Intn(mathx.MinI(len(game.PortraitFrameBig), len(game.PortraitFrameSmall))),

		CurrentHealth:        17,
		BaseHealth:           25,
		StrengthPerLevel:     1.25 + g.r.Float64()*2.00,
		AgilityPerLevel:      1.25 + g.r.Float64()*2.00,
		IntelligencePerLevel: 0.75 + g.r.Float64()*1.25,
		VitalityPerLevel:     1.50 + g.r.Float64()*1.50,
		Masteries:            g.generateMasteries(),
	}
}

func (g *generator) generateWeapon() *item.Instance {
	// unarmed or sword or bow

	switch g.r.Intn(3) {
	case 1:
		return &item.Instance{
			Class:           item.SwordClass,
			Name:            "Skirmish Sword of Quickness",
			BaseChanceToHit: 0.99,
			Modifiers: map[item.Modifier]float64{
				item.BaseMinDamageModifier: 11,
				item.BaseMaxDamageModifier: 22,
				item.PreparationModifier:   599,
				item.ActionPointModifier:   21,
			},
			Skills: []skill.ID{
				"sword-attack",
				"sword-slash",
				"sword-stab",
				"sword-hew",
				"sword-rogue-slash",
			},
		}
	case 2:
		return &item.Instance{
			Class:           item.BowClass,
			Name:            "Ferocious Longbow",
			BaseChanceToHit: 0.6,
			Modifiers: map[item.Modifier]float64{
				item.BaseMinDamageModifier: 6,
				item.BaseMaxDamageModifier: 17,
				item.PreparationModifier:   461,
				item.ActionPointModifier:   25,
			},
			Skills: []skill.ID{
				"bow-attack",
				"bow-straight",
				"bow-quick",
				"bow-focus",
				"bow-ballistic",
				"bow-skirmish",
			},
		}
	}
	return nil
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
	names := g.archive.Names()
	keys := make([]string, len(names))

	for key := range names {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	g.r.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	var satisfactory func(key string) bool

	switch sex {
	case game.Male:
		satisfactory = func(key string) bool {
			tags := names[key]
			for _, tag := range tags {
				if tag == "M" {
					return true
				}
			}
			return false
		}
	case game.Female:
		satisfactory = func(key string) bool {
			tags := names[key]
			for _, tag := range tags {
				if tag == "F" {
					return true
				}
			}
			return false
		}
	}
	for _, name := range keys {
		if satisfactory(name) {
			return name
		}
	}
	return fmt.Sprintf("Unhandled Sex %s", sex)
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
