package nutritions

import (
	"encoding/json"
	"net/http"
	"shorvath/nutrition-tracker/repository"
	"shorvath/nutrition-tracker/util"
	"strconv"
)

const Prefix = "/nutritions"

func Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /":        listHandler,
		"GET /{id}":    findByIdHandler,
		"POST /":       createHandler,
		"PUT /":        updateHandler,
		"DELETE /{id}": deleteHandler,
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	list, err := repository.ListNutritions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, list); err != nil {
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
		if err = util.WriteJson(w, http.StatusOK, nutrition); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestNutrition *repository.Nutrition
	if err := json.NewDecoder(r.Body).Decode(requestNutrition); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nutrition, err := repository.CreateNutrition(r.Context(), requestNutrition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, nutrition); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var requestNutrition *repository.Nutrition
	if err := json.NewDecoder(r.Body).Decode(requestNutrition); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nutrition, err := repository.UpdateNutrition(r.Context(), requestNutrition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, nutrition); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam != "" {
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = repository.DeleteNutrition(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
	}
}
