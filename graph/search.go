package graph

type Searcher struct {
	cost CostFunc
	adj  EdgeFunc
	heur Heuristic
}

func NewSearcher(cost CostFunc, edge EdgeFunc, guess Heuristic) *Searcher {
	return &Searcher{
		cost: cost,
		adj:  edge,
		heur: guess,
	}
}

// Search finds the path between from and to. It returns nil when no path is available.
func (s *Searcher) Search(from, to Vertex) []Vertex {
	return nil
}
