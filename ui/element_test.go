package ui

import "testing"

func TestValign(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"top", "top"},
		{"bottom", "bottom"},
		{"middle", "middle"},
		{"banana", "top"},
	}

	for _, tc := range tests {
		t.Run(tc.value, func(t *testing.T) {
			am := AttributeMap{
				"valign": tc.value,
			}
			got := am.Valign()
			if got != tc.want {
				t.Fatalf("want %s, got %s", tc.want, got)
			}

		})
	}
}
