// Code generated by "enumer -output=./linearGradient_enumer.go -type=BlendType,StrategyTargetFilter -json"; DO NOT EDIT.

package procedural

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _BlendTypeName = "NoisySmoothSpiky"

var _BlendTypeIndex = [...]uint8{0, 5, 11, 16}

const _BlendTypeLowerName = "noisysmoothspiky"

func (i BlendType) String() string {
	if i < 0 || i >= BlendType(len(_BlendTypeIndex)-1) {
		return fmt.Sprintf("BlendType(%d)", i)
	}
	return _BlendTypeName[_BlendTypeIndex[i]:_BlendTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _BlendTypeNoOp() {
	var x [1]struct{}
	_ = x[Noisy-(0)]
	_ = x[Smooth-(1)]
	_ = x[Spiky-(2)]
}

var _BlendTypeValues = []BlendType{Noisy, Smooth, Spiky}

var _BlendTypeNameToValueMap = map[string]BlendType{
	_BlendTypeName[0:5]:        Noisy,
	_BlendTypeLowerName[0:5]:   Noisy,
	_BlendTypeName[5:11]:       Smooth,
	_BlendTypeLowerName[5:11]:  Smooth,
	_BlendTypeName[11:16]:      Spiky,
	_BlendTypeLowerName[11:16]: Spiky,
}

var _BlendTypeNames = []string{
	_BlendTypeName[0:5],
	_BlendTypeName[5:11],
	_BlendTypeName[11:16],
}

// BlendTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func BlendTypeString(s string) (BlendType, error) {
	if val, ok := _BlendTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _BlendTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to BlendType values", s)
}

// BlendTypeValues returns all values of the enum
func BlendTypeValues() []BlendType {
	return _BlendTypeValues
}

// BlendTypeStrings returns a slice of all String values of the enum
func BlendTypeStrings() []string {
	strs := make([]string, len(_BlendTypeNames))
	copy(strs, _BlendTypeNames)
	return strs
}

// IsABlendType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i BlendType) IsABlendType() bool {
	for _, v := range _BlendTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for BlendType
func (i BlendType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for BlendType
func (i *BlendType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("BlendType should be a string, got %s", data)
	}

	var err error
	*i, err = BlendTypeString(s)
	return err
}

const _StrategyTargetFilterName = "AnyTargetWidestNarrowest"

var _StrategyTargetFilterIndex = [...]uint8{0, 9, 15, 24}

const _StrategyTargetFilterLowerName = "anytargetwidestnarrowest"

func (i StrategyTargetFilter) String() string {
	if i < 0 || i >= StrategyTargetFilter(len(_StrategyTargetFilterIndex)-1) {
		return fmt.Sprintf("StrategyTargetFilter(%d)", i)
	}
	return _StrategyTargetFilterName[_StrategyTargetFilterIndex[i]:_StrategyTargetFilterIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _StrategyTargetFilterNoOp() {
	var x [1]struct{}
	_ = x[AnyTarget-(0)]
	_ = x[Widest-(1)]
	_ = x[Narrowest-(2)]
}

var _StrategyTargetFilterValues = []StrategyTargetFilter{AnyTarget, Widest, Narrowest}

var _StrategyTargetFilterNameToValueMap = map[string]StrategyTargetFilter{
	_StrategyTargetFilterName[0:9]:        AnyTarget,
	_StrategyTargetFilterLowerName[0:9]:   AnyTarget,
	_StrategyTargetFilterName[9:15]:       Widest,
	_StrategyTargetFilterLowerName[9:15]:  Widest,
	_StrategyTargetFilterName[15:24]:      Narrowest,
	_StrategyTargetFilterLowerName[15:24]: Narrowest,
}

var _StrategyTargetFilterNames = []string{
	_StrategyTargetFilterName[0:9],
	_StrategyTargetFilterName[9:15],
	_StrategyTargetFilterName[15:24],
}

// StrategyTargetFilterString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func StrategyTargetFilterString(s string) (StrategyTargetFilter, error) {
	if val, ok := _StrategyTargetFilterNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _StrategyTargetFilterNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to StrategyTargetFilter values", s)
}

// StrategyTargetFilterValues returns all values of the enum
func StrategyTargetFilterValues() []StrategyTargetFilter {
	return _StrategyTargetFilterValues
}

// StrategyTargetFilterStrings returns a slice of all String values of the enum
func StrategyTargetFilterStrings() []string {
	strs := make([]string, len(_StrategyTargetFilterNames))
	copy(strs, _StrategyTargetFilterNames)
	return strs
}

// IsAStrategyTargetFilter returns "true" if the value is listed in the enum definition. "false" otherwise
func (i StrategyTargetFilter) IsAStrategyTargetFilter() bool {
	for _, v := range _StrategyTargetFilterValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for StrategyTargetFilter
func (i StrategyTargetFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for StrategyTargetFilter
func (i *StrategyTargetFilter) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("StrategyTargetFilter should be a string, got %s", data)
	}

	var err error
	*i, err = StrategyTargetFilterString(s)
	return err
}
