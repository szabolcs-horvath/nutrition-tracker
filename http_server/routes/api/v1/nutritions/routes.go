package nutritions

import (
	"net/http"
	"shorvath/nutrition-tracker/helpers"
	"shorvath/nutrition-tracker/repository"
	"strconv"
)

const Prefix = "/nutritions"

func Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /":     listHandler,
		"GET /{id}": findByIdHandler,
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	list, err := repository.ListNutritions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = helpers.WriteJson(w, http.StatusOK, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam != "" {
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		nutrition, err := repository.FindNutritionById(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err = helpers.WriteJson(w, http.StatusOK, nutrition); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
	}
}
