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
			"f": func(id string) {
				x = id
			},
		}
		err := Call("f", f, "modified")
		if err != nil {
			t.Fatalf("execute error: %v", err)
		}

		if x == "original" {
			t.Error("original was not changed")
		}
		if x != "modified" {
			t.Errorf("want %q, got %q", "modified", x)
		}
	})
	t.Run("WithStruct", func(t *testing.T) {
		x := "original"
		type foo struct {
			Handler func(string)
		}
		f := foo{
			Handler: func(string) {
				x = "modified"
			},
		}
		Call("Handler", f, "")

		if x == "original" {
			t.Error("original was not changed")
		}
	})
	t.Run("WithPtr", func(t *testing.T) {
		x := "original"
		type foo struct {
			Handler func(string)
		}
		f := foo{
			Handler: func(string) {
				x = "modified"
			},
		}
		Call("Handler", &f, "")

		if x == "original" {
			t.Error("original was not changed")
		}
	})
	// Would I need to be able to descend through structures?
	// t.Run("Descendants", func(t *testing.T) {
	// 	x := "original"
	// 	type foo struct {
	// 		bar struct {
	// 			Handler func()
	// 		}
	// 	}
	// 	f := foo{
	// 		bar: struct{ Handler func() }{
	// 			Handler: func() {
	// 				x = "modified"
	// 			},
	// 		},
	// 	}
	// 	Call("bar.Handler", f)

	// 	if x == "original" {
	// 		t.Error("original was not changed")
	// 	}
	// })
}

func TestRanger(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		data     interface{}
		succeeds bool
		want     string
	}{
		{
			name:  "slice-of-int",
			field: "Foo",
			data: struct{ Foo []int }{
				Foo: []int{1, 2, 3},
			},
			succeeds: true,
			want:     "1,2,3",
		},
		{
			name:  "slice-of-struct",
			field: "Bar",
			data: struct {
				Bar []struct {
					id string
				}
			}{
				Bar: []struct{ id string }{
					{id: "AK-104F7"}, {"HW-66320"},
				},
			},
			succeeds: true,
			want:     "{AK-104F7},{HW-66320}",
		},
		{
			name:  "data-is-ptr",
			field: "Foo",
			data: &struct{ Foo []int }{
				Foo: []int{1, 2, 3},
			},
			succeeds: true,
			want:     "1,2,3",
		},
		{
			name:  "non-slice",
			field: "Baz",
			data: struct {
				Baz string
			}{
				Baz: "no slices here!",
			},
		},
		{
			name:  "not-found",
			field: "Qux",
			data: struct {
				Quux []int
			}{
				Quux: []int{8, 8, 8},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic: %v", r)
				}
			}()
			result, err := Ranger(tc.field, tc.data)
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
