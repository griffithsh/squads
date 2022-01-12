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
	"github.com/griffithsh/squads/skill"
	"github.com/griffithsh/squads/ui"
)

func setupCombatUIPreparing(mgr *ecs.World, archive *data.Archive) {
	f, err := os.Open("game/combat/turnQueue.xml")
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	e = mgr.NewEntity()
	nui := ui.NewUI(f)
	data := randHUDData(archive)
	nui.Data = &data
	mgr.AddComponent(e, nui)
}

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
	for i := 0; i < 3+rand.Intn(6); i++ {
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

		Skills: randSkills(archive),
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

func randSkills(archive *data.Archive) [7]combat.UISkillInfoRow {
	skillInfo := func(id string) combat.UISkillInfo {
		if id == "" {
			return combat.UISkillInfo{
				Texture: "hud.png",
				IconX:   184,
				IconY:   0,
				Id:      "",
				Handle: func(string) {
				},
			}
		}
		desc := archive.Skill(skill.ID(id))
		icon := desc.Icon.Frames[0]
		handler := func(string) {}
		if id != "" {
			handler = func(id string) {
				fmt.Printf("Skill %q!\n", id)
			}
		}
		return combat.UISkillInfo{
			Texture: icon.Texture,
			IconX:   icon.X,
			IconY:   icon.Y,
			Id:      id,
			Handle:  handler,
		}
	}
	skillInfoRow := func(id1, id2 string) combat.UISkillInfoRow {
		return combat.UISkillInfoRow{
			Skills: [2]combat.UISkillInfo{
				skillInfo(id1),
				skillInfo(id2),
			},
		}
	}
	return [7]combat.UISkillInfoRow{
		skillInfoRow("", ""),
		skillInfoRow("", ""),
		skillInfoRow("", ""),
		skillInfoRow("raise-skeleton", "debug-basic-attack"),
		skillInfoRow("", "debug-lightning"),
		skillInfoRow("", "debug-revive"),
		skillInfoRow("", ""),
	}
}
