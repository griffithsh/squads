// Code generated by "stringer -type=CombatResult"; DO NOT EDIT.

package game

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Victorious-0]
	_ = x[Defeated-1]
	_ = x[Escaped-2]
}

const _CombatResult_name = "VictoriousDefeatedEscaped"

var _CombatResult_index = [...]uint8{0, 10, 18, 25}

func (i CombatResult) String() string {
	if i < 0 || i >= CombatResult(len(_CombatResult_index)-1) {
		return "CombatResult(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CombatResult_name[_CombatResult_index[i]:_CombatResult_index[i+1]]
}
