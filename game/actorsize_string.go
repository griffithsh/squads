// Code generated by "stringer -type=ActorSize"; DO NOT EDIT.

package game

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SMALL-0]
	_ = x[MEDIUM-1]
	_ = x[LARGE-2]
}

const _ActorSize_name = "SMALLMEDIUMLARGE"

var _ActorSize_index = [...]uint8{0, 5, 11, 16}

func (i ActorSize) String() string {
	if i < 0 || i >= ActorSize(len(_ActorSize_index)-1) {
		return "ActorSize(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ActorSize_name[_ActorSize_index[i]:_ActorSize_index[i+1]]
}