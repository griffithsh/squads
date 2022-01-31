package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/skill"
)

type skillEffect struct {
	When      int // milliseconds
	WhenPoint string
	What      interface{}
}

type skillDescription struct {
	ID          skill.ID
	Name        string
	Explanation string

	// Tags critically includes Attack or Spell, and allows the game to select
	// an appropriate animation to use when using the skill.
	Tags []string

	Icon game.Sprite

	Targeting      string
	TargetingBrush string

	// Effects of triggering this skill.
	Effects []skillEffect

	Costs map[string]int
}

// convert to a skill.Description
func (sd *skillDescription) convert() (skill.Description, error) {
	tags := []skill.Classification{}
	for _, raw := range sd.Tags {
		tag := skill.ClassificationFromString(raw)
		if tag == nil {
			fmt.Fprintf(os.Stderr, "unknown skill.Classification %s", raw)
			continue
		}
		tags = append(tags, *tag)
	}

	targetingRule := skill.TargetingRuleFromString(sd.Targeting)
	if targetingRule == nil {
		return skill.Description{}, fmt.Errorf("convert %q to targetingRule", sd.Targeting)
	}

	targetingBrush := skill.TargetingBrushFromString(sd.TargetingBrush)
	if targetingBrush == nil {
		return skill.Description{}, fmt.Errorf("convert %q to targetingBrush", sd.TargetingBrush)
	}

	effects := []skill.Effect{}
	for _, raw := range sd.Effects {
		when := skill.Timing{}
		point := skill.TimingPointFromString(raw.WhenPoint)
		if point != nil {
			when = skill.NewTimingFromPoint(*point)
		} else {
			when = skill.NewTiming(time.Duration(raw.When) * time.Millisecond)
		}
		effect := skill.Effect{
			When: when,
			What: raw.What,
		}
		effects = append(effects, effect)
	}
	costs := map[skill.CostType]int{}
	for k, v := range sd.Costs {
		costType := skill.CostTypeFromString(k)
		if costType == nil {
			return skill.Description{}, fmt.Errorf("convert to costType %q", k)
		}

		costs[*costType] = v
	}

	return skill.Description{
		ID:             sd.ID,
		Name:           sd.Name,
		Explanation:    sd.Explanation,
		Tags:           tags,
		Icon:           *sd.Icon.AsAnimation(),
		Targeting:      *targetingRule,
		TargetingBrush: *targetingBrush,
		Effects:        effects,
		Costs:          costs,
	}, nil
}

func (d *skillEffect) UnmarshalJSON(data []byte) error {
	v := struct {
		When      int
		WhenPoint string
		What      struct {
			Type string `json:"_type"`
		}
	}{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.When = v.When
	d.WhenPoint = v.WhenPoint

	switch v.What.Type {
	case "DamageEffect":
		v2 := struct {
			What skill.DamageEffect
		}{}
		json.Unmarshal(data, &v2)
		d.What = v2.What

	default:
		// TODO: other effect types!
		return fmt.Errorf("unknown _type %q (%v)", v.What.Type, v)
	}
	return nil
}

func parseSkills(r io.Reader) ([]skill.Description, error) {
	dec := json.NewDecoder(r)

	result := []skill.Description{}
	var s skillDescription
	for {
		err := dec.Decode(&s)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		sd, err := s.convert()
		if err != nil {
			return nil, fmt.Errorf("convert skillDescription: %v", err)
		}
		result = append(result, sd)
	}
	return result, nil
}

// SkillsByProfession gets the skills for a profession.
func (a *Archive) SkillsByProfession(prof string) []*skill.Description {
	// FIXME: implementation
	return []*skill.Description{
		// a.Skill("debug-basic-attack"),
		// a.Skill("debug-lightning"),
		// a.Skill("debug-revive"),
		// a.Skill("raise-skeleton"),
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
	panic(fmt.Sprintf("unconfigured skill %q, %d loaded skills", id, len(a.skills)))
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
						{Operator: skill.AddOp, Variable: "70"},
						{Operator: skill.MultOp, Variable: "$LIGHTNING"},
						{Operator: skill.AddOp, Variable: "100"},
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
	{
		ID:          "raise-skeleton",
		Name:        "Raise Skeleton",
		Explanation: "Raise the bones of the dead to fight alongside you.",
		Tags:        []skill.Classification{skill.Spell},
		Icon: *game.Sprite{
			Texture: "hud.png",
			X:       184,
			Y:       48,
			W:       24,
			H:       24,
		}.AsAnimation(),
		Targeting:      skill.TargetAnywhere,
		TargetingBrush: skill.SingleHex,
		Costs: map[skill.CostType]int{
			skill.CostsActionPoints: 65,
		},
		Effects: []skill.Effect{
			{
				When: skill.NewTimingFromPoint(skill.AttackApexTimingPoint),
				What: skill.DefileEffect{},
			},
			{
				When: skill.NewTimingFromPoint(skill.AttackApexTimingPoint),
				What: skill.SpawnParticipantEffect{
					Profession: "Skeleton",
					Level: []skill.Operation{
						{Operator: skill.AddOp, Variable: "1"},
						{Operator: skill.MultOp, Variable: "$DARK"},
					},
				},
			},
		},
	},
}
