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
	combatHUDTag                       = "COMBAT_HUD"
	timePassingTag                     = combatHUDTag + ".TIME_PASSING"
	currentParticipantTag              = combatHUDTag + ".CURRENT_PARTICIPANT"
	currentParticipantHovererTag       = currentParticipantTag + ".HOVERER"
	currentParticipantPortraitTag      = currentParticipantTag + ".PORTRAIT"
	currentParticipantPortraitBGTag    = currentParticipantPortraitTag + ".BG"
	currentParticipantPortraitFrameTag = currentParticipantPortraitTag + ".FRAME"
	currentParticipantNameTag          = currentParticipantTag + ".NAME"
	currentParticipantStatsTag         = currentParticipantTag + ".STATS"
	skillsTag                          = currentParticipantTag + ".SKILLS"
	turnQueueTag                       = combatHUDTag + ".TURN_QUEUE"

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
		layer:           100,
		centerX:         float64(screenX) / 2,
		centerY:         float64(screenY) / 2,
		lastCombatState: Uninitialised,
		dormant:         true,

		uiEntity:             mgr.NewEntity(),
		fullUIComponent:      makeUI("game/combat/ui.xml"),
		turnQueueUIComponent: makeUI("game/combat/turnQueue.xml"),
	}

	bus.Subscribe(game.WindowSizeChanged{}.Type(), hud.handleWindowSizeChanged)
	bus.Subscribe(game.CombatBegan{}.Type(), hud.handleCombatBegan)
	bus.Subscribe(StatModified{}.Type(), hud.handleCombatStatModified)
	bus.Subscribe(StateTransition{}.Type(), hud.handleCombatStateTransition)
	bus.Subscribe(ParticipantTurnChanged{}.Type(), hud.handleParticipantTurnChanged)
	bus.Subscribe(DamageAccepted{}.Type(), hud.handleDamageAccepted)

	return &hud
}

// Enable is the opposite of Disable. It shows anything that should be shown
// based on the current state.
func (hud *HUD) Enable() {
	switch hud.lastCombatState {
	case AwaitingInputState, SelectingTargetState:
		hud.showSkills()
		hud.showCurrentParticipant()
		hud.showCurrentParticipantStats()

	case PreparingState:
		hud.showTimePassingIcon()
	}

	hud.showTurnQueue()

	hud.dormant = false
}

// Disable the hud and everything in it, and ignore events that would show parts of
// the hud until Enable() is called.
func (hud *HUD) Disable() {
	hud.hideSkills()
	hud.hideCurrentParticipant()
	hud.hideCurrentParticipantStats()
	hud.hideTimePassingIcon()
	hud.hideTurnQueue()

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

func (hud *HUD) handleCombatBegan(event.Typer) {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, combatHUDTag)

	if hud.dormant {
		return
	}

	hud.showTurnQueue()
}

func (hud *HUD) handleCombatStatModified(ev event.Typer) {
	if hud.dormant {
		return
	}

	csm := ev.(*StatModified)

	switch csm.Stat {
	case game.PrepStat:
		hud.showTurnQueue()
	}
}

