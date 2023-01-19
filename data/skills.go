package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/item"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/skill"
	"github.com/griffithsh/squads/targeting"
)

type targetingJSON struct {
	Selectable struct {
		Type     targeting.SelectableType
		MinRange int
		MaxRange int
	}
	Brush struct {
		Type            targeting.BrushType
		MinRange        int
		MaxRange        int
		LinearExtent    int
		LinearDirection geom.RelativeDirection
	}
}

type skillEffect struct {
	When      int // milliseconds
	WhenPoint string
	What      []interface{}
}

// skillDescription is the raw format from a .skills file.
type skillDescription struct {
	ID          skill.ID
	Name        string
	Explanation string

	// Tags critically includes Attack or Spell, and allows the game to select
	// an appropriate animation to use when using the skill.
	Tags []string

	Icon game.Sprite

	Targeting targetingJSON

	// Effects of triggering this skill.
	Effects []skillEffect

	Costs map[string]int

	// AttackChanceToHitModifier multiplies the base chance to hit of the skill.
	// A value of zero does not modify the chance to hit. A value of 0.1
	// improves the chance to hit by 10%. A value of -0.5 halves the chance to
	// hit.
	AttackChanceToHitModifier float64
}

// convert to a skill.Description
func (sd *skillDescription) convert() (skill.Description, error) {
	tags := []skill.Classification{}
	for _, raw := range sd.Tags {
		tag, err := skill.ClassificationString(raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unknown skill.Classification %s", raw)
			continue
		}
		tags = append(tags, tag)
	}

	targetingRule := targeting.Rule{
		Selectable: targeting.Selectable{
			Type:     sd.Targeting.Selectable.Type,
			MinRange: sd.Targeting.Selectable.MinRange,
			MaxRange: sd.Targeting.Selectable.MaxRange,
		},
		Brush: targeting.Brush{
			Type:            sd.Targeting.Brush.Type,
			MinRange:        sd.Targeting.Brush.MinRange,
			MaxRange:        sd.Targeting.Brush.MaxRange,
			LinearExtent:    sd.Targeting.Brush.LinearExtent,
			LinearDirection: sd.Targeting.Brush.LinearDirection,
		},
	}

	effects := []skill.Effect{}
	for _, raw := range sd.Effects {
		var when skill.Timing
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
		ID:                        sd.ID,
		Name:                      sd.Name,
		Explanation:               sd.Explanation,
		Tags:                      tags,
		Icon:                      *sd.Icon.AsAnimation(),
		Targeting:                 targetingRule,
		Effects:                   effects,
		Costs:                     costs,
		AttackChanceToHitModifier: sd.AttackChanceToHitModifier,
	}, nil
}

