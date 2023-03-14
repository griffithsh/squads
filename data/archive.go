package data

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/griffithsh/squads/embedded"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld"
	"github.com/griffithsh/squads/game/overworld/hbg"
	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/skill"
)

// Archive is a store of game data.
type Archive struct {
	overworldRecipes       []*overworld.Recipe
	skills                 skill.Map
	appearances            map[AppearanceKey]*game.Appearance
	hairColors             []string
	skinColors             []string
	names                  map[string][]string
	combatMaps             []game.CombatMapRecipe // eventually these would be keyed by their terrain in some way?
	images                 map[string]image.Image
	overworldBaseTiles     map[procedural.Code]hbg.BaseTile
	overworldEncroachments hbg.EncroachmentsCollection
}

// NewArchive constructs a new Archive.
func NewArchive() (*Archive, error) {
	archive := Archive{
		overworldRecipes:       []*overworld.Recipe{},
		skills:                 skill.Map{},
		appearances:            map[AppearanceKey]*game.Appearance{},
		names:                  map[string][]string{},
		images:                 map[string]image.Image{},
		overworldBaseTiles:     map[procedural.Code]hbg.BaseTile{},
		overworldEncroachments: hbg.EncroachmentsCollection{},
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

	// Load all the embedded contents in here.
	files, err := embedded.Filenames()
	if err != nil {
		return nil, fmt.Errorf("list embedded files: %v", err)
	}
	for _, filename := range files {
		b, err := embedded.Get(filename)
		if err != nil {
			return nil, fmt.Errorf("get file %q: %v", filename, err)
		}
		if err := archive.interpret(filename, bytes.NewReader(b)); err != nil {
			return nil, fmt.Errorf("interpret %q: %v", filename, err)
		}
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
			if err := a.interpret(head.Name, tr); err != nil {
				return fmt.Errorf("interpret %q: %v", head.Name, err)
			}
		}
	}

	return nil
}

func (a *Archive) interpret(filename string, r io.Reader) error {
	switch filepath.Ext(filename) {
	case ".names":
		// .names are a csv-like file that include names with tags for
		// how they should be used, like M and F for sexes etc.
		names, err := parseNames(r)
		if err != nil {
			return fmt.Errorf("parse %s.names: %v", filename, err)
		}
		for k, v := range names {
			a.names[k] = v
		}

	case ".overworld-recipe":
		recipe, err := overworld.ParseRecipe(r)
		if err != nil {
			return fmt.Errorf("parse %s: %v", filename, err)
		}
		a.overworldRecipes = append(a.overworldRecipes, recipe)

	case ".png":
		decoded, err := png.Decode(r)
		if err != nil {
			return fmt.Errorf("png.Decode %s: %s", filename, err)
		}
		a.images[filename] = decoded

	case ".skills":
		skills, err := parseSkills(r)
		if err != nil {
			return fmt.Errorf("parse %s: %v", filename, err)
		}
		for _, skill := range skills {
			if _, ok := a.skills[skill.ID]; ok {
				//duplicate - overwrites
				fmt.Fprintf(os.Stderr, "skill in %s overwrites skill ID %v", filename, skill.ID)
			}
			a.skills[skill.ID] = skill
			fmt.Println("loaded skill", skill.ID)
		}

	case ".terrain":
		dec := json.NewDecoder(r)
		var v game.CombatMapRecipe
		err := dec.Decode(&v)
		if err != nil {
			return fmt.Errorf("parse %s: %v", filename, err)
		}
		a.combatMaps = append(a.combatMaps, v)
	}

	switch {
	case strings.HasSuffix(filename, ".appearance"):
		dec := json.NewDecoder(r)
		var v struct {
			game.Appearance
			Sex        string
			Profession string
			HairColor  string
			SkinColor  string
		}
		err := dec.Decode(&v)
		if err != nil {
			return fmt.Errorf("parse %s: %v", filename, err)
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
			Hair:       v.HairColor,
			Skin:       v.SkinColor,
		}
		a.hairColors = append(a.hairColors, key.Hair)
		a.skinColors = append(a.skinColors, key.Skin)
		if _, ok := a.appearances[key]; ok {
			//stomp alert!
			return fmt.Errorf("duplicate appearance %v %s, %s-hair, %s skin", sex, v.Profession, v.HairColor, v.SkinColor)
		}
		a.appearances[key] = &v.Appearance
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

func (a *Archive) GetOverworldBaseTiles() map[procedural.Code]hbg.BaseTile {
	return a.overworldBaseTiles
}

func (a *Archive) GetOverworldEncroachments() hbg.EncroachmentsCollection {
	return a.overworldEncroachments
}
