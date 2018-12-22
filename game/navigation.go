package game

import (
	"fmt"
	"math"
)

func reconstruct(prevs map[*Hex]*Hex, current *Hex) ([]*Hex, error) {
	result := []*Hex{current}
	n, ok := prevs[current]
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
// Hexes, ignoring obstacles.
func heuristic(a, b *Hex) float64 {
	// pythagorean theorum, minus the sqrt.
	return math.Pow(a.X()-b.X(), 2) + math.Pow(a.Y()-b.Y(), 2)
}

// Navigate a path from start to the goal, avoiding Impassable Hexes.
func Navigate(start, goal *Hex) ([]*Hex, error) {
	oneStep := heuristic(&Hex{M: 0, N: 0}, &Hex{M: 0, N: 1})

	closed := map[*Hex]interface{}{}
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
			return reconstruct(cameFrom, current)
		}
		delete(open, current)
		closed[current] = struct{}{}

		for _, n := range current.Neighbors() {
			if _, ok := closed[n]; ok {
				continue
			}

			// TODO(griffithsh):
			// Here we probably need to continue if this hex is "occupied".

			tentative := costs[current] + oneStep

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
