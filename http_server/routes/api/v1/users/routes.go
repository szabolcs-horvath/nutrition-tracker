package users

import "net/http"

const Prefix = "/users"

func Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		//"GET /{$}":        listHandler,
		//"GET /{id}":    findByIdHandler,
		//"POST /{$}":       createHandler,
		//"PUT /{$}":        updateHandler,
		//"DELETE /{id}": deleteHandler,
	}
}
