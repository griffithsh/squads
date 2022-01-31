package combat

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/griffithsh/squads/skill"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui"
)

/*
There are two (pseudo-public) methods to control all UI groups.

ShowX
HideX

Another method is only called internally - by the Update method.

RepaintX

- Show must add all entities, their tags, and parent relationships and the
invalidated tag. Because it's adding entities, it needs to ensure that it
does not add duplicates. It's the only method that should call to
mgr.NewEntity().

- Hide must call DestroyEntity on all the Entities that compose the group.
Anything that is added by Show.

- Repaint must add/replace the Position, Sprite, Font etc Components,
and then remove the invalidated tag.

There is no hierarchy between groups, each set of these functions is an
island disconnected from the others, and is only responsible for its own
Entities.
*/

var (
	combatHUDTag   = "COMBAT_HUD"
	timePassingTag = combatHUDTag + ".TIME_PASSING"

	invalidatedTag = combatHUDTag + ".INVALIDATED"
)

// HUD is a heads up display for the combat game state.
type HUD struct {
	mgr              *ecs.World
	bus              *event.Bus
	archive          SkillArchive
	scale            float64
	layer            int
	centerX, centerY float64 // center of the game's window, or half the width and height.
	lastCombatState  State

	// Whose turn is it?
	turnToken ecs.Entity

	dormant bool

	// Beta stuff:

	uiEntity             ecs.Entity
	turnQueueUIComponent *ui.UI
	fullUIComponent      *ui.UI
}

// NewHUD constructs a HUD.
func NewHUD(mgr *ecs.World, bus *event.Bus, screenX int, screenY int, archive SkillArchive) *HUD {
	makeUI := func(file string) *ui.UI {
		f, err := os.Open(file)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		return ui.NewUI(f)

	}
	hud := HUD{
		mgr:             mgr,
		bus:             bus,
		archive:         archive,
		scale:           2,
		layer:           1000,
		centerX:         float64(screenX) / 2,
		centerY:         float64(screenY) / 2,
		lastCombatState: Uninitialised,
		dormant:         true,

		uiEntity:             mgr.NewEntity(),
		fullUIComponent:      makeUI("game/combat/ui.xml"),
		turnQueueUIComponent: makeUI("game/combat/turnQueue.xml"),
	}

	bus.Subscribe(game.WindowSizeChanged{}.Type(), hud.handleWindowSizeChanged)
	bus.Subscribe(StateTransition{}.Type(), hud.handleCombatStateTransition)
	bus.Subscribe(ParticipantTurnChanged{}.Type(), hud.handleParticipantTurnChanged)
	bus.Subscribe(DamageAccepted{}.Type(), hud.handleDamageAccepted)

	return &hud
}

// Enable is the opposite of Disable. It shows anything that should be shown
// based on the current state.
func (hud *HUD) Enable() {
	if hud.lastCombatState == PreparingState {
		hud.showTimePassingIcon()
	}

	hud.dormant = false
}

// Disable the hud and everything in it, and ignore events that would show parts of
// the hud until Enable() is called.
func (hud *HUD) Disable() {
	hud.hideTimePassingIcon()

	hud.dormant = true
}

func (hud *HUD) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*game.WindowSizeChanged)
	hud.centerX, hud.centerY = float64(wsc.NewW)/2, float64(wsc.NewH)/2

	if hud.dormant {
		return
	}

	// If the window size has changed, then the center's absolute pixel
	// coordinates have also changed, so we should invalidate the time passing
	// icon.
	if hud.mgr.AnyTagged(timePassingTag) != 0 {
		hud.showTimePassingIcon()
	}
}

func (hud *HUD) handleCombatStateTransition(ev event.Typer) {
	if hud.dormant {
		return
	}

	cst := ev.(*StateTransition)
	hud.lastCombatState = cst.New.Value()

	if cst.New == PreparingState {
		hud.showTimePassingIcon()
	} else {
		hud.hideTimePassingIcon()
	}
}

func (hud *HUD) handleParticipantTurnChanged(t event.Typer) {
	ev := t.(*ParticipantTurnChanged)
	hud.turnToken = ev.Entity
}

