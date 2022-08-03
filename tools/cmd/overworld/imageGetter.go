package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"
)

//go:embed temporary.png
var temporary []byte

var resource image.Image

func init() {
	decoded, err := png.Decode(bytes.NewReader(temporary))
	if err != nil {
		panic(fmt.Errorf("couldn't decode the image file %v", err))
	}
	resource = decoded
}

type imageGetter struct{}

func (ig imageGetter) GetImage(string) (val image.Image, ok bool) {
	return resource, true
}