func (hud *HUD) handleCombatStateTransition(ev event.Typer) {
	if hud.dormant {
		return
	}

	cst := ev.(*StateTransition)
	hud.lastCombatState = cst.New.Value()
	n, o := cst.New.Value(), cst.Old.Value()

	if n == AwaitingInputState || n == SelectingTargetState {
		hud.showSkills()
		if o != AwaitingInputState && o != SelectingTargetState {
			hud.showCurrentParticipant()
			hud.showCurrentParticipantStats()
		}
	} else {
		hud.hideSkills()
		hud.hideCurrentParticipant()
		hud.hideCurrentParticipantStats()
	}

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
		Layer: 100,
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
				// TODO
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

		// Skill slots 1, 2, 8, and 9 are reserved for skills provided by the
		// Character's equipped weapon.
		// weaponSlots := []int{1, 2, 8, 9}
		weaponSlots := []*UISkillInfo{
			&result[1].Skills[0],
			&result[1].Skills[1],
			&result[2].Skills[0],
			&result[2].Skills[1],
		}
		for i, sd := range hud.archive.SkillsByWeaponClass(p.EquippedWeaponClass) {
			if i >= len(weaponSlots) {
				break
			}
			info := convert(sd)
			slot := weaponSlots[i]
			slot.Id = info.Id
			slot.Texture = info.Texture
			slot.IconX = info.IconX
			slot.IconY = info.IconY
			slot.Handle = info.Handle
		}
		// Skill slots 3, 4, 5, 10, 11, and 12 are reserved for skills provided
		// by the Character's profession.
		// profSlots := []int{3, 4, 5, 10, 11, 12}
		profSlots := []*UISkillInfo{
			&result[3].Skills[0],
			&result[3].Skills[1],
			&result[4].Skills[0],
			&result[4].Skills[1],
			&result[5].Skills[0],
			&result[5].Skills[1],
		}
		for i, sd := range hud.archive.SkillsByProfession(p.Profession) {
			if i >= len(profSlots) {
				break
			}
			info := convert(sd)
			slot := profSlots[i]
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
	var e ecs.Entity

	e = hud.mgr.AnyTagged(timePassingTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintTimePassingIcon()
	}
	e = hud.mgr.AnyTagged(currentParticipantTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintCurrentParticipant()
	}
	e = hud.mgr.AnyTagged(currentParticipantStatsTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintCurrentParticipantStats()
	}
	e = hud.mgr.AnyTagged(skillsTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintSkills()
	}
	e = hud.mgr.AnyTagged(turnQueueTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintTurnQueue()
	}

	// TODO (re?)construct data for the UI
	// ...

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

func (hud *HUD) showCurrentParticipant() {
	return
	e := hud.mgr.AnyTagged(currentParticipantTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	parent := hud.mgr.NewEntity()
	hud.mgr.Tag(parent, currentParticipantTag)
	hud.mgr.Tag(parent, invalidatedTag)

	e = hud.mgr.NewEntity()
	hud.mgr.Dependency(parent, e)
	hud.mgr.Tag(e, currentParticipantHovererTag)

	e = hud.mgr.NewEntity()
	hud.mgr.Dependency(parent, e)
	hud.mgr.Tag(e, currentParticipantPortraitBGTag)
	e = hud.mgr.NewEntity()
	hud.mgr.Dependency(parent, e)
	hud.mgr.Tag(e, currentParticipantPortraitTag)
	e = hud.mgr.NewEntity()
	hud.mgr.Dependency(parent, e)
	hud.mgr.Tag(e, currentParticipantPortraitFrameTag)

	e = hud.mgr.NewEntity()
	hud.mgr.Dependency(parent, e)
	hud.mgr.Tag(e, currentParticipantNameTag)
}

func (hud *HUD) hideCurrentParticipant() {
	e := hud.mgr.AnyTagged(currentParticipantTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintCurrentParticipant() {
	parent := hud.mgr.AnyTagged(currentParticipantTag)
	if parent == 0 {
		return
	}

	// Repaint the hovering arrow that points to the current Character.
	ex := hud.turnToken
	if ex == 0 || hud.mgr.Component(ex, "Position") == nil {
		return
	}

	pos := hud.mgr.Component(ex, "Position").(*game.Position)

	e := hud.mgr.AnyTagged(currentParticipantHovererTag)
	hud.mgr.AddComponent(e, &game.Sprite{
		Texture: "hud.png",
		X:       32,
		Y:       0,
		W:       18,
		H:       5,
		OffsetY: -32,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: pos.Center,
		Layer:  hud.layer,
	})
	hud.mgr.AddComponent(e, game.NewHoverAnimation())

	// Repaint the current Participant's portrait.
	center := game.Center{
		X: 30 * hud.scale,
		Y: (hud.centerY * 2) - 30*hud.scale,
	}
	scale := game.Scale{
		X: hud.scale,
		Y: hud.scale,
	}
	participant := hud.mgr.Component(hud.turnToken, "Participant").(*Participant)
	e = hud.mgr.AnyTagged(currentParticipantPortraitBGTag)
	hud.mgr.AddComponent(e, &participant.BigPortraitBG)
	hud.mgr.AddComponent(e, &scale)
	hud.mgr.AddComponent(e, &game.Position{
		Center:   center,
		Layer:    hud.layer - 1,
		Absolute: true,
	})
	e = hud.mgr.AnyTagged(currentParticipantPortraitTag)
	hud.mgr.AddComponent(e, &participant.BigIcon)
	hud.mgr.AddComponent(e, &scale)
	hud.mgr.AddComponent(e, &game.Position{
		Center:   center,
		Layer:    hud.layer,
		Absolute: true,
	})
	e = hud.mgr.AnyTagged(currentParticipantPortraitFrameTag)
	hud.mgr.AddComponent(e, &participant.BigPortraitFrame)
	hud.mgr.AddComponent(e, &scale)
	hud.mgr.AddComponent(e, &game.Position{
		Center:   center,
		Layer:    hud.layer + 1,
		Absolute: true,
	})

	// Repaint the current Participant's name.
	e = hud.mgr.AnyTagged(currentParticipantNameTag)
	hud.mgr.AddComponent(e, &game.Font{
		Text: participant.Name,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 58 * hud.scale,
			Y: (hud.centerY * 2) - 54*hud.scale,
		},
		Layer:    hud.layer + 1,
		Absolute: true,
	})

	hud.mgr.RemoveTag(parent, invalidatedTag)
}

func (hud *HUD) showCurrentParticipantStats() {
	return
	e := hud.mgr.AnyTagged(currentParticipantStatsTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	parent := hud.mgr.NewEntity()
	hud.mgr.Tag(parent, currentParticipantStatsTag)
	hud.mgr.Tag(parent, invalidatedTag)

	// Create enough child entities for every stat and its label.
	children := make([]ecs.Entity, 8)
	for i := 0; i < 8; i++ {
		e := hud.mgr.NewEntity()
		hud.mgr.Dependency(parent, e)
		children[i] = e
	}
	hud.mgr.AddComponent(parent, &ecs.Children{
		Value: children,
	})
}

func (hud *HUD) hideCurrentParticipantStats() {
	e := hud.mgr.AnyTagged(currentParticipantStatsTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintCurrentParticipantStats() {
	parent := hud.mgr.AnyTagged(currentParticipantStatsTag)
	children := hud.mgr.Component(parent, "Children").(*ecs.Children)

	e := hud.turnToken
	if e == 0 {
		return
	}
	participant := hud.mgr.Component(e, "Participant").(*Participant)
	labels := []string{
		"Health:",
		fmt.Sprintf("%d/%d", participant.CurrentHealth, participant.maxHealth()),
		"Energy:",
		"?/N", // TODO
		"Action:",
		fmt.Sprintf("%d/%d", participant.ActionPoints.Cur, participant.ActionPoints.Max),
		"Prep:",
		fmt.Sprintf("%d/%d", participant.PreparationThreshold.Cur, participant.PreparationThreshold.Max),
	}
	for i, child := range children.Value {
		// Magical x,y coordinates for stats are:
		// x = 60 or 102
		// y = -43, -36, -29, or -22
		x, y := 60.0+float64(i%2)*42.0, -43.0+float64(i/2*7)
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: x * hud.scale,
				Y: (hud.centerY * 2) + y*hud.scale,
			},
			Layer:    hud.layer,
			Absolute: true,
		})

		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale,
			Y: hud.scale,
		})

		hud.mgr.AddComponent(child, &game.Font{
			Text: labels[i],
			Size: "small",
		})
	}
	hud.mgr.RemoveTag(parent, invalidatedTag)
}

func (hud *HUD) showSkills() {
	return
	// Create a parent entity tagged with invalidatedTag and skillsTag
	e := hud.mgr.AnyTagged(skillsTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	e = hud.mgr.NewEntity()
	hud.mgr.Tag(e, skillsTag)

	// Create 7x2 child Entities with parent Components
	children := make([]ecs.Entity, 14)
	for i := 0; i < 14; i++ {
		children[i] = hud.mgr.NewEntity()
		hud.mgr.Dependency(e, children[i])
	}

	// Append all child entities to parent
	hud.mgr.AddComponent(e, &ecs.Children{
		Value: children,
	})
	hud.mgr.Tag(e, invalidatedTag)
}

func (hud *HUD) hideSkills() {
	e := hud.mgr.AnyTagged(skillsTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintSkills() {
	parent := hud.mgr.AnyTagged(skillsTag)
	children := hud.mgr.Component(parent, "Children").(*ecs.Children)

	type v struct {
		icon        game.FrameAnimation
		interactive *ui.Interactive
	}
	convert := func(sd *skill.Description) v {
		var button *ui.Interactive
		icon := game.FrameAnimation{
			Timings: []time.Duration{1000000},
		}
		switch sd.ID {
		default:
			icon = sd.Icon
			button = &ui.Interactive{
				W: 24, H: 24,
				Trigger: func(x, y float64) {
					hud.bus.Publish(&SkillRequested{
						Code: sd.ID,
					})
				},
			}
		}
		return v{
			icon:        icon,
			interactive: button,
		}
	}

	skills := map[int]v{
		// Cancel button
		0: {
			icon: *game.Sprite{
				Texture: "hud.png",
				X:       208,
				Y:       0,
				W:       24,
				H:       24,
			}.AsAnimation(),
			interactive: &ui.Interactive{
				W: 24, H: 24,
				Trigger: func(x, y float64) {
					hud.bus.Publish(&CancelSkillRequested{})
				},
			},
		},
	}

	if hud.lastCombatState == AwaitingInputState {
		skills = map[int]v{
			// Zero is always Move.
			0: convert(hud.archive.Skill(skill.BasicMovement)),

			// Consumables
			7: {
				icon: *game.Sprite{
					Texture: "hud.png",
					X:       232,
					Y:       0,
					W:       24,
					H:       24,
				}.AsAnimation(),
			},

			// TODO: weapon skills should be provided by equipped weapon
			// 1 - weapon 1
			// 2 - weapon 2
			// 8 - weapon 3
			// 9 - weapon 4

			// TODO:
			// 3,4,5,10,11,12 are slots 1-6 for profession-provided skills

			// Flee
			6: {
				icon: *game.Sprite{
					Texture: "hud.png",
					X:       184,
					Y:       24,
					W:       24,
					H:       24,
				}.AsAnimation(),
				interactive: &ui.Interactive{
					W: 24, H: 24,
					Trigger: func(x, y float64) {
						hud.bus.Publish(&AttemptingEscape{Entity: hud.turnToken})
					},
				},
			},

			// End turn
			13: {
				icon: *game.Sprite{
					Texture: "hud.png",
					X:       208,
					Y:       24,
					W:       24,
					H:       24,
				}.AsAnimation(),
				interactive: &ui.Interactive{
					W: 24, H: 24,
					Trigger: func(x, y float64) {
						hud.bus.Publish(&EndTurnRequested{})
					},
				},
			},
		}

		participant := hud.mgr.Component(hud.turnToken, "Participant").(*Participant)

		// Skill slots 1, 2, 8, and 9 are reserved for skills provided by the
		// Character's equipped weapon.
		weaponSlots := []int{1, 2, 8, 9}
		for i, sd := range hud.archive.SkillsByWeaponClass(participant.EquippedWeaponClass) {
			if i >= len(weaponSlots) {
				break
			}
			key := weaponSlots[i]
			skills[key] = convert(sd)
		}
		// Skill slots 3, 4, 5, 10, 11, and 12 are reserved for skills provided
		// by the Character's profession.
		profSlots := []int{3, 4, 5, 10, 11, 12}
		for i, sd := range hud.archive.SkillsByProfession(participant.Profession) {
			if i >= len(profSlots) {
				break
			}
			key := profSlots[i]
			skills[key] = convert(sd)
		}
	}

	for i, child := range children.Value {
		x, y := i%7, i/7

		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale,
			Y: hud.scale,
		})
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: (150 + float64(26*x)) * hud.scale,
				Y: hud.centerY*2 + (-42+float64(26*y))*hud.scale,
			},
			Layer:    hud.layer,
			Absolute: true,
		})

		s, ok := skills[i]
		if !ok {
			// Add empty skill slot sprite component
			hud.mgr.AddComponent(child, &game.Sprite{
				Texture: "hud.png",
				X:       184,
				Y:       0,
				W:       24,
				H:       24,
			})
			continue
		}
		hud.mgr.AddComponent(child, &s.icon)
		if s.interactive != nil {
			hud.mgr.AddComponent(child, s.interactive)
		}
	}
	hud.mgr.RemoveTag(parent, invalidatedTag)
}

