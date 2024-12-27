package users

import (
	"github.com/szabolcs-horvath/nutrition-tracker/repository"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"strconv"
)

const Prefix = "/users"

func HandlerFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /{$}":     listHandler,
		"GET /{id}":    findByIdHandler,
		"POST /{$}":    createHandler,
		"PUT /{$}":     updateHandler,
		"DELETE /{id}": deleteHandler,
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	list, err := repository.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := repository.FindUserById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var requestUser repository.CreateUserRequest
	if err := util.ReadJson(r, &requestUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := repository.CreateUser(r.Context(), requestUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusCreated, user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var requestUser repository.UpdateUserRequest
	if err := util.ReadJson(r, &requestUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := repository.UpdateUser(r.Context(), requestUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = util.WriteJson(w, http.StatusOK, user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = repository.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
