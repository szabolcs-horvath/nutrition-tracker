package routes

import (
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/htmx", http.StatusTemporaryRedirect)
}
