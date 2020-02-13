package data

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"

	"github.com/griffithsh/squads/game/overworld"
)

// Archive is a store of game data.
type Archive struct {
	overworldRecipes []*overworld.Recipe
}

// NewArchive constructs a new Archive.
func NewArchive() (*Archive, error) {
	return &Archive{
		overworldRecipes: []*overworld.Recipe{},
	}, nil
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
