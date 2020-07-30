package graph

import (
	"math"
)

// Searcher provides a way to find a path through a graph of Vertices with
// edges.
type Searcher struct {
	cost CostFunc
	adj  EdgeFunc
	heur Heuristic
}

// NewSearcher create a new Searcher.
func NewSearcher(cost CostFunc, edge EdgeFunc, guess Heuristic) *Searcher {
	return &Searcher{
		cost: cost,
		adj:  edge,
		heur: guess,
	}
}

// Search finds the path between two Vertices. It returns nil when no path is
// available.
func (s *Searcher) Search(start, goal Vertex) []Step {
	closed := map[Vertex]interface{}{}
	open := map[Vertex]interface{}{
		start: struct{}{},
	}
	origins := map[Vertex]Vertex{}
	costs := map[Vertex]float64{
		start: 0,
	}
	guesses := map[Vertex]float64{
		start: s.heur(start, goal),
	}

	for len(open) > 0 {
		low := math.MaxFloat64
		var current Vertex
		for k := range open {
			if guesses[k] < low {
				current = k
				low = guesses[k]
			}
		}
		var unassigned Vertex
		if current == unassigned {
			// We've gone through the entire open list without finding a path.
			return nil
		}
		if current == goal {
			m := map[Vertex]Vertex{}
			for k, v := range origins {
				m[k] = v
			}
			return reconstruct(m, costs, goal)
		}

		delete(open, current)
		closed[current] = struct{}{}

		for _, n := range s.adj(current) {
			if _, ok := closed[n]; ok {
				continue
			}

			tentative := costs[current] + s.cost(current, n)

			if _, ok := open[n]; !ok {
				open[n] = struct{}{}
			} else if tentative >= costs[n] {
				continue
			}

			origins[n] = current
			costs[n] = tentative
			guesses[n] = costs[n] + s.heur(n, goal)
		}
	}
	return nil
}

// Step in a path Search result.
type Step struct {
	V    Vertex
	Cost float64
}

// reconstruct the path to the start by following backwards from goal through
// the origins to a Vertex that has no origin (i.e, the start).
func reconstruct(origins map[Vertex]Vertex, costs map[Vertex]float64, goal Vertex) []Step {
	result := []Step{
		{
			V:    goal,
			Cost: costs[goal],
		},
	}

	origin, ok := origins[goal]
	for ok {
		result = append(result, Step{
			V:    origin,
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

	return result
}
