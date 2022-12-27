package procedural

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/squad"
)

// Generator holds configuration data to construct overworld maps.
type Generator struct {
	Name             string
	MakePaths        pathFunc        `json:"pathGeneration"`
	Terrain          TerrainBuilder  `json:"terrainBuilder"`
	TerrainOverrides map[string]Code `json:"terrainSpecialOverrides"`

	Baddies OpponentSquads
}

func (g *Generator) UnmarshalJSON(b []byte) error {
	var v struct {
		Name                    string
		PathGeneration          pathFunc
		TerrainBuilder          json.RawMessage
		TerrainSpecialOverrides map[string]Code
		Baddies                 OpponentSquads
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	g.Name = v.Name
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
	case "RadialGradientTerrainStrategy":
		var tb RadialGradientTerrainStrategy
		if err = json.Unmarshal(v.TerrainBuilder, &tb); err != nil {
			return fmt.Errorf("unmarshal RadialGradientTerrainStrategy: %v", err)
		}
		g.Terrain = &tb
	case "NoiseTerrainStrategy":
		var n NoiseTerrainStrategy
		if err = json.Unmarshal(v.TerrainBuilder, &n); err != nil {
			return fmt.Errorf("unmarshal NoiseTerrainStrategy: %v", err)
		}
		g.Terrain = &n
	default:
		return fmt.Errorf("unknown terrainBuilder.type value: %s", ty.Value)
	}
	g.TerrainOverrides = v.TerrainSpecialOverrides

	g.Baddies = v.Baddies

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

	Opponents          map[geom.Key]squad.RecipeID
	GenerationDuration time.Duration
}

func (g *Generated) Complexity() int {
	return len(g.Paths.Nodes)
}

// Generate should take a recipe and output an overworld map.
func (g Generator) Generate(seed int64, level int) Generated {
	start := time.Now()
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

	// TODO: overwrite standard terrain with doodads

	// TODO: misc other encounters - treasures, recruitable characters,
	// merchants, rest stops, etc

	// Return an object that contains info on how to render the overworld as
	// well as programmatic info on what hexes are navigable.
	return Generated{
		Paths:              paths,
		PathExtents:        extentsOf(keysOf(paths.Nodes)),
		Terrain:            terrainCodes,
		Opponents:          g.Baddies.Generate(prng, paths, terrainCodes),
		GenerationDuration: time.Since(start),
	}
}
