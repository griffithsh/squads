// Code generated by "stringer -output=./targeting_string.go -type=TargetingRule,TargetingBrush"; DO NOT EDIT.

package skill

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TargetAnywhere-0]
	_ = x[TargetAdjacent-1]
	_ = x[Untargeted-2]
}

const _TargetingRule_name = "TargetAnywhereTargetAdjacentUntargeted"

var _TargetingRule_index = [...]uint8{0, 14, 28, 38}

func (i TargetingRule) String() string {
	if i < 0 || i >= TargetingRule(len(_TargetingRule_index)-1) {
		return "TargetingRule(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TargetingRule_name[_TargetingRule_index[i]:_TargetingRule_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SingleHex-0]
	_ = x[Pathfinding-1]
	_ = x[AreaOfEffect-2]
	_ = x[None-3]
}

const _TargetingBrush_name = "SingleHexPathfindingAreaOfEffectNone"

var _TargetingBrush_index = [...]uint8{0, 9, 20, 32, 36}

func (i TargetingBrush) String() string {
	if i < 0 || i >= TargetingBrush(len(_TargetingBrush_index)-1) {
		return "TargetingBrush(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TargetingBrush_name[_TargetingBrush_index[i]:_TargetingBrush_index[i+1]]
}
