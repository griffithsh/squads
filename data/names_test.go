package data

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNames(t *testing.T) {
	input := `Tom:F,M,  Villanous
Jerry: ,M,,


	Mark `

	r := bytes.NewReader([]byte(input))

	got, err := parseNames(r)

	if err != nil {
		t.Fatalf("parsing test data: %v", err)
	}

	want := map[string][]string{
		"Tom":   []string{"F", "M", "Villanous"},
		"Jerry": []string{"M"},
		"Mark":  []string{},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}
