// Code generated by "enumer -type=SelectableType -json"; DO NOT EDIT.

package targeting

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _SelectableTypeName = "SelectAnywhereSelectWithinUntargeted"

var _SelectableTypeIndex = [...]uint8{0, 14, 26, 36}

const _SelectableTypeLowerName = "selectanywhereselectwithinuntargeted"

func (i SelectableType) String() string {
	if i < 0 || i >= SelectableType(len(_SelectableTypeIndex)-1) {
		return fmt.Sprintf("SelectableType(%d)", i)
	}
	return _SelectableTypeName[_SelectableTypeIndex[i]:_SelectableTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _SelectableTypeNoOp() {
	var x [1]struct{}
	_ = x[SelectAnywhere-(0)]
	_ = x[SelectWithin-(1)]
	_ = x[Untargeted-(2)]
}

var _SelectableTypeValues = []SelectableType{SelectAnywhere, SelectWithin, Untargeted}

var _SelectableTypeNameToValueMap = map[string]SelectableType{
	_SelectableTypeName[0:14]:       SelectAnywhere,
	_SelectableTypeLowerName[0:14]:  SelectAnywhere,
	_SelectableTypeName[14:26]:      SelectWithin,
	_SelectableTypeLowerName[14:26]: SelectWithin,
	_SelectableTypeName[26:36]:      Untargeted,
	_SelectableTypeLowerName[26:36]: Untargeted,
}

var _SelectableTypeNames = []string{
	_SelectableTypeName[0:14],
	_SelectableTypeName[14:26],
	_SelectableTypeName[26:36],
}

// SelectableTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SelectableTypeString(s string) (SelectableType, error) {
	if val, ok := _SelectableTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _SelectableTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to SelectableType values", s)
}

// SelectableTypeValues returns all values of the enum
func SelectableTypeValues() []SelectableType {
	return _SelectableTypeValues
}

// SelectableTypeStrings returns a slice of all String values of the enum
func SelectableTypeStrings() []string {
	strs := make([]string, len(_SelectableTypeNames))
	copy(strs, _SelectableTypeNames)
	return strs
}

// IsASelectableType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i SelectableType) IsASelectableType() bool {
	for _, v := range _SelectableTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for SelectableType
func (i SelectableType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for SelectableType
func (i *SelectableType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("SelectableType should be a string, got %s", data)
	}

	var err error
	*i, err = SelectableTypeString(s)
	return err
}