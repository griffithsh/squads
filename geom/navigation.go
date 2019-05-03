package geom

import (
	"errors"
	"fmt"
	"math"
)

// ContextualObstacle captures how much of an obstacle this is to the navigator.
// A bird can fly right over a tree, a snake is not impeded by a swamp. A horse
// runs fastest when the ground is level and clear. The Cost multiplies the
// normal traversal time. A Cost of 2 implies that taking this path is twice as
// long as it normally would be. A cost of Infinity marks something that is completely impassable.
type ContextualObstacle struct {
	M, N int

	Cost float64
}

type Positioned interface {
	X() float64
	Y() float64
}

func reconstruct(prevs map[Positioned]Positioned, goal Positioned) ([]Positioned, error) {
	result := []Positioned{goal}
	n, ok := prevs[goal]
	for ok {
		result = append(result, n)

		// Next!
		n, ok = prevs[n]
	}
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}
	return result, nil
}

// heuristic determines the comparitive "as the crow flies" distance between two
// Positioned things, ignoring obstacles.
func heuristic(a, b Positioned) float64 {
	// pythagorean theorum, minus the sqrt.
	return math.Pow(a.X()-b.X(), 2) + math.Pow(a.Y()-b.Y(), 2)
}

// Navigate a path from start to the goal, avoiding Impassable Hexes.
func Navigate(start, goal *Hex, obstacles []ContextualObstacle) ([]Positioned, error) {
	if start == nil {
		return nil, errors.New("no start")
	}
	if goal == nil {
		return nil, errors.New("no goal")
	}
	oneStep := heuristic(&Hex{M: 0, N: 0}, &Hex{M: 0, N: 1})

	closed := map[Key]interface{}{}
	open := map[*Hex]interface{}{
		start: struct{}{},
	}
	cameFrom := map[*Hex]*Hex{}
	costs := map[*Hex]float64{
		start: 0,
	}
	guesses := map[*Hex]float64{
		start: heuristic(start, goal),
	}

	for len(open) > 0 {
		var current *Hex
		low := math.MaxFloat64
		for k := range open {
			if guesses[k] < low {
				current = k
				low = guesses[k]
			}
		}
		if current == goal {
			m := map[Positioned]Positioned{}
			for k, v := range cameFrom {
				m[k] = v
			}
			return reconstruct(m, goal)
		}

		if current == nil {
			break
		}

		delete(open, current)
		closed[Key{M: current.M, N: current.N}] = struct{}{}

		for _, n := range current.Neighbors() {
			if _, ok := closed[Key{M: n.M, N: n.N}]; ok {
				continue
			}

			tentative := oneStep

			// The cost of passing through this hex might be affected by any
			// obstacles occupying the Hex.
			for _, o := range obstacles {
				if o.M == n.M && o.N == n.N {
					if o.Cost == math.Inf(0) {
						tentative = math.MaxFloat64
					} else {
						tentative *= o.Cost
					}
					break
				}
			}
			tentative += costs[current]

			if _, ok := open[n]; !ok {
				open[n] = struct{}{}
			} else if tentative >= costs[n] {
				continue
			}

			cameFrom[n] = current
			costs[n] = tentative
			guesses[n] = costs[n] + heuristic(n, goal)
		}
	}
	return nil, fmt.Errorf("no path available from %d,%d to %d,%d", start.M, start.N, goal.M, goal.N)
}