func (hud *HUD) handleDamageAccepted(t event.Typer) {
	ev := t.(*DamageAccepted)

	targetPosition := hud.mgr.Component(ev.Target, "Position").(*game.Position)

	e := hud.mgr.NewEntity()
	text := strconv.Itoa(ev.Amount)
	hud.mgr.AddComponent(e, &game.Font{
		Text: text,
	})

	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: targetPosition.Center.X - float64(len(text)*5)/2,
			Y: targetPosition.Center.Y - 32,
		},
		Layer: hud.layer,
	})

	hud.mgr.AddComponent(e, &game.FloatAwayAnimation{
		Rate: 6.5,
	})

	hud.mgr.AddComponent(e, &ecs.Expiry{
		Remaining: time.Millisecond * 1500,
	})
}

func (hud *HUD) skillsForParticipant(p *Participant) [7]UISkillInfoRow {
	// convert a *skill.Description to a UISkillInfo
	convert := func(sd *skill.Description) UISkillInfo {
		spr := sd.Icon.Frames[sd.Icon.Index()]
		return UISkillInfo{
			Id:      string(sd.ID),
			Texture: spr.Texture,
			IconX:   spr.X,
			IconY:   spr.Y,
			Handle: func(string) {
				hud.bus.Publish(&SkillRequested{
					Code: sd.ID,
				})
			},
		}
	}

	var result [7]UISkillInfoRow
	for i := 0; i < 7; i++ {
		var row [2]UISkillInfo
		for i := 0; i < 2; i++ {
			info := &row[i]
			info.Texture = "hud.png"
			info.IconX = 184
			info.IconY = 0
			info.Handle = func(string) {}
		}
		result[i].Skills = row
	}
	if hud.lastCombatState == AwaitingInputState {
		// Movement
		result[0].Skills[0] = convert(hud.archive.Skill(skill.BasicMovement))

		// Consumables
		result[0].Skills[1] = UISkillInfo{
			Texture: "hud.png",
			IconX:   232,
			IconY:   0,
			Handle: func(string) {
				// TODO: implement consumables
			},
		}

		// Flee
		result[6].Skills[0] = UISkillInfo{
			Texture: "hud.png",
			IconX:   184,
			IconY:   24,
			Handle: func(string) {
				hud.bus.Publish(&AttemptingEscape{Entity: hud.turnToken})
			},
		}

		// End turn
		result[6].Skills[1] = UISkillInfo{
			Texture: "hud.png",
			IconX:   208,
			IconY:   24,
			Handle: func(string) {
				hud.bus.Publish(&EndTurnRequested{})
			},
		}

		freeSlots := []*UISkillInfo{
			&result[1].Skills[0],
			&result[1].Skills[1],
			&result[2].Skills[0],
			&result[2].Skills[1],
			&result[3].Skills[0],
			&result[3].Skills[1],
			&result[4].Skills[0],
			&result[4].Skills[1],
			&result[5].Skills[0],
			&result[5].Skills[1],
		}
		for i, id := range p.Skills {
			sd := hud.archive.Skill(id)
			if i >= len(freeSlots) {
				break
			}
			info := convert(sd)
			slot := freeSlots[i]
			slot.Id = info.Id
			slot.Texture = info.Texture
			slot.IconX = info.IconX
			slot.IconY = info.IconY
			slot.Handle = info.Handle
		}

	} else {
		result[0].Skills[0] = UISkillInfo{
			Texture: "hud.png",
			IconX:   208,
			IconY:   0,
			Handle: func(string) {
				hud.bus.Publish(&CancelSkillRequested{})
			},
		}
	}
	return result
}

