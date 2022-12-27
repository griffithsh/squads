package procedural

import (
	"errors"
	"math/rand"
	"sort"

	"github.com/griffithsh/squads/geom"
)

// buildDoubleBlobPaths is intended to implement maps where there should be a
// single intersection point between two distinct mazes. Imagine a dark forest
// with a river running through the middle where there is only one bridge.
func buildDoubleBlobPaths(seed int64, level int) (Paths, error) {
	riverLength := 36
	lobeDensity := 10
	lobeSize := 14
	prng := rand.New(rand.NewSource(seed))

	// pick a direction - N,S,NE,SE,NW,SW
	// grab its opposite
	// from center, travel between X and Y hexes in those directions
	// from those keys, meander back towards 0,0, saving the paths
	// expand the paths by 1, randomly omitting some of the keys
	// this forms the Blockage in the center of the map.

	// from 0,0, use one of the perpendicular directions to the picked direction and opposite to move from the start by two hexes.
	// these form the seeds of the two lobes, or blobs
	// grow out in some fashion
	// assign a start from one lobe, and a goal from the other

	directions := []geom.DirectionType{geom.N, geom.S, geom.SE, geom.SW, geom.NE, geom.NW}
	shuffleSlice(prng, directions)
	sourceDir, destinationDir := directions[0], geom.Opposite[directions[0]]

	source, destination := geom.Key{}, geom.Key{}
	for i := 0; i < riverLength; i++ {
		source = source.ToDirection(sourceDir)
		destination = destination.ToDirection(destinationDir)
	}

	// zero = { M:0, N:0 }
	zero := geom.Key{}
	fromSource := meander(prng, source, zero, []geom.Key{})
	fromDestination := meander(prng, destination, zero, fromSource)

	riverSlice := append(fromSource, fromDestination...)

	river := map[geom.Key]struct{}{}
	for _, k := range riverSlice {
		river[k] = struct{}{}
		neighbors := k.ExpandBy(1, 1)
		sortKeys(neighbors)
		shuffleSlice(prng, neighbors)
		// Every hex has 6 neighbors, if we only pick four, we'll get a small
		// chance of gaps. Hopefully this makes it look more organic.
		for i := 0; i < 4; i++ {
			river[neighbors[i]] = struct{}{}
		}
	}

	perpendiculars := []geom.RelativeDirection{geom.BackLeft, geom.ForwardLeft}
	shuffleSlice(prng, perpendiculars)
	perpendicularDir := perpendiculars[0]
	backwards := geom.Actualize(sourceDir, perpendicularDir)
	forwards := geom.Opposite[backwards]

	paths := Paths{
		Nodes: map[geom.Key]Placement{
			{}: {Connections: map[geom.DirectionType]struct{}{}},
		},
		Specials: map[string][]geom.Key{
			"DividerKeys": keysOf(river),
		},
	}

	var startSeed, goalSeed geom.Key // so 0,0
	for {
		if _, ok := river[startSeed]; !ok {
			break
		}
		next1 := startSeed.ToDirection(backwards)
		paths.Connect(startSeed, next1)
		startSeed = next1
	}
	for {
		if _, ok := river[goalSeed]; !ok {
			break
		}
		next2 := goalSeed.ToDirection(forwards)
		paths.Connect(goalSeed, next2)
		goalSeed = next2
	}

	// project away from the goal
	o := goalSeed
	for i := 0; i < prng.Intn(3)+3; i++ {
		o = o.ToDirection(forwards)
	}
	sz := prng.Intn(5) + 5
	blob := o.ExpandBy(sz, sz)
	sortKeys(blob)
	shuffleSlice(prng, blob)

	contenders := []geom.Key{}
	for _, k := range blob {
		_, ok := river[k]
		if ok {
			// don't use it if it's in the river.
			continue
		}
		contenders = append(contenders, k)
		if len(contenders) >= 8 {
			// stop looking if we've got enough contenders
			break
		}
	}

	leftBank := func(k geom.Key) bool {
		x, y := geom.FlatField.Ktow(k)
		ax, ay := geom.FlatField.Ktow(source)
		bx, by := geom.FlatField.Ktow(destination)
		return ((bx-ax)*(y-ay) - (by-ay)*(x-ax)) > 0
	}

	// generate a series of random Keys, dividing them into each bank, discarding any that lie in the river
	bank1 := []geom.Key{}
	bank2 := []geom.Key{}

	options := geom.Key{}.ExpandBy(1, lobeSize)
	sortKeys(options)
	shuffleSlice(prng, options)

	if leftBank(startSeed) {
		bank1 = append(bank1, startSeed)
	} else {
		bank2 = append(bank2, startSeed)
	}
	if leftBank(goalSeed) {
		bank1 = append(bank1, goalSeed)
	} else {
		bank2 = append(bank2, goalSeed)
	}

	for {
		contender := options[0]

		// pop!
		options = options[1:]

		// Don't use start or goal seeds, they're included elsewhere.
		if contender == startSeed || contender == goalSeed {
			continue
		}

		// Don't use anything that's in the river.
		if _, inTheRiver := river[contender]; inTheRiver {
			continue
		}

		if leftBank(contender) {
			bank1 = append(bank1, contender)
		} else {
			bank2 = append(bank2, contender)
		}

		if len(options) == 0 {
			return Paths{}, errors.New("running out of options")
		}
		if len(bank1)+len(bank2) >= lobeDensity {
			break
		}
	}

	for _, contenders := range [][]geom.Key{bank1, bank2} {
		// Connect all these to their three nearest neighbors
		for _, poi := range contenders {
			neighbors := []Pair[geom.Key, float64]{}
			for _, other := range contenders {
				if poi == other {
					continue
				}
				dist := geom.FlatField.DistanceBetween(poi, other)
				neighbors = append(neighbors, Pair[geom.Key, float64]{
					key:   other,
					value: dist,
				})
			}
			sort.Slice(neighbors, func(i, j int) bool {
				return neighbors[i].value < neighbors[j].value
			})
			fails := 0
			for i, neighbor := range neighbors {
				if i == 3 {
					break
				}

				// Connect poi to neighbor.key
				path := meander(prng, poi, neighbor.key, keysOf(river))
				if len(path) == 0 {
					// retry!
					path = meander(prng, poi, neighbor.key, keysOf(river))
				}
				if len(path) == 0 {
					// fail!
					fails++
				}
				for i := 1; i < len(path); i++ {
					paths.Connect(path[i-1], path[i])
				}
			}
			if fails == len(neighbors) || fails == 3 {
				// Did not connect!
				return Paths{}, errors.New("unable to meander from this poi")
			}
		}
	}
	paths.Start = bank1[1]
	paths.Goal = bank2[1]

	return paths, nil
}

type Pair[K, V any] struct {
	key   K
	value V
}
