package pixelutil

import (
	"image"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
)

// LoadPicture is based on the loadPicture function from faiface/pixel's
// tutorials. See https://github.com/faiface/pixel/wiki
func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
