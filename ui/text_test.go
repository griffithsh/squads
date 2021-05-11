package ui

import (
	"fmt"
	"testing"
)

func TestNewTextNumWords(t *testing.T) {
	tests := []struct {
		input string
		want  int // want is the number of words we expect
	}{
		{"cat", 1},
		{"the cat", 2},
		{"sonorous,\nindefatigable", 2},
		{"", 0},
		{"Sphinx of black quartz,\njudge my vow.", 7},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			text := NewText(tc.input, TextSizeNormal)

			if len(text.Value) != tc.want {
				t.Fatalf("want %d, got %d", tc.want, len(text.Value))
			}
		})
	}
}
