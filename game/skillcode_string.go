// Code generated by "stringer -type=SkillCode"; DO NOT EDIT.

package game

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BasicMovement-0]
	_ = x[BasicAttack-1]
}

const _SkillCode_name = "BasicMovementBasicAttack"

var _SkillCode_index = [...]uint8{0, 13, 24}

func (i SkillCode) String() string {
	if i < 0 || i >= SkillCode(len(_SkillCode_index)-1) {
		return "SkillCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SkillCode_name[_SkillCode_index[i]:_SkillCode_index[i+1]]
}
