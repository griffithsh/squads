package hbg

import (
	"encoding/json"
	"time"
)

type Frame struct {
	Duration time.Duration `json:"duration"`
	Left     int           `json:"left"`
	Top      int           `json:"top"`
}

func (f Frame) MarshalJSON() ([]byte, error) {
	dummy := struct {
		Duration int `json:"duration"`
		Left     int `json:"left"`
		Top      int `json:"top"`
	}{
		Duration: int(f.Duration / time.Millisecond),
		Left:     f.Left,
		Top:      f.Top,
	}

	return json.Marshal(dummy)
}

func (f *Frame) UnmarshalJSON(data []byte) error {
	var dummy struct {
		Duration int `json:"duration"`
		Left     int `json:"left"`
		Top      int `json:"top"`
	}

	if err := json.Unmarshal(data, &dummy); err != nil {
		return err
	}

	f.Duration = time.Duration(dummy.Duration) * time.Millisecond
	f.Left = dummy.Left
	f.Top = dummy.Top
	return nil
}

type Option struct {
	Chance      int     `json:"chance"`
	Frames      []Frame `json:"frames"`
	ExtraHeight int     `json:"extraHeight"`
	ExtraDepth  int     `json:"extraDepth"`
}

type Options []Option

// Intner is anything that provides an Intn method -- i.e. a *rand.Rand.
type Intner interface {
	Intn(int) int
}

// Roll the Options, taking the weighting of the relative Chances to pick one
// based on the randomized input from the passed intn function. Pass rand.Intn
// to quickly get non-reproducible, pseudo-random results.
// Will return nil if there are no Options. Will skip the call to the passed
// random number source if the collection of Options has less than 2 members.
func (o Options) Roll(intn func(int) int) *Option {
	c := len(o)
	if c == 0 {
		return nil
	} else if c == 1 {
		return &o[0]
	}

	sum := o[0].Chance
	for i := 1; i < len(o); i++ {
		sum += o[i].Chance
	}
	roll := intn(sum)
	running := 0
	for _, opt := range o {
		if roll < opt.Chance+running {
			return &opt
		}
		running += opt.Chance
	}

	return &o[0]
}
