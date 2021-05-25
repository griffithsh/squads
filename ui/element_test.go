package ui

import "testing"

func TestNormaliseTwelfths(t *testing.T) {
	tests := []struct {
		scenario        string
		elementTwelfths []string
		wantTwelfths    []int
		wantOffsets     []int
	}{
		{
			scenario:        "simple-happy",
			elementTwelfths: []string{"6", "6"},
			wantTwelfths:    []int{6, 6},
			wantOffsets:     []int{0, 6},
		},
		{
			scenario:        "infers-everything",
			elementTwelfths: []string{"", ""},
			wantTwelfths:    []int{6, 6},
			wantOffsets:     []int{0, 6},
		},
		{
			scenario:        "multiples",
			elementTwelfths: []string{"1", "5", "5", "1"},
			wantTwelfths:    []int{1, 5, 5, 1},
			wantOffsets:     []int{0, 1, 6, 11},
		},
		{
			scenario:        "divides-end-remainder",
			elementTwelfths: []string{"3", "6", ""},
			wantTwelfths:    []int{3, 6, 3},
			wantOffsets:     []int{0, 3, 9},
		},
		{
			scenario:        "divides-start-remainder",
			elementTwelfths: []string{"", "6", "3"},
			wantTwelfths:    []int{3, 6, 3},
			wantOffsets:     []int{0, 3, 9},
		},
		{
			scenario:        "infers-multiple",
			elementTwelfths: []string{"5", "", ""},
			wantTwelfths:    []int{5, 4, 3},
			wantOffsets:     []int{0, 5, 9},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			elements := []*Element{}
			for _, attr := range tc.elementTwelfths {
				elements = append(elements, &Element{
					Attributes: AttributeMap{
						"twelfths": attr,
					},
				})
			}

			normaliseTwelfths(elements)

			for i, el := range elements {
				if tc.wantTwelfths[i] != el.Attributes.Twelfths() {
					t.Errorf("for %d, want twelfths %d, but got %d", i, tc.wantTwelfths[i], el.Attributes.Twelfths())
				}
				if tc.wantOffsets[i] != el.Attributes.TwelfthsOffset() {
					t.Errorf("for %d, want twelfths-offset %d, but got %d", i, tc.wantOffsets[i], el.Attributes.TwelfthsOffset())
				}
			}
		})
	}
}
