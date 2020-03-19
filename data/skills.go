package data

import (
	"fmt"
	"time"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/skill"
)

// SkillsByProfession gets the skills for a profession.
// FIXME: the game.CharacterProfession type must be removed for a
// runtime-configurable one. This should accept a string parameter instead.
func (a *Archive) SkillsByProfession(prof game.CharacterProfession) []*skill.Description {
	// FIXME: implementation
	return []*skill.Description{
		a.Skill("debug-basic-attack"),
		a.Skill("debug-lightning"),
		a.Skill("debug-revive"),
	}
}

// SkillsByWeaponClass provides the skills of a weapon class.
// FIXME: the skills should be provided by the instance of a class instead, so
// that a rapier can have different skills to a broadsword and to a scimitar.
func (a *Archive) SkillsByWeaponClass(weap game.ItemClass) []*skill.Description {
	// FIXME: implementation
	switch weap {
	case game.SwordClass:
		return []*skill.Description{
			a.Skill("debug-basic-attack"),
		}
	case game.UnarmedClass:
		fallthrough
	default:
		// Because other ItemClasses are armor, they provide no skills.
		return []*skill.Description{}
	}
}

// Skill retrieves a skill by its ID.
func (a *Archive) Skill(id skill.ID) *skill.Description {
	if val, ok := a.skills[id]; ok {
		return &val
	}
	panic(fmt.Sprintf("unconfigured skill %s", id))
}

// internalSkills are skills that are compiled into the binary instead of loaded
// in at run-time from a data file.
var internalSkills = []skill.Description{
	{
		ID:          skill.BasicMovement,
		Name:        "Move",
		Explanation: "Move to another tile",
		Tags:        []skill.Classification{skill.Skill},
		Icon: *game.Sprite{
			Texture: "hud.png",
			X:       232,
			Y:       24,
			W:       24,
			H:       24,
		}.AsAnimation(),

		Targeting:      skill.TargetAnywhere,
		TargetingBrush: skill.Pathfinding,
	},
	// consumables is a skill?
	// flee is a skill?
	// end turn is a skill?

	// Configure some test skills to develop and debug with.
	{
		ID:          "debug-basic-attack",
		Name:        "Attack",
		Explanation: "Attack an adjacent tile",
		Tags:        []skill.Classification{skill.Attack},
		Icon: *game.Sprite{
			Texture: "hud.png",
			X:       160,
			Y:       0,
			W:       24,
			H:       24,
		}.AsAnimation(),
		Targeting:      skill.TargetAdjacent,
		TargetingBrush: skill.SingleHex,
		Costs: map[skill.CostType]int{
			skill.CostsActionPoints: 20,
		},
		Effects: []skill.Effect{
			{
				When: skill.NewTiming(time.Millisecond * 500),
				What: skill.DamageEffect{
					Min: []skill.Operation{
						{Operator: skill.AddOp, Variable: "$DMG-MIN"},
					},
					Max: []skill.Operation{
						{Operator: skill.AddOp, Variable: "$DMG-MAX"},
					},
					Classification: skill.Attack,
				},
			},
		},
	},
	{
		ID:          "debug-lightning",
		Name:        "Mage Lightning",
		Explanation: "A lightning bolt strikes the target dealing 1-10 damage",
		Tags:        []skill.Classification{skill.Spell},
		Icon: *game.Sprite{
			Texture: "hud.png",
			X:       160,
			Y:       24,
			W:       24,
			H:       24,
		}.AsAnimation(),
		Targeting:      skill.TargetAnywhere,
		TargetingBrush: skill.SingleHex,
		Costs: map[skill.CostType]int{
			skill.CostsActionPoints: 45,
		},
		Effects: []skill.Effect{
			{
				When: skill.NewTimingFromPoint(skill.AttackApexTimingPoint),
				What: skill.DamageEffect{
					Min: []skill.Operation{
						{Operator: skill.AddOp, Variable: "1"},
						{Operator: skill.MultOp, Variable: "$LIGHTNING"},
						{Operator: skill.AddOp, Variable: "1"},
					},
					Max: []skill.Operation{
						{Operator: skill.AddOp, Variable: "7"},
						{Operator: skill.MultOp, Variable: "$LIGHTNING"},
						{Operator: skill.AddOp, Variable: "10"},
					},
					Classification: skill.Spell,
				},
			},
		},
	},
	{
		ID:          "debug-revive",
		Name:        "Pheonix form",
		Explanation: "A pheonix feather lands on the target, reviving it.",
		Tags:        []skill.Classification{skill.Spell},
		Icon: *game.Sprite{
			Texture: "hud.png",
			X:       160,
			Y:       48,
			W:       24,
			H:       24,
		}.AsAnimation(),
		Targeting:      skill.TargetAnywhere,
		TargetingBrush: skill.SingleHex,
		Costs: map[skill.CostType]int{
			skill.CostsActionPoints: 45,
		},
		Effects: []skill.Effect{
			{
				When: skill.NewTimingFromPoint(skill.AttackApexTimingPoint),
				What: skill.ReviveEffect{},
			},
			{
				When: skill.NewTimingFromPoint(skill.AttackApexTimingPoint),
				What: skill.HealEffect{
					Amount:       0.15,
					IsPercentage: true,
				},
			},
		},
	},
}
