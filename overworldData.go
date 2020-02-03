package main

import (
	"fmt"
	"math/rand"

	"github.com/griffithsh/squads/squad"

	"github.com/griffithsh/squads/game/overworld"
	"github.com/griffithsh/squads/geom"
)

/*
 * The idea is for this to move into the overworld package when I have something
 * semi-solid as to an implementation.
 */

// connect two Nodes by adding each to the other's Connected map.
func connect(n1 *overworld.Node, n2 *overworld.Node) error {
	n1Neighbors := n1.ID.Neighbors()
	dirOfN2, ok := n1Neighbors[n2.ID]
	if !ok {
		fmt.Printf("neighbors of n1 (%v): %v\n", n1.ID, n1Neighbors)
		fmt.Printf("neighbors of n2 (%v): %v\n", n2.ID, n2.ID.Neighbors())

		return fmt.Errorf("n1 (%v) is not a neighbor of n2 (%v)", n1, n2)
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

// available generates random terrain.
func available() map[geom.Key]overworld.TileID {
	result := make(map[geom.Key]overworld.TileID)

	for m := 0; m < 4; m++ {
		for n := 0; n < 16; n++ {
			id := overworld.TileID((n * m) % 3)
			result[geom.Key{M: m, N: n}] = id
		}
	}
	return result
}

// recurse generates pathways between nodes by iterating through hexes adjacent to start.
func recurse(rng *rand.Rand, start geom.Key, d *overworld.Map, potentials map[geom.Key]overworld.TileID) {
	dirs := []geom.DirectionType{geom.S, geom.SW, geom.NW, geom.N, geom.NE, geom.SE}
	rng.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})
	neighbors := start.Adjacent()
	for _, dir := range dirs {
		// threshold is how likely this neighbor is to not continue in this
		// direction
		threshold := 0.0
		switch len(d.Nodes[start].Connected) {
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
			if _, ok := potentials[neigh]; !ok {
				continue
			}

			// Does this direction already have a node?
			if _, ok := d.Nodes[neigh]; ok {
				// We have a chance to continue, and not connect the nodes
				if rng.Float64() > 0.25 {
					continue
				}
			} else {
				d.Nodes[neigh] = &overworld.Node{ID: neigh}
			}
			err := connect(d.Nodes[start], d.Nodes[neigh])
			if err != nil {
				fmt.Printf("connect %v to %v failed: %v\n", start, neigh, err)
				continue
			}
			recurse(rng, neigh, d, potentials)
		}
	}
}

// data generates a Map from a Recipe by calling randomness from rng.
func data(rng *rand.Rand, recipe overworld.Recipe) overworld.Map {
	d := overworld.Map{
		Terrain: recipe.Terrain,
		Nodes:   map[geom.Key]*overworld.Node{},
		Enemies: map[geom.Key]squad.RecipeID{},
	}

	origin := geom.Key{M: 0, N: 0}
	d.Nodes[origin] = &overworld.Node{ID: origin}
	recurse(rng, origin, &d, recipe.Terrain)

	if len(d.Nodes) < 2 {
		panic("impossible Map generated: not enough room for both start and exit")
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

		// FIXME: The selection of squads should be controlled by the
		// overworld Recipe.
		// 1 in 3 chance of adding an enemy squad here.
		if rng.Intn(3) == 0 {
			switch rng.Intn(3) {
			case 0:
				d.Enemies[key] = squad.SoloSkellington
			case 1:
				d.Enemies[key] = squad.WolfPack1
			case 2:
				d.Enemies[key] = squad.SoloGiant
			}
		}
	}

	return d
}
