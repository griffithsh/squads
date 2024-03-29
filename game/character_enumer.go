// Code generated by "enumer -output=./character_enumer.go -type=CharacterSex,CharacterPerformance -json"; DO NOT EDIT.

package game

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _CharacterSexName = "MaleFemale"

var _CharacterSexIndex = [...]uint8{0, 4, 10}

const _CharacterSexLowerName = "malefemale"

func (i CharacterSex) String() string {
	if i < 0 || i >= CharacterSex(len(_CharacterSexIndex)-1) {
		return fmt.Sprintf("CharacterSex(%d)", i)
	}
	return _CharacterSexName[_CharacterSexIndex[i]:_CharacterSexIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _CharacterSexNoOp() {
	var x [1]struct{}
	_ = x[Male-(0)]
	_ = x[Female-(1)]
}

var _CharacterSexValues = []CharacterSex{Male, Female}

var _CharacterSexNameToValueMap = map[string]CharacterSex{
	_CharacterSexName[0:4]:       Male,
	_CharacterSexLowerName[0:4]:  Male,
	_CharacterSexName[4:10]:      Female,
	_CharacterSexLowerName[4:10]: Female,
}

var _CharacterSexNames = []string{
	_CharacterSexName[0:4],
	_CharacterSexName[4:10],
}

// CharacterSexString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func CharacterSexString(s string) (CharacterSex, error) {
	if val, ok := _CharacterSexNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _CharacterSexNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to CharacterSex values", s)
}

// CharacterSexValues returns all values of the enum
func CharacterSexValues() []CharacterSex {
	return _CharacterSexValues
}

// CharacterSexStrings returns a slice of all String values of the enum
func CharacterSexStrings() []string {
	strs := make([]string, len(_CharacterSexNames))
	copy(strs, _CharacterSexNames)
	return strs
}

// IsACharacterSex returns "true" if the value is listed in the enum definition. "false" otherwise
func (i CharacterSex) IsACharacterSex() bool {
	for _, v := range _CharacterSexValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for CharacterSex
func (i CharacterSex) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for CharacterSex
func (i *CharacterSex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("CharacterSex should be a string, got %s", data)
	}

	var err error
	*i, err = CharacterSexString(s)
	return err
}

const _CharacterPerformanceName = "PerformIdlePerformMovePerformSkill1PerformSkill2PerformSkill3PerformHurtPerformDyingPerformVictory"

var _CharacterPerformanceIndex = [...]uint8{0, 11, 22, 35, 48, 61, 72, 84, 98}

const _CharacterPerformanceLowerName = "performidleperformmoveperformskill1performskill2performskill3performhurtperformdyingperformvictory"

func (i CharacterPerformance) String() string {
	if i < 0 || i >= CharacterPerformance(len(_CharacterPerformanceIndex)-1) {
		return fmt.Sprintf("CharacterPerformance(%d)", i)
	}
	return _CharacterPerformanceName[_CharacterPerformanceIndex[i]:_CharacterPerformanceIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _CharacterPerformanceNoOp() {
	var x [1]struct{}
	_ = x[PerformIdle-(0)]
	_ = x[PerformMove-(1)]
	_ = x[PerformSkill1-(2)]
	_ = x[PerformSkill2-(3)]
	_ = x[PerformSkill3-(4)]
	_ = x[PerformHurt-(5)]
	_ = x[PerformDying-(6)]
	_ = x[PerformVictory-(7)]
}

var _CharacterPerformanceValues = []CharacterPerformance{PerformIdle, PerformMove, PerformSkill1, PerformSkill2, PerformSkill3, PerformHurt, PerformDying, PerformVictory}

var _CharacterPerformanceNameToValueMap = map[string]CharacterPerformance{
	_CharacterPerformanceName[0:11]:       PerformIdle,
	_CharacterPerformanceLowerName[0:11]:  PerformIdle,
	_CharacterPerformanceName[11:22]:      PerformMove,
	_CharacterPerformanceLowerName[11:22]: PerformMove,
	_CharacterPerformanceName[22:35]:      PerformSkill1,
	_CharacterPerformanceLowerName[22:35]: PerformSkill1,
	_CharacterPerformanceName[35:48]:      PerformSkill2,
	_CharacterPerformanceLowerName[35:48]: PerformSkill2,
	_CharacterPerformanceName[48:61]:      PerformSkill3,
	_CharacterPerformanceLowerName[48:61]: PerformSkill3,
	_CharacterPerformanceName[61:72]:      PerformHurt,
	_CharacterPerformanceLowerName[61:72]: PerformHurt,
	_CharacterPerformanceName[72:84]:      PerformDying,
	_CharacterPerformanceLowerName[72:84]: PerformDying,
	_CharacterPerformanceName[84:98]:      PerformVictory,
	_CharacterPerformanceLowerName[84:98]: PerformVictory,
}

var _CharacterPerformanceNames = []string{
	_CharacterPerformanceName[0:11],
	_CharacterPerformanceName[11:22],
	_CharacterPerformanceName[22:35],
	_CharacterPerformanceName[35:48],
	_CharacterPerformanceName[48:61],
	_CharacterPerformanceName[61:72],
	_CharacterPerformanceName[72:84],
	_CharacterPerformanceName[84:98],
}

// CharacterPerformanceString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func CharacterPerformanceString(s string) (CharacterPerformance, error) {
	if val, ok := _CharacterPerformanceNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _CharacterPerformanceNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to CharacterPerformance values", s)
}

// CharacterPerformanceValues returns all values of the enum
func CharacterPerformanceValues() []CharacterPerformance {
	return _CharacterPerformanceValues
}

// CharacterPerformanceStrings returns a slice of all String values of the enum
func CharacterPerformanceStrings() []string {
	strs := make([]string, len(_CharacterPerformanceNames))
	copy(strs, _CharacterPerformanceNames)
	return strs
}

// IsACharacterPerformance returns "true" if the value is listed in the enum definition. "false" otherwise
func (i CharacterPerformance) IsACharacterPerformance() bool {
	for _, v := range _CharacterPerformanceValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for CharacterPerformance
func (i CharacterPerformance) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for CharacterPerformance
func (i *CharacterPerformance) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("CharacterPerformance should be a string, got %s", data)
	}

	var err error
	*i, err = CharacterPerformanceString(s)
	return err
}
