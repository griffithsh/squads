package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/griffithsh/squads/data"
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/combat"
	"github.com/griffithsh/squads/mathx"
	"github.com/griffithsh/squads/ui"
)

func setupCombatUI(mgr *ecs.World, archive *data.Archive) {
	f, err := os.Open("game/combat/ui.xml")
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	e = mgr.NewEntity()
	nui := ui.NewUI(f)
	data := randHUDData(archive)
	nui.Data = &data
	mgr.AddComponent(e, nui)
}

func randHUDData(archive *data.Archive) combat.HUDData {
	bgBig := game.PortraitBGBig[0]
	appearance := archive.Appearance("Villager", game.Female, "black", "pale")
	portBig := appearance.BigIcon()
	overlayBig := game.PortraitFrameBig[0]

	turnQueue := []combat.QueuedParticipant{}
	for i := 0; i < 2+rand.Intn(3); i++ {
		turnQueue = append(turnQueue, randQueuedParticipant(archive))
	}
	return combat.HUDData{
		Background:    bgBig.Texture,
		BackgroundX:   bgBig.X,
		BackgroundY:   bgBig.Y,
		Portrait:      portBig.Texture,
		PortraitX:     portBig.X,
		PortraitY:     portBig.Y,
		OverlayFrame:  overlayBig.Texture,
		OverlayFrameX: overlayBig.X,
		OverlayFrameY: overlayBig.Y,

		Name: "Neuvo",

		Health: 80, HealthMax: 123,
		Energy: 803, EnergyMax: 1017,
		Action: 202, ActionMax: 202,
		Prep: 0, PrepMax: 91,

		TurnQueue: turnQueue,

		Skills: randSkills(),
	}
}

func randQueuedParticipant(archive *data.Archive) combat.QueuedParticipant {
	bg := game.PortraitBGSmall[rand.Intn(mathx.MinI(len(game.PortraitBGBig), len(game.PortraitBGSmall)))]
	sex := game.CharacterSex(rand.Int() % 2)
	hairs := archive.HairVariations()
	skins := archive.SkinVariations()

	hair := hairs[rand.Intn(len(hairs))]
	skin := skins[rand.Intn(len(skins))]
	appearance := archive.Appearance("Villager", sex, hair, skin)
	port := appearance.SmallIcon()
	frame := game.PortraitFrameSmall[rand.Intn(mathx.MinI(len(game.PortraitFrameBig), len(game.PortraitFrameSmall)))]
	return combat.QueuedParticipant{
		Background:    bg.Texture,
		BackgroundX:   bg.X,
		BackgroundY:   bg.Y,
		Portrait:      port.Texture,
		PortraitX:     port.X,
		PortraitY:     port.Y,
		OverlayFrame:  frame.Texture,
		OverlayFrameX: frame.X,
		OverlayFrameY: frame.Y,

		Prep: rand.Intn(180) + 589, PrepMax: rand.Intn(500) + 812,
	}
}

func randSkills() [7]combat.UISkillInfoRow {
	randSkillInfo := func() combat.UISkillInfo {
		return combat.UISkillInfo{
			Texture: "hud.png",
			IconX:   184,
			IconY:   0,
		}
	}
	randSkillInfoRow := func() combat.UISkillInfoRow {
		return combat.UISkillInfoRow{
			Skills: [2]combat.UISkillInfo{
				randSkillInfo(),
				randSkillInfo(),
			},
		}
	}
	return [7]combat.UISkillInfoRow{
		randSkillInfoRow(),
		randSkillInfoRow(),
		randSkillInfoRow(),
		randSkillInfoRow(),
		randSkillInfoRow(),
		randSkillInfoRow(),
		randSkillInfoRow(),
	}
}
