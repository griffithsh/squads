package hbg

import (
	"testing"
	"time"
)

func TestMarshalBaseTile(t *testing.T) {

	if err := testMarshaling(BaseTile{
		Code:    "WATER",
		Texture: "water.png",
		Variations: Options{
			{
				Chance: 10,
				Frames: []Frame{
					{Duration: time.Nanosecond * 12, Left: 128, Top: 0},
				},
				ExtraHeight: 0,
				ExtraDepth:  0,
			},
			{
				Chance: 1,
				Frames: []Frame{
					{Duration: time.Nanosecond * 250, Left: 42, Top: 12},
				},
				ExtraHeight: 12,
				ExtraDepth:  4,
			},
		},
	}); err != nil {
		t.Errorf("%v", err)
	}
}
