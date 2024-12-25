package meallogs

import (
	"encoding/json"
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
)

const Prefix = "/meallogs"

func Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"POST /": createHandler,
		//"PUT /":          updateHandler,
		//"DELETE /{id}":   deleteHandler,
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestItem repository.CreateMealLogRequest
	if err := json.NewDecoder(r.Body).Decode(&requestItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := repository.CreateMealLog(r.Context(), requestItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
