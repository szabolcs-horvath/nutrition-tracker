package repository

import (
	"github.com/szabolcs-horvath/nutrition-tracker/custom_types"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"html/template"
	"io"
	"time"
)

var templates *template.Template

func init() {
	templates = template.Must(template.New("templates").Funcs(util.TemplateFuncs()).Funcs(TemplateFuncs()).ParseGlob("web/templates/*.gohtml"))
}

func Render(w io.Writer, name string, data any) error {
	return templates.ExecuteTemplate(w, name, data)
}

func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"sumForQuota":                 sumForQuota,
		"remainingForMealQuota":       remainingForMealQuota,
		"remainingForDailyQuota":      remainingForDailyQuota,
		"percentageForMealQuota":      percentageForMealQuota,
		"percentageForDailyQuota":     percentageForDailyQuota,
		"percentageRemainingForQuota": percentageRemainingForQuota,
		"isClosestToCurrentTime":      isClosestToCurrentTime,
	}
}

func sumForQuota(quota custom_types.Quota, meallogs []*MealLog) float64 {
	var acc float64
	util.Reduce(meallogs, func(ml *MealLog, sum *float64) {
		*sum += ml.GetByQuota(quota)
	}, &acc)
	return acc
}

func remainingForMealQuota(quota custom_types.Quota, meallogs []*MealLog, meal *Meal) float64 {
	return *meal.Quotas[quota] - sumForQuota(quota, meallogs)
}

func remainingForDailyQuota(quota custom_types.Quota, meallogs []*MealLog, dailyQuota *DailyQuota) float64 {
	return *dailyQuota.Quotas[quota] - sumForQuota(quota, meallogs)
}

func percentageForMealQuota(quota custom_types.Quota, meallogs []*MealLog, meal *Meal) int64 {
	return util.Percentage(sumForQuota(quota, meallogs), *meal.Quotas[quota])
}

func percentageForDailyQuota(quota custom_types.Quota, meallogs []*MealLog, dailyQuota *DailyQuota) int64 {
	return util.Percentage(sumForQuota(quota, meallogs), *dailyQuota.Quotas[quota])
}

func percentageRemainingForQuota(quota custom_types.Quota, meallogs []*MealLog, meal *Meal) int64 {
	return util.PercentageRemaining(sumForQuota(quota, meallogs), *meal.Quotas[quota])
}

func isClosestToCurrentTime(meal *Meal, meals []*Meal) bool {
	type Closest struct {
		ID       *int64
		TimeDiff *time.Duration
	}

	now, err := custom_types.NewTime(time.Now())
	if err != nil {
		return false
	}
	winner := util.Reduce(meals, func(m *Meal, closest *Closest) {
		timeDiff := custom_types.TimeDiffAbs(m.Time, *now)
		if closest.TimeDiff == nil || timeDiff < *closest.TimeDiff {
			closest.ID = &m.ID
			closest.TimeDiff = &timeDiff
		}
	}, &Closest{})

	return winner.ID != nil && *winner.ID == meal.ID
}
