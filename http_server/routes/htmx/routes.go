package htmx

import (
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
)

const Prefix = "/htmx"

func Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /":              rootHandler,
		"GET /meals":         mealsHandler,
		"GET /notifications": notificationsHandler,
		"GET /items":         itemsHandler,
	}
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	dailyQuota, err := repository.FindDailyQuotaByOwnerAndCurrentDay(r.Context(), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	meals, err := repository.FindMealsForUser(r.Context(), 1, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mealLogs, err := repository.FindMealLogsForUserAndCurrentDay(r.Context(), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	items, err := repository.ListItems(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mealLogsByMeal := util.GroupByKeys(mealLogs, meals, func(ml *repository.MealLog) *repository.Meal {
		meal, _ := util.FindFirst(meals, func(m *repository.Meal) bool { return ml.Meal.ID == m.ID })
		return meal
	})

	data := map[string]any{
		"Title":   "Meals",
		"TabName": "meals_tab",
		"Data": map[string]any{
			"DailyQuota":     dailyQuota,
			"Meals":          meals,
			"MealLogs":       mealLogs,
			"MealLogsByMeal": mealLogsByMeal,
			"Items":          items,
		},
	}

	err = repository.Render(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mealsHandler(w http.ResponseWriter, r *http.Request) {
	dailyQuota, err := repository.FindDailyQuotaByOwnerAndCurrentDay(r.Context(), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	meals, err := repository.FindMealsForUser(r.Context(), 1, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mealLogs, err := repository.FindMealLogsForUserAndCurrentDay(r.Context(), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mealLogsByMeal := util.GroupByKeys(mealLogs, meals, func(ml *repository.MealLog) *repository.Meal {
		meal, _ := util.FindFirst(meals, func(m *repository.Meal) bool { return ml.Meal.ID == m.ID })
		return meal
	})

	data := map[string]any{
		"Title":   "Meals",
		"TabName": "meals_tab",
		"Data": map[string]any{
			"DailyQuota":     dailyQuota,
			"Meals":          meals,
			"MealLogs":       mealLogs,
			"MealLogsByMeal": mealLogsByMeal,
		},
	}

	err = repository.Render(w, "meals_tab", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func notificationsHandler(w http.ResponseWriter, r *http.Request) {
	notifications, err := repository.ListNotificationsByUserId(r.Context(), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]any{
		"Title":   "Notifications",
		"TabName": "notifications_tab",
		"Data": map[string]any{
			"Notifications": notifications,
		},
	}

	err = repository.Render(w, "notifications_tab", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	items, err := repository.ListItems(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]any{
		"Title":   "Items",
		"TabName": "items_tab",
		"Data": map[string]any{
			"Items": items,
		},
	}

	err = repository.Render(w, "items_tab", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
