// Code generated by "stringer -type=Class"; DO NOT EDIT.

package item

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnarmedClass-0]
	_ = x[SwordClass-1]
	_ = x[AxeClass-2]
	_ = x[ClubClass-3]
	_ = x[DaggerClass-4]
	_ = x[SlingClass-5]
	_ = x[BowClass-6]
	_ = x[SpearClass-7]
	_ = x[PolearmClass-8]
	_ = x[StaffClass-9]
	_ = x[WandClass-10]
	_ = x[HatClass-11]
	_ = x[BodyArmorClass-12]
	_ = x[AmuletClass-13]
	_ = x[RingClass-14]
	_ = x[GloveClass-15]
	_ = x[BootClass-16]
	_ = x[BeltClass-17]
}

const _Class_name = "UnarmedClassSwordClassAxeClassClubClassDaggerClassSlingClassBowClassSpearClassPolearmClassStaffClassWandClassHatClassBodyArmorClassAmuletClassRingClassGloveClassBootClassBeltClass"

var _Class_index = [...]uint8{0, 12, 22, 30, 39, 50, 60, 68, 78, 90, 100, 109, 117, 131, 142, 151, 161, 170, 179}

func (i Class) String() string {
	if i < 0 || i >= Class(len(_Class_index)-1) {
		return "Class(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Class_name[_Class_index[i]:_Class_index[i+1]]
}
