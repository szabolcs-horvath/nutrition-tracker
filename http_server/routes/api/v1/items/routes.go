package items

import (
	"encoding/json"
	"net/http"
	"shorvath/nutrition-tracker/helpers"
	"shorvath/nutrition-tracker/repository"
	"strconv"
)

const Prefix = "/items"

func Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /": func(w http.ResponseWriter, r *http.Request) {
			list, err := repository.ListItems(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err = helpers.WriteJson(w, http.StatusOK, list); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		},
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
			var requestItem *repository.Item
			if err := json.NewDecoder(r.Body).Decode(requestItem); err != nil {
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
		"PUT /": func(w http.ResponseWriter, r *http.Request) {
			var requestItem *repository.Item
			if err := json.NewDecoder(r.Body).Decode(requestItem); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			item, err := repository.UpdateItem(r.Context(), requestItem)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err = helpers.WriteJson(w, http.StatusOK, item); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		},
		"DELETE /{id}": func(w http.ResponseWriter, r *http.Request) {
			idParam := r.PathValue("id")
			if idParam != "" {
				id, err := strconv.ParseInt(idParam, 10, 64)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				err = repository.DeleteItem(r.Context(), id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			} else {
				http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
			}
		},
	}
}
