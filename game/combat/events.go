package combat

import (
	"fmt"
	"reflect"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/skill"
)

// StateTransition occurs when the combat's state changes
type StateTransition struct {
	Old, New StateContext
}

// Type of the Event.
func (StateTransition) Type() event.Type {
	return "combat.StateTransition"
}

// DifferentHexSelected occurs when the user has selected a different hex - i.e.
// via mousing over.
type DifferentHexSelected struct {
	K       *geom.Key
	Context StateContext
}

// Type of the Event.
func (DifferentHexSelected) Type() event.Type {
	return "combat.DifferentHexSelected"
}

// ParticipantTurnChanged occurs when the Character whose turn it is changes.
type ParticipantTurnChanged struct {
	Entity ecs.Entity
}

// Type of the Event.
func (v ParticipantTurnChanged) Type() event.Type {
	t := reflect.TypeOf(v)
	return event.Type(fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()))
}

// StatModified occurs when an Participant's current stats changed.
type StatModified struct {
	Entity ecs.Entity
	Stat   game.StatType
	Amount int
}

// Type of the Event.
func (StatModified) Type() event.Type {
	return "combat.StatModified"
}

// EndTurnRequested occurs when the player indicates that they are finished
// commanding the current Character.
type EndTurnRequested struct{}

// Type of the Event.
func (EndTurnRequested) Type() event.Type {
	return "combat.EndTurnRequested"
}

// CancelSkillRequested occurs when the player indicates they want to cancel
// targeting of the skill they selected.
type CancelSkillRequested struct{}

// Type of the Event.
func (CancelSkillRequested) Type() event.Type {
	return "combat.CancelSkillRequested"
}

// SkillRequested occurs when the player indicates they would like to use a
// skill.
type SkillRequested struct {
	Code skill.ID
}

// Type of the Event.
func (SkillRequested) Type() event.Type {
	return "combat.SkillRequested"
}

// ParticipantMoving occurs when a Character has begun their movement.
type ParticipantMoving struct {
	Entity               ecs.Entity
	NewSpeed, OldSpeed   float64
	OldFacing, NewFacing geom.DirectionType
}

// Type of the Event.
func (ParticipantMoving) Type() event.Type {
	return "combat.ParticipantMoving"
}

// ParticipantMovementConcluded occurs when a Character has finished their movement.
type ParticipantMovementConcluded struct {
	Entity ecs.Entity
}

// Type of the Event.
func (ParticipantMovementConcluded) Type() event.Type {
	return "combat.ParticipantMovementConcluded"
}

// AttemptingEscape occurs when a character is attempting to escape from combat.
type AttemptingEscape struct {
	Entity ecs.Entity
}

// Type of the Event.
func (AttemptingEscape) Type() event.Type {
	return "combat.AttemptingEscape"
}

// CharacterCelebrating occurs when a character has something to shout about.
type CharacterCelebrating struct {
	Entity ecs.Entity
}

// Type of the Event.
func (CharacterCelebrating) Type() event.Type {
	return "combat.CharacterCelebrating"
}

// UsingSkill occurs when a character has triggered a skill.
type UsingSkill struct {
	User     ecs.Entity
	Skill    skill.ID
	Selected *geom.Hex
}

// Type of the Event.
func (UsingSkill) Type() event.Type {
	return "combat.UsingSkill"
}

// SkillUseConcluded occurs when a character has finished using their skill.
type SkillUseConcluded struct {
	User     ecs.Entity
	Skill    skill.ID
	Selected *geom.Hex
}

// Type of the Event.
func (SkillUseConcluded) Type() event.Type {
	return "combat.SkillUseConcluded"
}

// DamageApplied is an Event that occurs when something has sent raw, un-reduced
// damage to be applied to a Partcipant.
type DamageApplied struct {
	Amount     int
	Target     ecs.Entity
	DamageType game.DamageType
	SkillType  skill.Classification
}

// Type of the Event.
func (DamageApplied) Type() event.Type {
	return "combat.DamageApplied"
}

// DamageAccepted is an event that occurs when a Participant has lost health
// points.
type DamageAccepted struct {
	Target     ecs.Entity
	Amount     int
	Reduced    int
	DamageType game.DamageType
}

// Type of the Event.
func (DamageAccepted) Type() event.Type {
	return "combat.DamageAccepted"
}

// DamageFailed occurs when the application of damage fails to result in
// accepting any of it.
type DamageFailed struct {
	Target ecs.Entity
	Reason string // Negated|Dodged|Miss even maybe?
}

// Type of the Event.
func (DamageFailed) Type() event.Type {
	return "combat.DamageFailed"
}

// ParticipantDied occurs when a combat Participant is KnockedDown.
type ParticipantDied struct {
	Entity ecs.Entity
}

// Type of the Event.
func (ParticipantDied) Type() event.Type {
	return "combat.ParticipantDied"
}

// ParticipantRevived occurs when a combat Participant is KnockedDown.
type ParticipantRevived struct {
	Entity ecs.Entity
}

// Type of the Event.
func (ParticipantRevived) Type() event.Type {
	return "combat.ParticipantRevived"
}

// ParticipantDefiled occurs when a combat Participant is Defiled and can no longer be resurrected.
type ParticipantDefiled struct {
	Entity ecs.Entity
}

// Type of the Event.
func (ParticipantDefiled) Type() event.Type {
	return "combat.ParticipantDefiled"
}

// CharacterEnteredCombat occurs when a combat new Participant (that was not
// present at the start of the combat) has entered combat.
type CharacterEnteredCombat struct {
	Level      int
	Profession string
	Team       *game.Team
	At         geom.Key
}

// Type of the Event.
func (CharacterEnteredCombat) Type() event.Type {
	return "combat.CharacterEnteredCombat"
}
