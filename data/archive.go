package data

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld"
	"github.com/griffithsh/squads/skill"
)

// Archive is a store of game data.
type Archive struct {
	overworldRecipes []*overworld.Recipe
	skills           skill.Map
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
			skill.DamageEffect{
				Min: []skill.Operation{
					{Operator: skill.AddOp, Variable: "$DMG-MIN"},
				},
				Max: []skill.Operation{
					{Operator: skill.AddOp, Variable: "$DMG-MAX"},
				},
				Classification: skill.Attack,
				ScheduleTime:   200 * time.Millisecond,
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
			skill.DamageEffect{
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
				ScheduleTime:   1200 * time.Millisecond,
			},
		},
	},
}

// NewArchive constructs a new Archive.
func NewArchive() (*Archive, error) {
	archive := Archive{
		overworldRecipes: []*overworld.Recipe{},
		skills:           skill.Map{},
	}
	for _, sd := range internalSkills {
		archive.skills[sd.ID] = sd
	}
	return &archive, nil
}

// Load data into the Archive from a tar.gz archive, replacing data already in
// the Archive if the provided files share the same filenames.
func (a *Archive) Load(r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("new reader: %v", err)
	}

	tr := tar.NewReader(gzr)

	for {
		head, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("read next file from tar: %v", err)
		}

		if head.Typeflag == tar.TypeReg {
			switch filepath.Ext(head.Name) {
			case ".overworld-recipe":
				recipe, err := overworld.ParseRecipe(tr)
				if err != nil {
					return fmt.Errorf("parse overworld recipe from %s: %v", head.Name, err)
				}
				a.overworldRecipes = append(a.overworldRecipes, recipe)
			}
		}
	}

	return nil
}

// GetRecipes returns overworld recipes.
func (a *Archive) GetRecipes() []*overworld.Recipe {
	return a.overworldRecipes
}

func (a *Archive) SkillsByProfession(prof game.CharacterProfession) []*skill.Description {
	// FIXME: implementation
	return []*skill.Description{
		a.Skill("debug-basic-attack"),
		a.Skill("debug-lightning"),
	}
}

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
