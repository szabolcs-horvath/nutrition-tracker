package v1

import (
	"net/http"
	"shorvath/nutrition-tracker/http_server/routes"
	"shorvath/nutrition-tracker/http_server/routes/api/v1/items"
	"shorvath/nutrition-tracker/http_server/routes/api/v1/notifications"
	"shorvath/nutrition-tracker/http_server/routes/api/v1/portions"
	"shorvath/nutrition-tracker/http_server/routes/api/v1/users"
)

const Prefix = "/v1"

func Routes() map[string]*http.ServeMux {
	return map[string]*http.ServeMux{
		users.Prefix:         routes.SubRouteHandlers(users.Handlers()),
		items.Prefix:         routes.SubRouteHandlers(items.Handlers()),
		portions.Prefix:      routes.SubRouteHandlers(portions.Handlers()),
		notifications.Prefix: routes.SubRouteHandlers(notifications.Handlers()),
	}
}
