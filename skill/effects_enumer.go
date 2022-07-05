// Code generated by "enumer -type=InjuryType,Operator -json -output effects_enumer.go"; DO NOT EDIT.

package skill

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _InjuryTypeName = "BleedingInjury"

var _InjuryTypeIndex = [...]uint8{0, 14}

const _InjuryTypeLowerName = "bleedinginjury"

func (i InjuryType) String() string {
	if i < 0 || i >= InjuryType(len(_InjuryTypeIndex)-1) {
		return fmt.Sprintf("InjuryType(%d)", i)
	}
	return _InjuryTypeName[_InjuryTypeIndex[i]:_InjuryTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _InjuryTypeNoOp() {
	var x [1]struct{}
	_ = x[BleedingInjury-(0)]
}

var _InjuryTypeValues = []InjuryType{BleedingInjury}

var _InjuryTypeNameToValueMap = map[string]InjuryType{
	_InjuryTypeName[0:14]:      BleedingInjury,
	_InjuryTypeLowerName[0:14]: BleedingInjury,
}

var _InjuryTypeNames = []string{
	_InjuryTypeName[0:14],
}

// InjuryTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func InjuryTypeString(s string) (InjuryType, error) {
	if val, ok := _InjuryTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _InjuryTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to InjuryType values", s)
}

// InjuryTypeValues returns all values of the enum
func InjuryTypeValues() []InjuryType {
	return _InjuryTypeValues
}

// InjuryTypeStrings returns a slice of all String values of the enum
func InjuryTypeStrings() []string {
	strs := make([]string, len(_InjuryTypeNames))
	copy(strs, _InjuryTypeNames)
	return strs
}

// IsAInjuryType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i InjuryType) IsAInjuryType() bool {
	for _, v := range _InjuryTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for InjuryType
func (i InjuryType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for InjuryType
func (i *InjuryType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("InjuryType should be a string, got %s", data)
	}

	var err error
	*i, err = InjuryTypeString(s)
	return err
}

const _OperatorName = "AddOpMultOp"

var _OperatorIndex = [...]uint8{0, 5, 11}

const _OperatorLowerName = "addopmultop"

func (i Operator) String() string {
	if i < 0 || i >= Operator(len(_OperatorIndex)-1) {
		return fmt.Sprintf("Operator(%d)", i)
	}
	return _OperatorName[_OperatorIndex[i]:_OperatorIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _OperatorNoOp() {
	var x [1]struct{}
	_ = x[AddOp-(0)]
	_ = x[MultOp-(1)]
}

var _OperatorValues = []Operator{AddOp, MultOp}

var _OperatorNameToValueMap = map[string]Operator{
	_OperatorName[0:5]:       AddOp,
	_OperatorLowerName[0:5]:  AddOp,
	_OperatorName[5:11]:      MultOp,
	_OperatorLowerName[5:11]: MultOp,
}

var _OperatorNames = []string{
	_OperatorName[0:5],
	_OperatorName[5:11],
}

// OperatorString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func OperatorString(s string) (Operator, error) {
	if val, ok := _OperatorNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _OperatorNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Operator values", s)
}

// OperatorValues returns all values of the enum
func OperatorValues() []Operator {
	return _OperatorValues
}

// OperatorStrings returns a slice of all String values of the enum
func OperatorStrings() []string {
	strs := make([]string, len(_OperatorNames))
	copy(strs, _OperatorNames)
	return strs
}

// IsAOperator returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Operator) IsAOperator() bool {
	for _, v := range _OperatorValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Operator
func (i Operator) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Operator
func (i *Operator) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Operator should be a string, got %s", data)
	}

	var err error
	*i, err = OperatorString(s)
	return err
}
