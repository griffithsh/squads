package graph

import "testing"

type node struct {
	m, n int
}

func TestSearch(t *testing.T) {
	s := NewSearcher(nil, nil, nil)

	start := node{0, 0}
	dest := node{12, 8}
	s.Search(start, dest)
}
