// Code generated by "stringer -type=State"; DO NOT EDIT.

package combat

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AwaitingInputState-0]
	_ = x[ExecutingState-1]
	_ = x[ThinkingState-2]
	_ = x[PreparingState-3]
}

const _State_name = "AwaitingInputStateExecutingStateThinkingStatePreparingState"

var _State_index = [...]uint8{0, 18, 32, 45, 59}

func (i State) String() string {
	if i < 0 || i >= State(len(_State_index)-1) {
		return "State(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}
