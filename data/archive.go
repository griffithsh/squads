package data

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld"
	"github.com/griffithsh/squads/skill"
)

// Archive is a store of game data.
type Archive struct {
	overworldRecipes []*overworld.Recipe
	skills           skill.Map
	performances     map[PerformanceKey]*game.PerformanceSet
}

// NewArchive constructs a new Archive.
func NewArchive() (*Archive, error) {
	archive := Archive{
		overworldRecipes: []*overworld.Recipe{},
		skills:           skill.Map{},
		performances:     map[PerformanceKey]*game.PerformanceSet{},
	}
	for _, sd := range internalSkills {
		archive.skills[sd.ID] = sd
	}

	for k, v := range internalPerformances {
		archive.performances[k] = v
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
					return fmt.Errorf("parse %s.overworld-recipe: %v", head.Name, err)
				}
				a.overworldRecipes = append(a.overworldRecipes, recipe)
			case ".performance-set":
				dec := json.NewDecoder(tr)
				var v game.PerformanceSet
				err := dec.Decode(&v)
				if err != nil {
					return fmt.Errorf("parse %s.performance-set: %v", head.Name, err)
				}
				for _, sex := range v.Sexes {
					a.performances[PerformanceKey{sex, head.Name}] = &v
				}
			case ".skill-thingy?":
				// TODO:
			}
		}
	}

	return nil
}

// GetRecipes returns overworld recipes.
func (a *Archive) GetRecipes() []*overworld.Recipe {
	return a.overworldRecipes
}
