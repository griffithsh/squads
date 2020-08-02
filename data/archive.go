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
	appearances      map[AppearanceKey]*game.Appearance
	hairColors       []string
	skinColors       []string
	names            map[string][]string
	combatMaps       []game.CombatMapRecipe // eventually these would be keyed by their terrain in some way?
	images           map[string]image.Image
}

// NewArchive constructs a new Archive.
func NewArchive() (*Archive, error) {
	archive := Archive{
		overworldRecipes: []*overworld.Recipe{},
		skills:           skill.Map{},
		appearances:      map[AppearanceKey]*game.Appearance{},
		names:            map[string][]string{},
		images:           map[string]image.Image{},
	}
	for _, sd := range internalSkills {
		archive.skills[sd.ID] = sd
	}

	for k, v := range internalAppearances {
		archive.appearances[k] = v
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

	hairColors := map[string]struct{}{}
	skinColors := map[string]struct{}{}
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
			} else if strings.HasPrefix(head.Name, "character-appearance/") {
				switch {
				case strings.HasSuffix(head.Name, ".variations.png"):
					decoded, err := png.Decode(tr)
					if err != nil {
						return fmt.Errorf("png.Decode %s: %s", head.Name, err)
					}
					a.images[head.Name] = decoded
				case strings.HasSuffix(head.Name, ".appearance"):
					dec := json.NewDecoder(tr)
					var v struct {
						game.Appearance
						Sex        string
						Profession string
						Haircolor  string
						Skincolor  string
					}
					err := dec.Decode(&v)
					if err != nil {
						return fmt.Errorf("parse %s: %v", head.Name, err)
					}
					var sex game.CharacterSex
					if v.Sex == "XX" {
						sex = game.Female
					} else if v.Sex == "XY" {
						sex = game.Male
					} else {
						return fmt.Errorf("unknown Sex %s", v.Sex)
					}

					key := AppearanceKey{
						Sex:        sex,
						Profession: v.Profession,
						Hair:       v.Haircolor,
						Skin:       v.Skincolor,
					}
					hairColors[key.Hair] = struct{}{}
					skinColors[key.Skin] = struct{}{}
					a.appearances[key] = &v.Appearance
				}
			} else {
				switch filepath.Ext(head.Name) {
				case ".overworld-recipe":
					recipe, err := overworld.ParseRecipe(tr)
					if err != nil {
						return fmt.Errorf("parse %s: %v", head.Name, err)
					}
					a.overworldRecipes = append(a.overworldRecipes, recipe)
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

	for color := range hairColors {
		a.hairColors = append(a.hairColors, color)
	}

	for color := range skinColors {
		a.skinColors = append(a.skinColors, color)
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
