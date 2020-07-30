package overworld

import (
	"fmt"
	"math/rand"

	"github.com/griffithsh/squads/graph"

	"github.com/griffithsh/squads/baddy"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/squad"
)

// connect two Nodes by adding each to the other's Connected map.
func connect(n1 *Node, n2 *Node) error {
	if n1 == nil {
		return fmt.Errorf("cannot connect n1(nil) to n2(%v)", n2)
	}
	if n2 == nil {
		return fmt.Errorf("cannot connect n1(%v) to n2(nil)", n1)
	}
	n1Neighbors := n1.ID.Neighbors()
	dirOfN2, ok := n1Neighbors[n2.ID]
	if !ok {
		fmt.Printf("neighbors of n1 (%v): %v\n", n1.ID, n1Neighbors)
		fmt.Printf("neighbors of n2 (%v): %v\n", n2.ID, n2.ID.Neighbors())

		return fmt.Errorf("cannot connect n1 (%v) to n2 (%v): not neighbors", n1, n2)
	}

	dirOfN1 := geom.Opposite[dirOfN2]

	if n1.Connected == nil {
		n1.Connected = map[geom.DirectionType]geom.Key{}
	}
	n1.Connected[dirOfN2] = n2.ID

	if n2.Connected == nil {
		n2.Connected = map[geom.DirectionType]geom.Key{}
	}
	n2.Connected[dirOfN1] = n1.ID

	return nil
}

// recurse generates pathways between nodes by iterating through hexes adjacent to start.
func recurse(rng *rand.Rand, origin geom.Key, d *Map, paths map[KeyPair]struct{}) {
	dirs := []geom.DirectionType{geom.S, geom.SW, geom.NW, geom.N, geom.NE, geom.SE}
	rng.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})
	neighbors := origin.Adjacent()
	for _, dir := range dirs {
		// threshold is how likely this neighbor is to not continue in this
		// direction
		threshold := 0.0
		switch len(d.Nodes[origin].Connected) {
		case 0: // 100% chance of this node connecting to anything
		case 1: // 65% chance of this node connecting to 2 things
			threshold = 0.45
		case 2:
			threshold = 0.85
		case 3:
			threshold = 0.925
		default:
			threshold = 0.99999
		}
		if rng.Float64() > threshold {
			neigh := neighbors[dir]
			// if this direction is not in potentials, then pull out.

			if _, ok := paths[KeyPair{origin, neigh}]; !ok {
				continue
			}

			// Does this direction already have a node?
			if _, ok := d.Nodes[neigh]; ok {
				// We have a chance to continue, and not connect the nodes
				if rng.Float64() > 0.25 {
					continue
				}
			} else {
				d.Nodes[neigh] = &Node{ID: neigh}
			}
			err := connect(d.Nodes[origin], d.Nodes[neigh])
			if err != nil {
				fmt.Printf("connect %v to %v failed: %v\n", origin, neigh, err)
				continue
			}
			recurse(rng, neigh, d, paths)
		}
	}
}