const turnQueueSlots int = 8
const entitiesPerTurnQueueSlot int = 5

func (hud *HUD) showTurnQueue() {
	return
	e := hud.mgr.AnyTagged(turnQueueTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	e = hud.mgr.NewEntity()
	hud.mgr.Tag(e, turnQueueTag)
	hud.mgr.Tag(e, invalidatedTag)

	children := make([]ecs.Entity, turnQueueSlots*entitiesPerTurnQueueSlot)
	for i := 0; i < turnQueueSlots*entitiesPerTurnQueueSlot; i++ {
		child := hud.mgr.NewEntity()
		hud.mgr.Dependency(e, child)
		children[i] = child
	}
	hud.mgr.AddComponent(e, &ecs.Children{
		Value: children,
	})
}

func (hud *HUD) hideTurnQueue() {
	e := hud.mgr.AnyTagged(turnQueueTag)
	if e == 0 {
		return
	}

	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintTurnQueue() {
	parents := hud.mgr.Tagged(turnQueueTag)
	if len(parents) > 1 {
		panic("incorrect use of " + turnQueueTag + ": multiple found")
	}
	parent := parents[0]
	if parent == 0 {
		return
	}
	type v struct {
		e             ecs.Entity
		remaining     int
		current, max  int
		icon          *game.Sprite
		bg            *game.Sprite
		frame         *game.Sprite
		disambiguator float64
	}

	var q []v
	for _, e := range hud.mgr.Get([]string{"Participant"}) {
		participant := hud.mgr.Component(e, "Participant").(*Participant)

		// Knocked Down and Escaped Characters cannot take turns.
		if participant.Status != Alive {
			continue
		}

		// An awkward way of not including the Character with the TurnToken.
		if participant.PreparationThreshold.Cur == 0 {
			continue
		}

		q = append(q, v{
			e:             e,
			remaining:     participant.PreparationThreshold.Max - participant.PreparationThreshold.Cur,
			current:       participant.PreparationThreshold.Cur,
			max:           participant.PreparationThreshold.Max,
			icon:          &participant.SmallIcon,
			bg:            &participant.SmallPortraitBG,
			frame:         &participant.SmallPortraitFrame,
			disambiguator: participant.Disambiguator,
		})
	}
	sort.Slice(q, func(i, j int) bool {
		if q[i].remaining != q[j].remaining {
			return q[i].remaining < q[j].remaining
		}
		return q[i].disambiguator < q[j].disambiguator
	})

	children := hud.mgr.Component(parent, "Children").(*ecs.Children)
	x, y := 10, 10+13
	stride := 42
	for i := 0; i < turnQueueSlots; i++ {
		// If we have no more Characters for this slot, then hide it.
		if i >= len(q) {
			hud.mgr.RemoveComponent(children.Value[i*entitiesPerTurnQueueSlot+0], &game.Sprite{})
			hud.mgr.RemoveComponent(children.Value[i*entitiesPerTurnQueueSlot+1], &game.Sprite{})
			hud.mgr.RemoveComponent(children.Value[i*entitiesPerTurnQueueSlot+2], &game.Sprite{})
			hud.mgr.RemoveComponent(children.Value[i*entitiesPerTurnQueueSlot+3], &game.Sprite{})
			hud.mgr.AddComponent(children.Value[i*entitiesPerTurnQueueSlot+4], &game.Font{})
			continue
		}

		v := q[i]

		// participant's portrait icon
		center := game.Center{
			X: float64(13+x+i*stride) * hud.scale,
			Y: float64(y) * hud.scale,
		}
		scale := game.Scale{
			X: hud.scale,
			Y: hud.scale,
		}
		child := children.Value[i*entitiesPerTurnQueueSlot]
		hud.mgr.AddComponent(child, v.bg)
		hud.mgr.AddComponent(child, &game.Position{
			Center:   center,
			Layer:    hud.layer - 1,
			Absolute: true,
		})
		hud.mgr.AddComponent(child, &scale)
		child = children.Value[i*entitiesPerTurnQueueSlot+1]
		hud.mgr.AddComponent(child, v.icon)
		hud.mgr.AddComponent(child, &game.Position{
			Center:   center,
			Layer:    hud.layer,
			Absolute: true,
		})
		hud.mgr.AddComponent(child, &scale)
		child = children.Value[i*entitiesPerTurnQueueSlot+2]
		hud.mgr.AddComponent(child, v.frame)
		hud.mgr.AddComponent(child, &game.Position{
			Center:   center,
			Layer:    hud.layer + 1,
			Absolute: true,
		})
		hud.mgr.AddComponent(child, &scale)

		// current preparation progressbar
		prepPerc := float64(v.current) / float64(v.max)
		child = children.Value[i*entitiesPerTurnQueueSlot+3]
		hud.mgr.Dependency(parent, child)
		hud.mgr.AddComponent(child, &game.Sprite{
			Texture: "tranquility-plus-39-palette.png",
			X:       1,
			Y:       2,
			W:       1,
			H:       1,
		})
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: (13*prepPerc + float64(x+i*stride)) * hud.scale,
				Y: float64(y+13+2) * hud.scale,
			},
			Layer:    hud.layer + 1,
			Absolute: true,
		})
		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale * 26 * prepPerc,
			Y: hud.scale * 4,
		})

		// current/max preparation text
		child = children.Value[i*entitiesPerTurnQueueSlot+4]
		hud.mgr.Dependency(parent, child)
		hud.mgr.AddComponent(child, &game.Font{
			Text: fmt.Sprintf("%d/%d", v.current, v.max),
			Size: "small",
		})
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: float64(x+i*stride) * hud.scale,
				Y: float64(y+13+2+4+2) * hud.scale,
			},
			Layer:    hud.layer,
			Absolute: true,
		})
		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale,
			Y: hud.scale,
		})
	}
	hud.mgr.RemoveTag(parent, invalidatedTag)
}
