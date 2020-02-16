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

// KeyPair is a pair of geom.Keys that represents a linkage between two
// overworld nodes.
type KeyPair struct {
	First  geom.Key
	Second geom.Key
}

// Recipe stores configuration of how to roll a map.
type Recipe struct {
	Label string

	// Terrain stores the visible tiles of an overworld.
	Terrain map[geom.Key]TileID

	// Paths between nodes that are permitted to be generated.
	Paths []KeyPair

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

		if strings.HasPrefix(line, "paths:") {
			paths, err := parsePaths(br)
			if err != nil {
				return nil, fmt.Errorf("parsePaths: %v", err)
			}
			result.Paths = paths
		}

	}
	if result.Label == "" {
		return nil, errors.New("missing recipe label")
	}
	if len(result.Terrain) == 0 {
		return nil, errors.New("missing recipe terrain")
	}
	if len(result.Paths) == 0 {
		return nil, errors.New("missing recipe paths")
	}
	return &result, nil
}

func parseTerrain(r *bufio.Reader) (map[geom.Key]TileID, error) {
	// values accumulates comma separated M, N, and TileID integer triplets.
	values := ""
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("read line: %v", err)
		}

		line = strings.Trim(line, "\t \n")
		if line == "" {
			break
		}

		values += line
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

func parsePaths(r *bufio.Reader) ([]KeyPair, error) {
	var values string
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("read line: %v", err)
		}

		line = strings.Trim(line, "\t \n")
		if line == "" {
			break
		}

		values += line
	}
	paths := strings.Split(values, ",")
	result := []KeyPair{}
	for i, path := range paths {
		path = strings.Trim(path, "\t \n")
		if path == "" {
			continue
		}
		raw := strings.SplitN(path, " ", 4)

		m1, err := strconv.Atoi(raw[0])
		if err != nil {
			return nil, fmt.Errorf("non-integer value for m1 \"%s\" in tile number %d (\"%s\")", raw[0], i, path)
		}
		n1, err := strconv.Atoi(raw[1])
		if err != nil {
			return nil, fmt.Errorf("non-integer value for n1 \"%s\" in tile number %d (\"%s\")", raw[1], i, path)
		}
		m2, err := strconv.Atoi(raw[2])
		if err != nil {
			return nil, fmt.Errorf("non-integer value for m2 \"%s\" in tile number %d (\"%s\")", raw[2], i, path)
		}
		n2, err := strconv.Atoi(raw[3])
		if err != nil {
			return nil, fmt.Errorf("non-integer value for n2 \"%s\" in tile number %d (\"%s\")", raw[3], i, path)
		}
		result = append(result, KeyPair{
			First:  geom.Key{M: m1, N: n1},
			Second: geom.Key{M: m2, N: n2},
		})
	}

	return result, nil
}
