package dynamic

import "testing"

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
