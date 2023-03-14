package hbg

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	test := `{
		"chance":4,
		"frames":[{
			"duration": 250,
			"left":3,
			"top":8
		}
		],
		"extraHeight": 7,
		"extraDepth":18
	}`

	var got Option

	err := json.Unmarshal([]byte(test), &got)
	if err != nil {
		t.Fatalf("unmarshal test json: %v", err)
	}

	want := Option{
		Chance: 4,
		Frames: []Frame{
			{
				Duration: time.Millisecond * 250,
				Left:     3,
				Top:      8,
			},
		},
		ExtraHeight: 7,
		ExtraDepth:  18,
	}
	if got.Chance != want.Chance {
		t.Errorf("Chance: want %d, got %d", want.Chance, got.Chance)
	}

	if len(got.Frames) != len(want.Frames) {
		t.Errorf("len of Frames: want %d, got %d", len(want.Frames), len(got.Frames))
	}
	if got.Frames[0].Duration != want.Frames[0].Duration {
		t.Errorf("Frame Duration: want %d, got %d", want.Frames[0].Duration*time.Millisecond, got.Frames[0].Duration*time.Millisecond)
	}
	if got.Frames[0].Left != want.Frames[0].Left {
		t.Errorf("Frame Left: want %d, got %d", want.Frames[0].Left, got.Frames[0].Left)
	}
	if got.Frames[0].Top != want.Frames[0].Top {
		t.Errorf("Frame Top: want %d, got %d", want.Frames[0].Top, got.Frames[0].Top)
	}

	if got.ExtraHeight != want.ExtraHeight {
		t.Errorf("ExtraHeight: want %d, got %d", want.ExtraHeight, got.ExtraHeight)
	}

	if got.ExtraDepth != want.ExtraDepth {
		t.Errorf("ExtraDepth: want %d, got %d", want.ExtraDepth, got.ExtraDepth)
	}
}
