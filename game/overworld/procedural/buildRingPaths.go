package procedural

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

func meander(prng *rand.Rand, start, goal geom.Key, avoid []geom.Key) []geom.Key {
	f := geom.NewField(3, 5, 12)
	current := start
	banned := map[geom.Key]struct{}{}
	for _, k := range avoid {
		banned[k] = struct{}{}
	}

	// Yo if we hit a length longer than this something is real wrong.
	max := start.HexesFrom(goal) * 5

	result := []geom.Key{start}
	for {
		n := current.Neighbors()

		// Check for goal in site.
		if _, ok := n[goal]; ok {
			// We found the goal!
			return append(result, goal)
		}

		if len(result) > max {
			fmt.Printf("\tgave up looking for path from %v to %v after %d steps %v\n", start, goal, max, result)
			return []geom.Key{}
		}

		// Remove any banned keys from neighbors.
		for k := range n {
			if _, ok := banned[k]; ok {
				delete(n, k)
			}
		}

		if len(n) == 0 {
			fmt.Printf("\tmeander from %v to %v: no available neighbors for %v\n", start, goal, current)
			return []geom.Key{}
		}

		// Score neighbors.
		type score struct {
			k   geom.Key
			val float64
		}
		scored := []score{}
		// distance from current
		currentDistance := f.DistanceBetween(current, goal)

		for _, k := range shuffledGeomKeys(prng, n) {
			nd := f.DistanceBetween(k, goal)
			s := score{k: k, val: nd - currentDistance}
			s.val *= -1
			if s.val < 0 {
				s.val = 0.1
			}
			scored = append(scored, s)
		}

		// roll for next
		sum := 0.0
		for _, s := range scored {
			sum += s.val
		}
		roll := prng.Float64() * sum
		consumed := 0.0
		for _, s := range scored {
			if s.val+consumed >= roll {
				// Found it.
				// append next to result and add it to banned
				result = append(result, s.k)
				banned[s.k] = struct{}{}
				// assign new current
				current = s.k
				break
			}
			consumed += s.val
		}
	}
}

// buildRingPaths might implement lakes, islands with impassable interiors, or atolls.
func buildRingPaths(seed int64, level int) (Paths, error) {
	// drive these by level in some way?
	rotations := 10
	minRing := 7
	maxRing := 7

	// pick N points on a ring, connect each to the next in turn, connecting the
	// last to the first.

	fmt.Printf("buildRingPaths(%d)\n", seed)
	prng := rand.New(rand.NewSource(seed))

	f := geom.NewField(3, 5, 12)

	sin, cos := math.Sincos(math.Pi * 2 / float64(rotations))

	result := Paths{
		Nodes: map[geom.Key]Placement{},
	}

	contenders := geom.Key{}.ExpandBy(minRing, maxRing)
	sortKeys(contenders)
	i := DeterministicIndexOf(prng, contenders)
	current := contenders[i]

	result.Start = current
	goalIndex := rotations/2 + prng.Intn((rotations/3)*2) - (rotations / 3)

	toConnect := []struct{ a, b geom.Key }{}
	for i = 0; i < rotations-1; i++ {
		x, y := f.Ktow(current)
		x2, y2 := x*cos-y*sin, x*sin+y*cos
		next := f.Wtok(x2, y2)

		// Blur next a bit.
		next = shuffledGeomKeys(prng, next.Neighbors())[0]

		toConnect = append(toConnect, struct{ a, b geom.Key }{a: current, b: next})

		if goalIndex == i {
			result.Goal = next
		}

		current = next
	}
	if geom.Equal(&result.Goal, &geom.Key{}) {
		fmt.Printf("\tdid not assign goal from index %d\n", goalIndex)
	}
	toConnect = append(toConnect, struct{ a, b geom.Key }{
		// Link last to first.
		a: toConnect[len(toConnect)-1].b,
		b: toConnect[0].a,
	})

	for _, c := range toConnect {
		banned := keysOf(result.Nodes)
		if geom.Equal(&c.a, &c.b) {
			// nothing to connect!
			fmt.Printf("\tnot connecting %v to itself\n", c.a)
			continue
		}

		aPath := []geom.Key{}
		for i := 2; i >= 0; i-- {
			aPath = meander(prng, c.a, c.b, banned)
			if len(aPath) > 0 {
				break
			}
			fmt.Printf("\tmeander from %v to %v failed\n", c.a, c.b)
			if i == 0 {
				return Paths{}, fmt.Errorf("buildRingPaths: could not meander from %v to %v: exhausted retries", c.a, c.b)
			}
		}
		if len(aPath) == 0 {
			return Paths{}, fmt.Errorf("buildRingPaths: could not meander from %v to %v: exhausted retries", c.a, c.b)
		}

		fmt.Printf("\tmeandering from %v to %v: %v\n", c.a, c.b, aPath)
		step := aPath[0]
		for i := 1; i < len(aPath); i++ {
			result.Connect(step, aPath[i])
			step = aPath[i]
		}

	}

	// Right, ring done, let's add tendrils for interest
	for _, c := range toConnect {
		avoid := keysOf(result.Nodes)
		banned := map[geom.Key]struct{}{}
		for _, k := range avoid {
			banned[k] = struct{}{}
		}

		branch := c.a
		for i := 0; i < 4; i++ {
			type score struct {
				k   geom.Key
				val float64
			}
			scored := []score{}
			// distance from current
			currentDistance := f.DistanceBetween(branch, geom.Key{})
			n := branch.Neighbors()
			for k := range n {
				if _, ok := banned[k]; ok {
					delete(n, k)
				}
			}

			if len(n) == 0 {
				break
			}

			for _, k := range shuffledGeomKeys(prng, n) {
				nd := f.DistanceBetween(k, geom.Key{})
				s := score{k: k, val: math.Abs(nd - currentDistance)}
				scored = append(scored, s)
			}

			// roll for next
			sum := 0.0
			for _, s := range scored {
				sum += s.val
			}
			roll := prng.Float64() * sum
			consumed := 0.0
			for _, s := range scored {
				if s.val+consumed >= roll {
					// Found it.
					result.Connect(branch, s.k)
					banned[s.k] = struct{}{}
					branch = s.k
					break
				}
				consumed += s.val
			}
		}
	}

	return result, nil
}
