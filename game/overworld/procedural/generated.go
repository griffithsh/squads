package procedural

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/squad"
)

type Generated struct {
	Recipe string
	Paths  Paths

	// PathExtents are only being exposed for debugging purposes.
	PathExtents map[geom.DirectionType]geom.Key
	Terrain     map[geom.Key]Code

	Opponents          map[geom.Key]squad.RecipeID
	GenerationDuration time.Duration
}

func (g *Generated) Complexity() int {
	return len(g.Paths.Nodes)
}

// JSON marshaling support

type jsonPaths struct {
	Algorithm string
	Seed      int64

	Start    geom.Key
	Goal     geom.Key
	Nodes    map[string]Placement
	Specials map[string][]geom.Key
}
type t struct {
	Paths              jsonPaths
	PathExtents        map[geom.DirectionType]geom.Key
	Terrain            map[string]Code
	Opponents          map[string]squad.RecipeID
	GenerationDuration time.Duration
}

func packKey(k geom.Key) string {
	return fmt.Sprintf("%d|%d", k.M, k.N)
}

func (g *Generated) MarshalJSON() ([]byte, error) {

	nodes := map[string]Placement{}
	for k, v := range g.Paths.Nodes {
		nodes[packKey(k)] = v
	}
	paths := jsonPaths{
		Algorithm: g.Paths.Algorithm,
		Seed:      g.Paths.Seed,
		Start:     g.Paths.Start,
		Goal:      g.Paths.Goal,
		Nodes:     nodes,
		Specials:  g.Paths.Specials,
	}
	opponents := map[string]squad.RecipeID{}
	for k, v := range g.Opponents {
		opponents[packKey(k)] = v
	}
	terrain := map[string]Code{}
	for key, value := range g.Terrain {
		terrain[packKey(key)] = value
	}
	v := t{
		Paths:              paths,
		PathExtents:        g.PathExtents,
		Terrain:            terrain,
		Opponents:          opponents,
		GenerationDuration: g.GenerationDuration,
	}
	return json.Marshal(v)
}

func unpackKey(k string) (geom.Key, error) {
	s := strings.Split(k, "|")
	m, err := strconv.Atoi(s[0])
	if err != nil {
		return geom.Key{}, fmt.Errorf("get M from %s: %v", k, err)
	}
	n, err := strconv.Atoi(s[1])
	if err != nil {
		return geom.Key{}, fmt.Errorf("get N from %s: %v", k, err)
	}
	return geom.Key{m, n}, nil
}

func (g *Generated) UnmarshalJSON(data []byte) error {
	v := t{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	nodes := map[geom.Key]Placement{}
	for k, v := range v.Paths.Nodes {
		key, err := unpackKey(k)
		if err != nil {
			return fmt.Errorf("unpack Node with key %q: %v", k, err)
		}
		nodes[key] = v
	}
	terrain := map[geom.Key]Code{}
	for k, v := range v.Terrain {
		key, err := unpackKey(k)
		if err != nil {
			return fmt.Errorf("unpack Terrain with key %q: %v", k, err)
		}
		terrain[key] = v
	}
	opponents := map[geom.Key]squad.RecipeID{}
	for k, v := range v.Opponents {
		key, err := unpackKey(k)
		if err != nil {
			return fmt.Errorf("unpack Opponent with key %q: %v", k, err)
		}
		opponents[key] = v
	}

	g.Paths = Paths{
		Algorithm: v.Paths.Algorithm,
		Seed:      v.Paths.Seed,
		Start:     v.Paths.Start,
		Goal:      v.Paths.Goal,
		Nodes:     nodes,
		Specials:  v.Paths.Specials,
	}
	g.PathExtents = v.PathExtents
	g.Terrain = terrain
	g.Opponents = opponents
	g.GenerationDuration = v.GenerationDuration
	return nil
}
