package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/griffithsh/squads/data"
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/embark"
	"github.com/griffithsh/squads/ui"
)

func setupEmbarkFocusCharacter(mgr *ecs.World, archive *data.Archive) {
	f, err := os.Open("output/demo.ui.xml")
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	e = mgr.NewEntity()
	nui := ui.NewUI(f)
	data := randomCharData(archive)
	nui.Data = &data
	mgr.AddComponent(e, nui)
}

func randomName(archive *data.Archive) string {
	names := archive.Names()
	// use go's randomised map iteration to pick a random one.
	for k := range names {
		return k
	}
	return "Noman"
}

func randomSkin(archive *data.Archive) string {
	skins := archive.SkinVariations()
	return skins[rand.Intn(len(skins))]
}
func randomHair(archive *data.Archive) string {
	hairs := archive.HairVariations()
	return hairs[rand.Intn(len(hairs))]
}
func randomSex() game.CharacterSex {
	if rand.Int()%2 == 0 {
		return game.Male
	}
	return game.Female
}

func randomCharData(archive *data.Archive) embark.CharacterSheetData {
	bgBig := game.PortraitBGBig[3]
	appearance := archive.Appearance("Villager", randomSex(), randomHair(archive), randomSkin(archive))
	portBig := appearance.BigIcon()
	overlayBig := game.PortraitFrameBig[0]

	masteries := []string{
		"Fishing: 2\n",
		"Hunting: 1\n",
		"Foightin: 5\n",
	}
	return embark.CharacterSheetData{
		Name:          randomName(archive),
		Profession:    "Villager",
		Lvl:           1,
		Sex:           "Male",
		Background:    bgBig.Texture,
		BackgroundX:   bgBig.X,
		BackgroundY:   bgBig.Y,
		Portrait:      portBig.Texture,
		PortraitX:     portBig.X,
		PortraitY:     portBig.Y,
		OverlayFrame:  overlayBig.Texture,
		OverlayFrameX: overlayBig.X,
		OverlayFrameY: overlayBig.Y,
		Prep:          91,
		AP:            202,
		Strlvl:        fmt.Sprintf("%.2f", 1.2),
		Agilvl:        fmt.Sprintf("%.2f", 0.876),
		Intlvl:        fmt.Sprintf("%.2f", 0.9012),
		Vitlvl:        fmt.Sprintf("%.2f", 2.005),
		Masteries:     masteries,

		HandleCancel: func(string) { fmt.Println("Cancel!") },
		HandleAction: func(string) { fmt.Println("Confirm") },
		ActionButton: "Do it!",
	}
}
