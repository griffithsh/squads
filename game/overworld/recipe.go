package overworld

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/griffithsh/squads/geom"
)

// Recipe stores configuration of how to roll a map.
type Recipe struct {
	Label string

	// Terrain stores the visible tiles of an overworld.
	Terrain map[geom.Key]TileID

	// Interesting stores locations in the map that should always be included as
	// part of the generated Nodes.
	Interesting []int
}

// ParseRecipe parses the bytes of a recipe file into a Recipe.
func ParseRecipe(r io.Reader) (*Recipe, error) {
	result := Recipe{}
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			// give up?
			break
		}

		if strings.HasPrefix(line, "label:") {
			// this is the name of the recipe!
			result.Label = strings.Trim(strings.TrimLeft(line, "label:"), "\t \n")
		}
		if strings.HasPrefix(line, "terrain:") {
			terrain, err := parseTerrain(br)
			if err != nil {
				return nil, fmt.Errorf("parseTerrain: %v", err)
			}
			result.Terrain = terrain
		}

	}
	if result.Label == "" || len(result.Terrain) == 0 {
		return nil, errors.New("no recipe data found")
	}
	return &result, nil
}

func parseTerrain(r *bufio.Reader) (map[geom.Key]TileID, error) {
	// values accumulates comma separated M, N, and TileID integer triplets.
	values := ""
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read line: %v", err)
		}

		if strings.Trim(line, "\t \n") == "" {
			break
		}

		values += strings.Trim(line, "\t \n")
	}
	tiles := strings.Split(values, ",")

	result := make(map[geom.Key]TileID)
	for i, tile := range tiles {
		// Expecting tile = " 0 0 1\n" or similar.
		tile = strings.Trim(tile, "\t \n")
		if tile == "" {
			continue
		}

		raw := strings.SplitN(tile, " ", 3)
		m, err := strconv.Atoi(raw[0])
		if err != nil {
			return nil, fmt.Errorf("non-integer value for m \"%s\" in tile number %d (\"%s\")", raw[0], i, tile)
		}
		n, err := strconv.Atoi(raw[1])
		if err != nil {
			return nil, fmt.Errorf("non-integer value for n \"%s\" in tile number %d (\"%s\")", raw[1], i, tile)
		}
		tileID, err := strconv.Atoi(raw[2])
		if err != nil {
			return nil, fmt.Errorf("non-integer value for tile id \"%s\" in tile number %d (\"%s\")", raw[2], i, tile)
		}

		result[geom.Key{M: m, N: n}] = TileID(tileID)

	}
	return result, nil
}