// Update the HUD. Synchronise the current game state to the Entities that compose it.
func (hud *HUD) Update(elapsed time.Duration) {
	e := hud.mgr.AnyTagged(timePassingTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintTimePassingIcon()
	}

	if hud.dormant {
		return
	}

	participantEntities := hud.mgr.Get([]string{"Participant"})
	participants := make([]*Participant, len(participantEntities))
	for i, e := range participantEntities {
		participants[i] = hud.mgr.Component(e, "Participant").(*Participant)
	}
	sort.Slice(participants, func(i, j int) bool {
		ip, jp := participants[i], participants[j]
		ipRem := ip.PreparationThreshold.Max - ip.PreparationThreshold.Cur
		jpRem := jp.PreparationThreshold.Max - jp.PreparationThreshold.Cur
		if ipRem != jpRem {
			return ipRem < jpRem
		}
		return ip.Disambiguator < jp.Disambiguator
	})

	turnQueue := make([]QueuedParticipant, 0)
	for _, participant := range participants {
		// Knocked Down and Escaped Characters cannot take turns.
		if participant.Status != Alive {
			continue
		}

		turnQueue = append(turnQueue, QueuedParticipant{
			Background:    participant.SmallPortraitBG.Texture,
			BackgroundY:   participant.SmallPortraitBG.Y,
			BackgroundX:   participant.SmallPortraitBG.X,
			Portrait:      participant.SmallIcon.Texture,
			PortraitX:     participant.SmallIcon.X,
			PortraitY:     participant.SmallIcon.Y,
			OverlayFrame:  participant.SmallPortraitFrame.Texture,
			OverlayFrameX: participant.SmallPortraitFrame.X,
			OverlayFrameY: participant.SmallPortraitFrame.Y,

			Prep:    participant.PreparationThreshold.Cur,
			PrepMax: participant.PreparationThreshold.Max,
		})
	}

	// Add the right ui component given current combat state
	hud.mgr.RemoveComponent(hud.uiEntity, &ui.UI{})
	switch hud.lastCombatState {
	case PreparingState:
		hud.turnQueueUIComponent.Data = struct{ TurnQueue []QueuedParticipant }{
			// NB the turn queue contains ALL participants, and is the only relevant
			// field for the turnQueue UI.
			TurnQueue: turnQueue,
		}
		hud.mgr.AddComponent(hud.uiEntity, hud.turnQueueUIComponent)
	case AwaitingInputState, SelectingTargetState:
		participant := hud.mgr.Component(hud.turnToken, "Participant").(*Participant)
		hud.fullUIComponent.Data = HUDData{
			Background:    participant.BigPortraitBG.Texture,
			BackgroundY:   participant.BigPortraitBG.Y,
			BackgroundX:   participant.BigPortraitBG.X,
			Portrait:      participant.BigIcon.Texture,
			PortraitX:     participant.BigIcon.X,
			PortraitY:     participant.BigIcon.Y,
			OverlayFrame:  participant.BigPortraitFrame.Texture,
			OverlayFrameX: participant.BigPortraitFrame.X,
			OverlayFrameY: participant.BigPortraitFrame.Y,

			Name:      participant.Name,
			Health:    participant.CurrentHealth,
			HealthMax: participant.maxHealth(),
			Energy:    0, // FIXME: implement energy
			EnergyMax: 0, // FIXME: implement energy
			Action:    participant.ActionPoints.Cur,
			ActionMax: participant.ActionPoints.Max,
			Prep:      participant.PreparationThreshold.Cur,
			PrepMax:   participant.PreparationThreshold.Max,

			// NB the turn queue contains all but the prepared participant.
			TurnQueue: turnQueue[:len(turnQueue)-1],

			Skills: hud.skillsForParticipant(participant),
		}
		hud.mgr.AddComponent(hud.uiEntity, hud.fullUIComponent)
	default:
		// Do neither
	}
}

func (hud *HUD) showTimePassingIcon() {
	e := hud.mgr.AnyTagged(timePassingTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	e = hud.mgr.NewEntity()
	hud.mgr.Tag(e, timePassingTag)
	hud.mgr.Tag(e, invalidatedTag)
}

func (hud *HUD) hideTimePassingIcon() {
	e := hud.mgr.AnyTagged(timePassingTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintTimePassingIcon() {
	e := hud.mgr.AnyTagged(timePassingTag)
	if e == 0 {
		return
	}
	hud.mgr.AddComponent(e, &game.Sprite{
		Texture: "hud.png",
		X:       16,
		Y:       0,
		W:       16,
		H:       24,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: hud.centerX,
			Y: hud.centerY,
		},
		Layer:    hud.layer,
		Absolute: true,
	})
	hud.mgr.RemoveTag(e, invalidatedTag)
}
