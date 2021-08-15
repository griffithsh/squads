package embark

import (
	"fmt"
	"sort"

	"github.com/griffithsh/squads/game"
)

type CharacterSheetData struct {
	Name         string
	Profession   string
	Lvl          int
	Sex          string
	Portrait     string
	PortraitX    int
	PortraitY    int
	Prep         int
	AP           int
	Strlvl       string
	Agilvl       string
	Intlvl       string
	Vitlvl       string
	Masteries    []string
	HandleCancel func()
	HandleAction func()
	ActionButton string
}

func AsCharacterSheetData(char *game.Character, equip *game.Equipment, prof *game.ProfessionDetails, app *game.Appearance) CharacterSheetData {
	var masteries []string
	for mastery, lvl := range char.Masteries {
		if lvl == 0 {
			continue
		}
		masteries = append(masteries, fmt.Sprintf("%v: %d\n", mastery, lvl))
	}
	sort.Strings(masteries)
	// bg := game.PortraitBGBig[char.PortraitBG]
	// overlay := game.PortraitFrameBig[char.PortraitFrame]
	return CharacterSheetData{
		Name:       char.Name,
		Profession: char.Profession,
		Lvl:        char.Level,
		Sex:        char.Sex.String(),
		Portrait:   app.BigIcon().Texture,
		PortraitX:  app.BigIcon().X,
		PortraitY:  app.BigIcon().Y,
		Prep:       char.InherantPreparation + prof.Preparation + equip.WeaponPreparation(),
		AP:         char.InherantActionPoints + prof.ActionPoints + equip.WeaponActionPoints(),
		Strlvl:     fmt.Sprintf("%.2f", char.StrengthPerLevel),
		Agilvl:     fmt.Sprintf("%.2f", char.AgilityPerLevel),
		Intlvl:     fmt.Sprintf("%.2f", char.IntelligencePerLevel),
		Vitlvl:     fmt.Sprintf("%.2f", char.VitalityPerLevel),
		Masteries:  masteries,
	}
}
