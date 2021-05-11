package ui

import (
	"fmt"
	"testing"
)

func TestNewTextNumWords(t *testing.T) {
	tests := []struct {
		input string
		want  []int // want is the number of words we expect
	}{
		{"cat", []int{1}},
		{"the cat", []int{2}},
		{"sonorous,\nindefatigable", []int{1, 1}},
		{"", []int{0}},
		{"Sphinx of black quartz,\njudge my vow.", []int{4, 3}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			text := NewText(tc.input, TextSizeNormal)

			for i, line := range text.Lines {
				if tc.want[i] != len(line) {
					t.Fatalf("for %q, want %d, got %d", tc.input, tc.want[i], len(line))
				}
			}
		})
	}
}