func (d *skillEffect) UnmarshalJSON(data []byte) error {
	v := struct {
		When      int
		WhenPoint string
		What      []struct {
			Type string `json:"_type"`
		}
	}{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.When = v.When
	d.WhenPoint = v.WhenPoint

	v2 := struct {
		What []json.RawMessage
	}{}
	if err := json.Unmarshal(data, &v2); err != nil {
		return err
	}

	var whats []interface{}
	for i, typer := range v.What {
		ty := typer.Type
		switch ty {
		case "DamageEffect":
			var what skill.DamageEffect
			if err := json.Unmarshal(v2.What[i], &what); err != nil {
				return fmt.Errorf("unmarshal DamageEffect: %q\n\n%s", err, v2.What[i])
			}
			whats = append(whats, what)

		case "InjuryEffect":
			what := struct {
				InjuryType string `json:"type"`
				Value      int
			}{}
			if err := json.Unmarshal(v2.What[i], &what); err != nil {
				return fmt.Errorf("unmarshal InjuryEffect: %q\n\n%s", err, v2.What[i])
			}
			injuryType := skill.InjuryTypeFromString(what.InjuryType)
			if injuryType == nil {
				return fmt.Errorf("unknown InjuryType %q", what.InjuryType)
			}
			whats = append(whats, skill.InjuryEffect{
				Type:  *injuryType,
				Value: what.Value,
			})

		// case "HealEffect":
		// case "ReviveEffect":
		// case "DefileEffect":
		// case "SpawnParticipantEffect":
		default:
			// TODO: other effect types!
			return fmt.Errorf("unknown _type %q (%v)", ty, v)
		}
	}
	d.What = whats
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
func (a *Archive) SkillsByWeaponClass(weap item.Class) []*skill.Description {
	// FIXME: implementation
	switch weap {
	case item.SwordClass:
		return []*skill.Description{
			a.Skill("debug-basic-attack"),
		}
	case item.UnarmedClass:
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
			Texture: "combat/hud.png",
			X:       232,
			Y:       24,
			W:       24,
			H:       24,
		}.AsAnimation(),
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
			Texture: "combat/hud.png",
			X:       160,
			Y:       0,
			W:       24,
			H:       24,
		}.AsAnimation(),
		Targeting: targeting.Rule{
			Selectable: targeting.Selectable{
				Type:     targeting.SelectWithin,
				MinRange: 1,
				MaxRange: 1,
			},
			Brush: targeting.Brush{
				Type: targeting.SingleHex,
			},
		},
		Costs: map[skill.CostType]int{
			skill.CostsActionPoints: 20,
		},
		Effects: []skill.Effect{
			{
				When: skill.NewTiming(time.Millisecond * 500),
				What: []interface{}{skill.DamageEffect{
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
	},
	{
		ID:          "debug-lightning",
		Name:        "Mage Lightning",
		Explanation: "A lightning bolt strikes the target dealing 1-10 damage",
		Tags:        []skill.Classification{skill.Spell},
		Icon: *game.Sprite{
			Texture: "combat/hud.png",
			X:       160,
			Y:       24,
			W:       24,
			H:       24,
		}.AsAnimation(),
		Targeting: targeting.Rule{
			Selectable: targeting.Selectable{
				Type: targeting.SelectAnywhere,
			},
			Brush: targeting.Brush{
				Type: targeting.SingleHex,
			},
		},
		Costs: map[skill.CostType]int{
			skill.CostsActionPoints: 45,
		},
		Effects: []skill.Effect{
			{
				When: skill.NewTimingFromPoint(skill.AttackApexTimingPoint),
				What: []interface{}{skill.DamageEffect{
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
				}},
			},
		},
	},
	{
		ID:          "debug-revive",
		Name:        "Pheonix form",
		Explanation: "A pheonix feather lands on the target, reviving it.",
		Tags:        []skill.Classification{skill.Spell},
		Icon: *game.Sprite{
			Texture: "combat/hud.png",
			X:       160,
			Y:       48,
			W:       24,
			H:       24,
		}.AsAnimation(),
		Targeting: targeting.Rule{
			Selectable: targeting.Selectable{
				Type: targeting.SelectAnywhere,
			},
			Brush: targeting.Brush{
				Type: targeting.SingleHex,
			},
		},
		Costs: map[skill.CostType]int{
			skill.CostsActionPoints: 45,
		},
		Effects: []skill.Effect{
			{
				When: skill.NewTimingFromPoint(skill.AttackApexTimingPoint),
				What: []interface{}{
					skill.ReviveEffect{},
					skill.HealEffect{
						Amount:       0.15,
						IsPercentage: true,
					},
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
			Texture: "combat/hud.png",
			X:       184,
			Y:       48,
			W:       24,
			H:       24,
		}.AsAnimation(),
		Targeting: targeting.Rule{
			Selectable: targeting.Selectable{
				Type: targeting.SelectAnywhere,
			},
			Brush: targeting.Brush{
				Type: targeting.SingleHex,
			},
		},
		Costs: map[skill.CostType]int{
			skill.CostsActionPoints: 65,
		},
		Effects: []skill.Effect{
			{
				When: skill.NewTimingFromPoint(skill.AttackApexTimingPoint),
				What: []interface{}{
					skill.DefileEffect{},
					skill.SpawnParticipantEffect{
						Profession: "Skeleton",
						Level: []skill.Operation{
							{Operator: skill.AddOp, Variable: "1"},
							{Operator: skill.MultOp, Variable: "$DARK"},
						},
					},
				},
			},
		},
	},
}
