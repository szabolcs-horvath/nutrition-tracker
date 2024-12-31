package htmx

import (
	"fmt"
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
)

const Prefix = "/htmx"

func Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /":                        rootHandler,
		"GET /today":                   todayHandler,
		"GET /notifications":           notificationsHandler,
		"GET /items":                   itemsHandler,
		"POST /items/search":           itemSearchHandler,
		"POST /meallogs/meal/{mealId}": addMealLogForMealHandler,
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
		"Title":   "Today",
		"TabName": "today_tab",
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

func todayHandler(w http.ResponseWriter, r *http.Request) {
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
		"Title":   "Today",
		"TabName": "today_tab",
		"Data": map[string]any{
			"DailyQuota":     dailyQuota,
			"Meals":          meals,
			"MealLogs":       mealLogs,
			"MealLogsByMeal": mealLogsByMeal,
		},
	}

	err = repository.Render(w, "today_tab", data)
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

func itemSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	if len(query) < 2 {
		http.Error(w, fmt.Errorf("the search string has to be at least 2 characters").Error(), http.StatusBadRequest)
		return
	}
	results, err := repository.SearchItemsByNameAndUser(r.Context(), 1, query)

	var meal *repository.Meal
	if mealIdParam := r.FormValue("meal_id"); mealIdParam != "" {
		mealId, intParseErr := strconv.ParseInt(mealIdParam, 10, 64)
		if intParseErr != nil {
			http.Error(w, fmt.Errorf("failure while trying to parse meal_id: %s", intParseErr.Error()).Error(), http.StatusBadRequest)
			return
		}
		meal, err = repository.FindMealById(r.Context(), mealId)
	}

	data := map[string]any{
		"Data": map[string]any{
			"Meal":          meal,
			"SearchResults": results,
		},
	}

	err = repository.Render(w, "item_search_results", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addMealLogForMealHandler(w http.ResponseWriter, r *http.Request) {
	mealIdParam := r.PathValue("mealId")
	mealId, err := strconv.ParseInt(mealIdParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var requestMealLog repository.CreateMealLogRequest
	if err = util.ReadJson(r, &requestMealLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = repository.CreateMealLog(r.Context(), requestMealLog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	meallogs, err := repository.FindMealLogsForMealAndCurrentDay(r.Context(), mealId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repository.Render(w, "meallogs_simple", meallogs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
