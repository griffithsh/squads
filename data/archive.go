package data

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld"
	"github.com/griffithsh/squads/skill"
)

// Archive is a store of game data.
type Archive struct {
	overworldRecipes []*overworld.Recipe
	skills           skill.Map
	performances     map[PerformanceKey]*game.PerformanceSet
	names            map[string][]string
	combatMaps       []game.CombatMapRecipe // eventually these would be keyed by their terrain in some way?
	images           map[string]image.Image
}

// NewArchive constructs a new Archive.
func NewArchive() (*Archive, error) {
	archive := Archive{
		overworldRecipes: []*overworld.Recipe{},
		skills:           skill.Map{},
		performances:     map[PerformanceKey]*game.PerformanceSet{},
		names:            map[string][]string{},
		images:           map[string]image.Image{},
	}
	for _, sd := range internalSkills {
		archive.skills[sd.ID] = sd
	}

	for k, v := range internalPerformances {
		archive.performances[k] = v
	}
	for k, v := range internalNames {
		archive.names[k] = v
	}

	archive.combatMaps = append(archive.combatMaps, internalCombatMaps...)

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
			if strings.HasPrefix(head.Name, "combat-terrain/") {
				fmt.Println("found file in combat-terrain", head.Name)
				switch filepath.Ext(head.Name) {
				case ".png":
					decoded, err := png.Decode(tr)
					if err != nil {
						return fmt.Errorf("png.Decode %s: %s", head.Name, err)
					}
					a.images[head.Name] = decoded

				case ".terrain":
					dec := json.NewDecoder(tr)
					var v game.CombatMapRecipe
					err := dec.Decode(&v)
					if err != nil {
						return fmt.Errorf("parse %s: %v", head.Name, err)
					}
					a.combatMaps = append(a.combatMaps, v)
				}
			} else {
				switch filepath.Ext(head.Name) {
				case ".overworld-recipe":
					recipe, err := overworld.ParseRecipe(tr)
					if err != nil {
						return fmt.Errorf("parse %s: %v", head.Name, err)
					}
					a.overworldRecipes = append(a.overworldRecipes, recipe)
				case ".performance-set":
					dec := json.NewDecoder(tr)
					v := defaultPerformanceSet()
					err := dec.Decode(v)
					if err != nil {
						return fmt.Errorf("parse %s: %v", head.Name, err)
					}
					for _, sex := range v.Sexes {
						a.performances[PerformanceKey{sex, v.Name}] = v
					}
				case ".skill-thingy?":
					// TODO:
				case ".names":
					// .names are a csv-like file that include names with tags for
					// how they should be used, like M and F for sexes etc.
					names, err := parseNames(tr)
					if err != nil {
						return fmt.Errorf("parse %s.names: %v", head.Name, err)
					}
					for k, v := range names {
						a.names[k] = v
					}
				}

			}
		}
	}

	return nil
}

// GetRecipes returns overworld recipes.
func (a *Archive) GetRecipes() []*overworld.Recipe {
	return a.overworldRecipes
}

// GetImage returns an image that has been loaded dynamically into the archive.
func (a *Archive) GetImage(name string) (val image.Image, ok bool) {
	val, ok = a.images[name]
	return val, ok
}
