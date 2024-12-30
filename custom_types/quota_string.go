// Code generated by "stringer -type=Quota"; DO NOT EDIT.

package custom_types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Calories-0]
	_ = x[Fats-1]
	_ = x[FatsSaturated-2]
	_ = x[Carbs-3]
	_ = x[CarbsSugar-4]
	_ = x[CarbsSlowRelease-5]
	_ = x[CarbsFastRelease-6]
	_ = x[Proteins-7]
	_ = x[Salt-8]
}

const _Quota_name = "CaloriesFatsFatsSaturatedCarbsCarbsSugarCarbsSlowReleaseCarbsFastReleaseProteinsSalt"

var _Quota_index = [...]uint8{0, 8, 12, 25, 30, 40, 56, 72, 80, 84}

func (i Quota) String() string {
	if i < 0 || i >= Quota(len(_Quota_index)-1) {
		return "Quota(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Quota_name[_Quota_index[i]:_Quota_index[i+1]]
}