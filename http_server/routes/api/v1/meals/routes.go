package meals

import (
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
	"strings"
)

const Prefix = "/meals"

func HandlerFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /{id}":       findByIdHandler,
		"GET /owner/{id}": listByOwnerHandler,
		"POST /{$}":       createHandler,
		"PUT /{$}":        updateHandler,
		"DELETE /{id}":    deleteHandler,
	}
}

func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	meal, err := repository.FindMealById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, meal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func listByOwnerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	archivedQueryParam := r.URL.Query().Get("archived")
	archived := "true" == strings.ToLower(archivedQueryParam)
	list, err := repository.FindMealsForUser(r.Context(), id, archived)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestMeal repository.CreateMealRequest
	if err := util.ReadJson(r, &requestMeal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	meal, err := repository.CreateMeal(r.Context(), requestMeal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, meal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var requestMeal repository.UpdateMealRequest
	if err := util.ReadJson(r, &requestMeal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	meal, err := repository.UpdateMeal(r.Context(), requestMeal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, meal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = repository.ArchiveMeal(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
