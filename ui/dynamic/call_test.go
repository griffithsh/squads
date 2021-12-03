package dynamic

import (
	"fmt"
	"strings"
	"testing"
)

func TestCall(t *testing.T) {
	t.Run("WithMap", func(t *testing.T) {
		x := "original"
		f := map[string]interface{}{
			"f": func() {
				x = "modified"
			},
		}
		err := Call("f", f)
		if err != nil {
			t.Fatalf("execute error: %v", err)
		}

		if x == "original" {
			t.Error("original was not changed")
		}
	})
	t.Run("WithStruct", func(t *testing.T) {
		x := "original"
		type foo struct {
			Handler func()
		}
		f := foo{
			Handler: func() {
				x = "modified"
			},
		}
		Call("Handler", f)

		if x == "original" {
			t.Error("original was not changed")
		}
	})
	t.Run("WithPtr", func(t *testing.T) {
		x := "original"
		type foo struct {
			Handler func()
		}
		f := foo{
			Handler: func() {
				x = "modified"
			},
		}
		Call("Handler", &f)

		if x == "original" {
			t.Error("original was not changed")
		}
	})
	// should it be possible to call unexported funcs?
	t.Run("WithStructUnexported", func(t *testing.T) {
		return
		x := "original"
		type foo struct {
			handler func()
		}
		f := foo{
			handler: func() {
				x = "modified"
			},
		}
		Call("handler", f)

		if x == "original" {
			t.Error("original was not changed")
		}
	})
	// Would I need to be able to descend through structures?
	t.Run("Descendants", func(t *testing.T) {
		return
		x := "original"
		type foo struct {
			bar struct {
				Handler func()
			}
		}
		f := foo{
			bar: struct{ Handler func() }{
				Handler: func() {
					x = "modified"
				},
			},
		}
		Call("bar.Handler", f)

		if x == "original" {
			t.Error("original was not changed")
		}
	})
}

func TestRanger(t *testing.T) {
	tests := []struct {
		name     string
		slice    string
		data     interface{}
		succeeds bool
		want     string
	}{
		{
			name:  "slice-of-int",
			slice: "foo",
			data: struct{ foo []int }{
				foo: []int{1, 2, 3},
			},
			succeeds: true,
			want:     "1,2,3",
		},
		{
			name:  "slice-of-struct",
			slice: "bar",
			data: struct {
				bar []struct {
					id string
				}
			}{
				bar: []struct{ id string }{
					{id: "AK-104F7"}, {"HW-66320"},
				},
			},
			succeeds: true,
			want:     "{AK-104F7},{HW-66320}",
		},
		{
			name:  "data-is-ptr",
			slice: "foo",
			data: &struct{ foo []int }{
				foo: []int{1, 2, 3},
			},
			succeeds: true,
			want:     "1,2,3",
		},
		{
			name:  "non-slice",
			slice: "baz",
			data: struct {
				baz string
			}{
				baz: "no slices here!",
			},
		},
		{
			name:  "not-found",
			slice: "qux",
			data: struct {
				quux []int
			}{
				quux: []int{8, 8, 8},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Ranger(tc.slice, tc.data)
			if err != nil {
				fmt.Println("error:", err)
				if tc.succeeds {
					t.Fatal()
				}
			}

			var got []string
			for _, item := range result {
				got = append(got, fmt.Sprintf("%v", item))
			}

			if strings.Join(got, ",") != tc.want {
				t.Error(fmt.Sprintf("want %s, got %s", tc.want, strings.Join(got, ",")))
			}
		})
	}
}
