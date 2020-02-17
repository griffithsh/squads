package geom

import (
	"errors"
	"fmt"
	"math"
)

type NavigateStep struct {
	M, N int
	Cost float64 // Cost to travel through all previous steps to this one.
}

// Navigate from start to goal.
// existsFunc should provide whether a Hex at M,N exists or not. costFunc should
// provide the cost multiplier for moving into M,N. It should return Infinity
// if M,N is not passable. The steps returned include the start and goal.
func Navigate(start, goal Key, existsFunc func(Key) bool, costFunc func(Key) float64) ([]NavigateStep, error) {
	oneStep := 10

	// heuristic should return a guess as to how far away the two Hex coordinates are.
	heuristic := func(ma, na, mb, nb int) float64 {
		a := Hex{
			M: ma,
			N: na,
		}
		b := Hex{
			M: mb,
			N: nb,
		}
		x := math.Abs(math.Abs(a.X()) - math.Abs(b.X()))
		y := math.Abs(math.Abs(a.Y()) - math.Abs(b.Y()))

		// We can use the pythagorean theorum without the sqrt here, because we
		// only need to use the output of this function to compare against
		// other outputs of this function.
		return math.Pow(x, 2) + math.Pow(y, 2)
	}

	reconstruct := func(origins map[Key]Key, costs map[Key]float64, goal Key) ([]NavigateStep, error) {
		// Follow back from goal through the origins to a Key that has no origin (i.e, the start).
		result := []NavigateStep{
			{
				M:    goal.M,
				N:    goal.N,
				Cost: costs[goal],
			},
		}

		origin, ok := origins[goal]
		for ok {
			result = append(result, NavigateStep{
				M:    origin.M,
				N:    origin.N,
				Cost: costs[origin],
			})

			// Continue by looking for the origin of the current origin.
			origin, ok = origins[origin]
		}

		// Reverse the result.
		// https://github.com/golang/go/wiki/SliceTricks#reversing
		for i := len(result)/2 - 1; i >= 0; i-- {
			opp := len(result) - 1 - i
			result[i], result[opp] = result[opp], result[i]
		}

		return result, nil
	}

	if !existsFunc(start) {
		return nil, fmt.Errorf("no start (%d,%d)", start.M, start.N)
	}
	if !existsFunc(goal) {
		return nil, errors.New("no goal")
	}

	closed := map[Key]interface{}{}
	open := map[Key]interface{}{
		start: struct{}{},
	}
	origins := map[Key]Key{}
	costs := map[Key]float64{
		start: 0,
	}
	guesses := map[Key]float64{
		start: heuristic(start.M, start.N, goal.M, goal.N),
	}

	for len(open) > 0 {
		var current Key // NB, this used to be a pointer to Hex, but is now a value...
		low := math.MaxFloat64
		for k := range open {
			if guesses[k] < low {
				current = k
				low = guesses[k]
			}
		}
		if current == goal {
			m := map[Key]Key{}
			for k, v := range origins {
				m[k] = v
			}
			return reconstruct(m, costs, goal)
		}

		delete(open, current)
		closed[Key{M: current.M, N: current.N}] = struct{}{}

		for _, n := range neighbors(current.M, current.N) {
			if !existsFunc(n) {
				continue
			}
			if _, ok := closed[n]; ok {
				continue
			}

			tentative := float64(oneStep)

			multiplier := costFunc(n)
			if math.IsInf(multiplier, 0) {
				// I think this should work, because we will avoid adding this hex to open.
				continue
			}
			tentative = tentative * multiplier

			tentative += costs[current]

			if _, ok := open[n]; !ok {
				open[n] = struct{}{}
			} else if tentative >= costs[n] {
				continue
			}

			origins[n] = current
			costs[n] = tentative
			guesses[n] = costs[n] + heuristic(n.M, n.N, goal.M, goal.N)
		}
	}
	return nil, fmt.Errorf("no path available from %d,%d to %d,%d", start.M, start.N, goal.M, goal.N)
}
