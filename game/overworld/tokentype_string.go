// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package overworld

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SquadToken-1]
}

const _TokenType_name = "SquadToken"

var _TokenType_index = [...]uint8{0, 10}

func (i TokenType) String() string {
	i -= 1
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
