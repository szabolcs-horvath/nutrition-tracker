package meallogs

import (
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
)

const Prefix = "/meallogs"

func Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"POST /{$}":    createHandler,
		"PUT /{$}":     updateHandler,
		"DELETE /{id}": deleteHandler,
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestMealLog repository.CreateMealLogRequest
	if err := util.ReadJson(r, requestMealLog); err != nil {
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
	if err := util.ReadJson(r, requestMealLog); err != nil {
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
