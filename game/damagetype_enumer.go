// Code generated by "enumer -type=DamageType -json"; DO NOT EDIT.

package game

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _DamageTypeName = "PhysicalDamageMagicalDamageFireDamage"

var _DamageTypeIndex = [...]uint8{0, 14, 27, 37}

const _DamageTypeLowerName = "physicaldamagemagicaldamagefiredamage"

func (i DamageType) String() string {
	if i < 0 || i >= DamageType(len(_DamageTypeIndex)-1) {
		return fmt.Sprintf("DamageType(%d)", i)
	}
	return _DamageTypeName[_DamageTypeIndex[i]:_DamageTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _DamageTypeNoOp() {
	var x [1]struct{}
	_ = x[PhysicalDamage-(0)]
	_ = x[MagicalDamage-(1)]
	_ = x[FireDamage-(2)]
}

var _DamageTypeValues = []DamageType{PhysicalDamage, MagicalDamage, FireDamage}

var _DamageTypeNameToValueMap = map[string]DamageType{
	_DamageTypeName[0:14]:       PhysicalDamage,
	_DamageTypeLowerName[0:14]:  PhysicalDamage,
	_DamageTypeName[14:27]:      MagicalDamage,
	_DamageTypeLowerName[14:27]: MagicalDamage,
	_DamageTypeName[27:37]:      FireDamage,
	_DamageTypeLowerName[27:37]: FireDamage,
}

var _DamageTypeNames = []string{
	_DamageTypeName[0:14],
	_DamageTypeName[14:27],
	_DamageTypeName[27:37],
}

// DamageTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DamageTypeString(s string) (DamageType, error) {
	if val, ok := _DamageTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _DamageTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to DamageType values", s)
}

// DamageTypeValues returns all values of the enum
func DamageTypeValues() []DamageType {
	return _DamageTypeValues
}

// DamageTypeStrings returns a slice of all String values of the enum
func DamageTypeStrings() []string {
	strs := make([]string, len(_DamageTypeNames))
	copy(strs, _DamageTypeNames)
	return strs
}

// IsADamageType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i DamageType) IsADamageType() bool {
	for _, v := range _DamageTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for DamageType
func (i DamageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for DamageType
func (i *DamageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("DamageType should be a string, got %s", data)
	}

	var err error
	*i, err = DamageTypeString(s)
	return err
}
