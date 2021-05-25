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

func TestSplit(t *testing.T) {
	tests := []struct {
		name  string
		have  string
		size  TextSize
		width int
		want  []string
	}{
		{"short", "This is pretty short.", TextSizeNormal, 100, []string{"This is pretty", "short."}},
		{"letter", "Dear sirs,\nHow are you?\nSincerely", TextSizeNormal, 250, []string{"Dear sirs,", "How are you?", "Sincerely"}},
		{"medium", "This is a medium-length sentence", TextSizeNormal, 150, []string{"This is a medium-length", "sentence"}},
		{"long", "Well, this is embarassing... I never thought that this could go so wrong! But here I find myself, humiliated.", TextSizeNormal, 210, []string{"Well, this is embarassing... I never", "thought that this could go so wrong! But", "here I find myself, humiliated."}},
		{"complex", "This is the title of my manuscript.\nIt's short.", TextSizeNormal, 125, []string{"This is the title of my", "manuscript.", "It's short."}},

		{"not-sure", "First second third fourth fifth sixth seventh eighth ninth tenth eleventh twelfth.", TextSizeSmall, 188, []string{"First second third fourth fifth sixth seventh", "eighth ninth tenth eleventh twelfth."}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s-%dpx", tc.name, tc.width), func(t *testing.T) {

			text := NewText(tc.have, tc.size)

			splitLines := SplitLines(text.Lines, tc.width)

			for i, line := range splitLines {
				got := line.String()
				want := ""
				if len(tc.want) > i {
					want = tc.want[i]
				}
				if got != want {
					t.Errorf("line %d: want %q, got %q (width %dpx)", i, want, got, line.Width())
				}
			}
			for i := len(splitLines); i < len(tc.want); i++ {
				t.Errorf("line %d: want %q, got nothing", i, tc.want[i])
			}
		})
	}
}
