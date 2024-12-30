package custom_types

import "github.com/szabolcs-horvath/nutrition-tracker/util"

//go:generate stringer -type=Quota

type Quota int64

const (
	Calories Quota = iota
	Fats
	FatsSaturated
	Carbs
	CarbsSugar
	CarbsSlowRelease
	CarbsFastRelease
	Proteins
	Salt
)

var AllQuotas = []Quota{
	Calories,
	Fats,
	FatsSaturated,
	Carbs,
	CarbsSugar,
	CarbsSlowRelease,
	CarbsFastRelease,
	Proteins,
	Salt,
}

var AllQuotaStrings = util.Map(AllQuotas, func(quota Quota) string { return quota.String() })
