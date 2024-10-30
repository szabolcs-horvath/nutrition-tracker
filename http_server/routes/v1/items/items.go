package items

import (
	"encoding/json"
	"net/http"
	"shorvath/nutrition-tracker/helpers"
	"shorvath/nutrition-tracker/repository"
	"strconv"
)

func Prefix() string {
	return "/items"
}

func Routes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return map[string]func(w http.ResponseWriter, r *http.Request){
		"GET /{id}": func(w http.ResponseWriter, r *http.Request) {
			idParam := r.PathValue("id")
			if idParam != "" {
				id, err := strconv.ParseInt(idParam, 10, 64)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				item, err := repository.FindItemByIdWithNutrition(r.Context(), id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				if err = helpers.WriteJson(w, http.StatusOK, item); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
			} else {
				http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
			}
		},
		"POST /": func(w http.ResponseWriter, r *http.Request) {
			var requestItem repository.Item
			if err := json.NewDecoder(r.Body).Decode(&requestItem); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			item, err := repository.CreateItem(r.Context(), requestItem)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err = helpers.WriteJson(w, http.StatusCreated, item); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		},
	}
}
