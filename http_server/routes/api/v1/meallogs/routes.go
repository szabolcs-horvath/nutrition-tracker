package meallogs

import (
	"github.com/szabolcs-horvath/nutrition-tracker/custom_types"
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
)

const Prefix = "/meallogs"

func HandlerFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /{id}":                   findByIdHandler,
		"GET /owner/{id}/date/{date}": listByOwnerAndDateHandler,
		"GET /owner/{id}/date/{$}":    listByOwnerAndDateHandler,
		"POST /{$}":                   createHandler,
		"PUT /{$}":                    updateHandler,
		"DELETE /{id}":                deleteHandler,
	}
}

func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	meallog, err := repository.FindMealLogById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, meallog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func listByOwnerAndDateHandler(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var result []*repository.MealLog
	dateParam := r.PathValue("date")
	if dateParam == "" {
		list, listErr := repository.FindMealLogsForUserAndCurrentDay(r.Context(), ownerId)
		if listErr != nil {
			http.Error(w, listErr.Error(), http.StatusInternalServerError)
			return
		}
		result = list
	} else {
		date, parseErr := custom_types.ParseDate(r.PathValue("date"))
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
			return
		}
		list, listErr := repository.FindMealLogsForUserAndDate(r.Context(), ownerId, date.UnderlyingTime())
		if listErr != nil {
			http.Error(w, listErr.Error(), http.StatusInternalServerError)
			return
		}
		result = list
	}
	if err = util.WriteJson(w, http.StatusOK, result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestMealLog repository.CreateMealLogRequest
	if err := util.ReadJson(r, &requestMealLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mealLog, err := repository.CreateMealLog(r.Context(), requestMealLog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, mealLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var requestMealLog repository.UpdateMealLogRequest
	if err := util.ReadJson(r, &requestMealLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mealLog, err := repository.UpdateMealLog(r.Context(), requestMealLog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, mealLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = repository.DeleteMealLog(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
