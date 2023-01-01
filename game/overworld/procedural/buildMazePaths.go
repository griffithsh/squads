package procedural

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/griffithsh/squads/geom"
)

func buildMazePaths(seed int64, level int) (Paths, error) {
	prng := rand.New(rand.NewSource(seed))
	// TODO: Tunables - passed level should drive these values.
	seeds := 10
	goalComplexity := 75
	expandMin := 3
	expandMax := 12

	start := time.Now()
	availableSeeds := geom.Key{M: 0, N: 0}.ExpandBy(expandMin, expandMax)
	sort.Slice(availableSeeds, func(i, j int) bool {
		if availableSeeds[i].M == availableSeeds[j].M {
			return availableSeeds[i].N < availableSeeds[j].N
		}
		return availableSeeds[i].M < availableSeeds[j].M
	})
	seedsSet := map[geom.Key]struct{}{}

	for len(seedsSet) < seeds {
		index := rand.Intn(len(availableSeeds))
		seedsSet[availableSeeds[index]] = struct{}{}
	}
	fmt.Printf("Seed: %v (%v)\n", prng, seedsSet)

	// toHandle is the queue of hexes we still need to connect to stuff.
	toHandle := shuffledGeomKeys(prng, seedsSet)

	placed := map[geom.Key]Placement{}

	for len(toHandle) > 0 {
		curr := toHandle[0]
		toHandle = toHandle[1:]
		if _, ok := placed[curr]; ok {
			fmt.Printf("\nskipping %v, already placed\n", curr)
			continue
		}
		fmt.Printf("\ndealing with %v -- complexity is %v\n", curr, len(placed))
		toPlace := Placement{
			Connections: map[geom.DirectionType]struct{}{},
		}
		adj := shuffledGeomKeys(prng, curr.Neighbors())
		for _, k := range adj {
			dir := curr.Neighbors()[k]

			// Is this adjacent hex already placed?
			neighbor, isPlaced := placed[k]
			if isPlaced {
				// Does it already expect a connection from curr`?
				if _, ok := neighbor.Connections[geom.Opposite[dir]]; ok {
					toPlace.Connections[dir] = struct{}{}
					fmt.Printf("\tneighbor %v already expects a connection\n", k)
				}
				continue
			}

			// Bail out logic .....
			segments := pathSegments(placed)
			if len(segments) > 0 && len(segments[0]) > goalComplexity {
				continue
			}

			if rollConnection(prng, len(toPlace.Connections), len(placed), goalComplexity) {
				toPlace.Connections[dir] = struct{}{}
				fmt.Printf("\tAdding %v\n", curr.Adjacent()[dir])
				toHandle = append(toHandle, curr.Adjacent()[dir])
			}
		}

		if len(toPlace.Connections) > 0 {
			fmt.Printf("\tadded with connections %v\n", keysOf(toPlace.Connections))
			// Done, chuck it into the placed set.
			placed[curr] = toPlace

		} else {
			fmt.Printf("\tAborting %v, not connected\n", curr)
		}
	}

	segments := pathSegments(placed)
	if len(segments) > 1 {
		fmt.Printf("dropping %v vestiges:", len(segments)-1)
		for i := 1; i < len(segments); i++ {
			fmt.Printf(" %v", len(segments[i]))
			for _, key := range segments[i] {
				delete(placed, key)
			}
		}
		fmt.Println()
	}
	duration := time.Since(start)
	fmt.Printf("complexity: %v of %v (%v)\n", len(placed), goalComplexity, duration)

	result := Paths{
		Algorithm: "maze-paths",
		Seed:      seed,
		Nodes:     placed,
	}
	for _, k := range shuffledGeomKeys(prng, placed) {
		if len(placed[k].Connections) == 1 {
			// found a contender
			if geom.Equal(&result.Start, &geom.Key{}) {
				result.Start = k
				continue
			} else if geom.Equal(&result.Goal, &geom.Key{}) {
				result.Goal = k
				break
			}
		}
	}

	return result, nil
}

func rollConnection(prng *rand.Rand, numConnections int, currComplexity, goalComplexity int) bool {
	completion := float64(currComplexity) / float64(goalComplexity)
	completion -= 0.8
	if completion < 0 {
		completion = 0
	}
	completion *= 5
	reduction := 1.0 - completion

	baseChance := 0.0
	switch numConnections {
	case 0:
		baseChance = 0.95
	case 1:
		baseChance = 0.2
	case 2:
		baseChance = 0.08
	case 3:
		baseChance = 0.03
	case 4:
		baseChance = 0.01
	default: // 5
		baseChance = 0.005
	}

	fmt.Printf("\treduction:%v\n", reduction)
	return prng.Float64() < baseChance //*reduction
}

func pathSegments(placed map[geom.Key]Placement) [][]geom.Key {
	// queue of keys that have not had their connections mapped yet.
	queue := map[geom.Key]struct{}{}
	for k := range placed {
		queue[k] = struct{}{}
	}

	segments := [][]geom.Key{}
	for len(queue) > 0 {
		curr := deterministicKeyFrom(queue)
		connected := []geom.Key{}
		segment := []geom.Key{}
		for {
			if _, ok := queue[curr]; ok {
				segment = append(segment, curr)
				delete(queue, curr)
				for dir := range placed[curr].Connections {
					neighbor := curr.Adjacent()[dir]
					connected = append(connected, neighbor)
				}
			}
			if len(connected) == 0 {
				break
			}
			curr, connected = connected[0], connected[1:]
		}
		segments = append(segments, segment)
	}

	sort.Slice(segments, func(i, j int) bool {
		return len(segments[i]) > len(segments[j])
	})

	return segments
}
