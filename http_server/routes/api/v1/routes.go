package v1

import (
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/items"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/meallogs"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/notifications"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/portions"
	"github.com/szabolcs-horvath/nutrition-tracker/http_server/routes/api/v1/users"
	"net/http"
)

const Prefix = "/v1"

func Routes() map[string]*http.ServeMux {
	return map[string]*http.ServeMux{
		users.Prefix:         routes.SubRouteHandlers(users.Handlers()),
		items.Prefix:         routes.SubRouteHandlers(items.Handlers()),
		portions.Prefix:      routes.SubRouteHandlers(portions.Handlers()),
		notifications.Prefix: routes.SubRouteHandlers(notifications.Handlers()),
		meallogs.Prefix:      routes.SubRouteHandlers(meallogs.Handlers()),
	}
}