// generate a Map from a Recipe and a base enemy level for opponents by calling
// randomness from rng.
func generate(rng *rand.Rand, recipe *Recipe, lvl int) Map {
	// TODO: use lvl to make the baddies stronger
	d := Map{
		Terrain: recipe.Terrain,
		Nodes:   map[geom.Key]*Node{},
		Enemies: map[geom.Key][]*game.Character{},
	}

	// Pick the first key of a random pair in the recipe's permissable paths.
	origin := recipe.Paths[rng.Intn(len(recipe.Paths)-1)].First
	d.Nodes[origin] = &Node{ID: origin}

	// unpack paths into a two-way set available paths.
	pathSet := map[KeyPair]struct{}{}
	keySet := map[geom.Key]struct{}{}
	for _, path := range recipe.Paths {
		pathSet[KeyPair{path.First, path.Second}] = struct{}{}
		pathSet[KeyPair{path.Second, path.First}] = struct{}{}

		keySet[path.First] = struct{}{}
		keySet[path.Second] = struct{}{}
	}

	// recurse through, generating a random path.
	recurse(rng, origin, &d, pathSet)

	// Roll for Points of interest that must be included in the Map's Nodes
	must := []geom.Key{}
	for _, poi := range recipe.Interesting {
		// Randomly select poi.Pick number of Options from each InterestRoll, by
		// constructing a shuffled sice of indices, then truncating the excess.
		indices := []int{}
		for i := range poi.Options {
			indices = append(indices, i)
		}
		rng.Shuffle(len(indices), func(i, j int) {
			indices[i], indices[j] = indices[j], indices[i]
		})
		indices = indices[:poi.Pick]

		for _, i := range indices {
			must = append(must, poi.Options[i])
		}
	}

	// Link every rolled point of interest by connecting it to the origin that
	// we started generating with.
	for _, poi := range must {
		if _, ok := d.Nodes[poi]; ok {
			// This poi is already in the existing nodes.
			continue
		}

		cost := func(_ graph.Vertex, v graph.Vertex) float64 {
			k := v.(geom.Key)
			if _, ok := d.Nodes[k]; ok {
				return 0.0
			}
			return 1.0
		}
		edge := func(v graph.Vertex) []graph.Vertex {
			k := v.(geom.Key)

			result := make([]graph.Vertex, 0, 6)
			neighbors := []geom.Key{
				k.ToN(),
				k.ToNE(),
				k.ToSE(),
				k.ToS(),
				k.ToSW(),
				k.ToNW(),
			}
			for _, k := range neighbors {
				if _, ok := keySet[k]; ok {
					result = append(result, k)
				}
			}
			return result
		}
		guess := func(v1, v2 graph.Vertex) float64 {
			squareDiff := func(a, b int) float64 {
				diff := a - b
				if diff < 0 {
					diff = -diff
				}
				return float64(diff * diff)
			}
			a, b := v1.(geom.Key), v2.(geom.Key)

			return squareDiff(a.M, b.M) + squareDiff(a.N, b.N)
		}
		steps := graph.NewSearcher(cost, edge, guess).Search(origin, poi)
		if steps == nil {
			// Points of interest should never be located in areas of the map
			// that are inaccessible? This can only panic on recipes with
			// non-contiguous Paths, right?

			panic(fmt.Sprintf("there was no path from %v to the PoI at %v", origin, poi))
		}

		var prev geom.Key
		prev = geom.Key{M: steps[0].V.(geom.Key).M, N: steps[0].V.(geom.Key).N}
		for i, step := range steps[1:] {
			k := geom.Key{M: step.V.(geom.Key).M, N: step.V.(geom.Key).N}

			// If this step is not already in nodes, add it.
			if _, ok := d.Nodes[k]; !ok {
				// Add a new node.
				current := Node{ID: k}
				d.Nodes[k] = &current

				// connect the new node to the previous step's node.
				err := connect(d.Nodes[prev], &current)
				if err != nil {
					fmt.Printf("connect step %d: %v\n", i, err)
				}
			}

			// Prepare for next loop.
			prev = k
		}
	}

	if len(d.Nodes) < 2 {
		fmt.Println("impossible Map generated: not enough room for both start and exit")
		// try again?
		return generate(rng, recipe, lvl)
	}
	// Sort then shuffle keys, so that the results of this function are
	// deterministic based on the provided PRNG.
	keys := d.SortedNodeKeys()
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	for i, key := range keys {
		if i == 0 {
			// First key is the player start.
			d.Start = key
			continue
		} else if i == 1 {
			// Second key is the exit gate.
			d.Gate = key
			continue
		}

		// TODO: The selection of squads should be controlled by the overworld
		// Recipe.
		if rng.Intn(3) == 0 { // 1 in 3 chance of adding an enemy squad here.
			rollCharacters := func(rng *rand.Rand, id squad.RecipeID) []*game.Character {
				result := []*game.Character{}
				for _, recipeID := range squad.Recipes[id].Construct(rng) {
					char := baddy.Recipes[recipeID].Construct(rng)
					char.Level = lvl
					vit := int(char.VitalityPerLevel * float64(char.Level))
					hp := char.BaseHealth
					char.CurrentHealth = game.MaxHealth(hp, vit)
					result = append(result, char)
				}
				return result
			}
			switch rng.Intn(4) {
			case 0:
				d.Enemies[key] = rollCharacters(rng, squad.SoloSkellington)
			case 1:
				d.Enemies[key] = rollCharacters(rng, squad.WolfPack1)
			case 2:
				d.Enemies[key] = rollCharacters(rng, squad.SoloNecro)
			case 3:
				d.Enemies[key] = rollCharacters(rng, squad.NecroCohort)
			}
		}
	}

	return d
}
