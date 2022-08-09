package procedural

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

// Generator holds configuration data to construct overworld maps.
type Generator struct {
	MakePaths pathFunc       `json:"pathGeneration"`
	Terrain   TerrainBuilder `json:"terrainBuilder"`
}

func (g *Generator) UnmarshalJSON(b []byte) error {
	var v struct {
		PathGeneration pathFunc
		TerrainBuilder json.RawMessage
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	g.MakePaths = v.PathGeneration

	var ty struct {
		Value string `json:"type"`
	}
	err = json.Unmarshal(v.TerrainBuilder, &ty)
	if err != nil {
		return fmt.Errorf("unmarshal terrainBuilder: %v", err)
	}

	switch ty.Value {
	case "LinearGradientTerrainStrategy":
		var tb LinearGradientTerrainStrategy
		if err = json.Unmarshal(v.TerrainBuilder, &tb); err != nil {
			return fmt.Errorf("unmarshal LinearGradientTerrainStrategy: %v", err)
		}
		g.Terrain = &tb
	default:
		return fmt.Errorf("unknown terrainBuilder.type value: %s", ty.Value)
	}

	return nil
}

// Placement holds info about what roads or paths have been placed on a Key.
type Placement struct {
	Connections map[geom.DirectionType]struct{}
}

type Generated struct {
	Paths Paths

	// PathExtents are only being exposed for debugging purposes.
	PathExtents map[geom.DirectionType]geom.Key
	Terrain     map[geom.Key]Code
}

// Generate should take a recipe and output an overworld map.
func (g Generator) Generate(seed int64, level int) Generated {
	prng := rand.New(rand.NewSource(seed))

	paths := Paths{}
	for i := 5; i >= 0; i-- {
		var err error
		paths, err = g.MakePaths(seed, level)
		if err == nil {
			break
		}
		fmt.Printf("path generation failed for bad seed %d: %v\n", seed, err)
		seed = prng.Int63()
		if i == 0 {
			panic(fmt.Sprintf("path generation retries exhausted: %v", err))
		}
	}

	terrainCodes := g.Terrain.Build(prng, paths)

	// overwrite standard terrain with doodads
	// TODO: ...

	// TODO: baddies

	// TODO: misc other encounters

	// return an object that contains info on how to render the overworld as
	// well as programmatic info on what hexes are navigable.
	return Generated{
		Paths:       paths,
		PathExtents: extentsOf(keysOf(paths)),
		Terrain:     terrainCodes,
	}
}
