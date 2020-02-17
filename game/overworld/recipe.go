package overworld

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/griffithsh/squads/geom"
)

// trim whitespace.
func trim(s string) string {
	return strings.Trim(s, "\t \n")
}

// KeyPair is a pair of geom.Keys that represents a linkage between two
// overworld nodes.
type KeyPair struct {
	First  geom.Key
	Second geom.Key
}

// InterestRoll contains a slice of Options to pick from, and a number of times
// to pick.
type InterestRoll struct {
	Pick    int
	Options []geom.Key
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
	Interesting []InterestRoll
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

		if strings.HasPrefix(line, "interesting:") {
			interesting, err := parseInteresting(br)
			if err != nil {
				return nil, fmt.Errorf("parseInteresting: %v", err)
			}
			result.Interesting = interesting
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

func parseInteresting(r *bufio.Reader) ([]InterestRoll, error) {
	var values string
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("read line: %v", err)
		}

		line = trim(line)
		if line == "" {
			break
		}

		values += line

	}
	re := regexp.MustCompile(`\d*? \([^)]*\)`)
	interestRolls := re.FindAllString(values, -1)

	results := make([]InterestRoll, 0, len(interestRolls))
	for i, raw := range interestRolls {
		raw = strings.Trim(raw, "\t \n")
		splits := strings.SplitN(raw, " ", 2)
		pick, rest := trim(splits[0]), strings.Trim(splits[1], "()\t \n,")
		ipick, err := strconv.Atoi(pick)
		if err != nil {
			return nil, fmt.Errorf("non-integer value for pick \"%s\" in Interesting %d (%s)", pick, i, raw)
		}

		result := InterestRoll{
			Pick: ipick,
		}
		keys := strings.Split(rest, ",")
		for j, key := range keys {
			key = trim(key)

			parts := strings.Split(key, " ")
			m, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("non-integer value for M in Option %d \"%s\" in Interesting %d (%s)", j, parts[0], i, raw)
			}
			n, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("non-integer value for N in Option %d \"%s\" in Interesting %d (%s)", j, parts[1], i, raw)
			}

			result.Options = append(result.Options, geom.Key{M: m, N: n})
		}

		results = append(results, result)
	}
	return results, nil
}
