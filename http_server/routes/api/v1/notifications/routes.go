package notifications

import "net/http"

const Prefix = "/notifications"

func Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		//"GET /{id}":       findByIdHandler,
		//"GET /owner/{id}": listByOwnerHandler,
		//"POST /{$}":          createHandler,
		//"PUT /{$}":           updateHandler,
		//"DELETE /{id}":    deleteHandler,
	}
}
