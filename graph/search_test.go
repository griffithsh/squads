package graph

import (
	"math"
	"testing"
)

type node struct {
	name string
	x, y float64
}

func TestSearch1(t *testing.T) {
	/*
		A 34,92
		B 65, 113
		C 94, 49
		D 122, 115

		  C
		 / \
		A   \
		 \   \
		  B - D

		Costs:
		AC = 5
		AB = 10
		BD = 10
		CD = 10
	*/
	A, B, C, D := node{"A", 34, 92}, node{"B", 65, 113}, node{"C", 94, 49}, node{"D", 122, 115}
	costs := func(v1, v2 Vertex) float64 {
		a, b := v1.(node), v2.(node)
		switch {
		case a.name == "A" && b.name == "C":
			fallthrough
		case a.name == "C" && b.name == "A":
			return 5
		default:
			return 10
		}
	}
	edges := func(v Vertex) []Vertex {
		n := v.(node)
		switch n.name {
		case "A":
			return []Vertex{B, C}
		case "B":
			return []Vertex{A, D}
		case "C":
			return []Vertex{A, D}
		case "D":
			return []Vertex{B, C}
		}
		return []Vertex{}
	}
	heur := func(v1 Vertex, v2 Vertex) float64 {
		// Pythagorean theorem can give us as-the-crow-flies distance.
		a := v1.(node).x - v1.(node).x
		b := v1.(node).y - v1.(node).y
		return math.Sqrt(a*a + b*b)
	}
	s := NewSearcher(costs, edges, heur)

	start := A
	dest := D
	path := s.Search(start, dest)

	if path == nil {
		t.Fatal("want path, got nil")
	}

	want := map[int]Step{
		0: {A, 0},
		1: {C, 5},
		2: {D, 15},
	}
	if len(want) != len(path) {
		t.Fatalf("want %d steps, got %d", len(want), len(path))
	}
	for i, step := range path {
		if step != want[i] {
			t.Errorf("want step %d as %s:%f, got %s:%f", i, want[i].V.(node).name, want[i].Cost, step.V.(node).name, step.Cost)
		}
	}
}

func TestSearch2(t *testing.T) {
	/*
			A 1,1
			B 2,2
			C 50,50

			   A
			  /|
			 / |
			C  |
		     \ |
			  \|
			   B

			Costs:
			AB = 100
			AC = 10
			CB = 10
	*/
	A, B, C := node{"A", 1, 1}, node{"B", 2, 2}, node{"C", 50, 50}
	costs := func(v1, v2 Vertex) float64 {
		a, b := v1.(node), v2.(node)
		switch {
		case a.name == "A" && b.name == "B":
			fallthrough
		case a.name == "B" && b.name == "A":
			return 100
		default:
			return 10
		}
	}
	edges := func(v Vertex) []Vertex {
		n := v.(node)
		switch n.name {
		case "A":
			return []Vertex{C, B}
		case "B":
			return []Vertex{A, C}
		case "C":
			return []Vertex{A, B}
		}
		return []Vertex{}
	}
	heur := func(v1 Vertex, v2 Vertex) float64 {
		// Pythagorean theorem can give us as-the-crow-flies distance.
		a := v1.(node).x - v1.(node).x
		b := v1.(node).y - v1.(node).y
		return math.Sqrt(a*a + b*b)
	}
	s := NewSearcher(costs, edges, heur)
	path := s.Search(A, B)

	if path == nil {
		t.Fatal("want path, got nil")
	}

	want := map[int]Step{
		0: {A, 0},
		1: {C, 10},
		2: {B, 20},
	}
	if len(want) != len(path) {
		t.Fatalf("want %d steps, got %d: %v", len(want), len(path), path)
	}
	for i, step := range path {
		if step != want[i] {
			t.Errorf("want step %d as %s:%f, got %s:%f", i, want[i].V.(node).name, want[i].Cost, step.V.(node).name, step.Cost)
		}
	}
}
