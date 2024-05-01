package calendar

import (
	"net/http"

	"github.com/Henry-Le-CS/gost"
)

func CalendarControllers() *gost.Controller {
	router := gost.DeclareRouter()

	router.Add(gost.DeclareRouteHandler(
		"POST",
		"/",
		generateCalendarHandler,
	))

	controllers := gost.DeclareController(gost.ControllerArgs{
		Prefix: "/calendar",
		Router: router,
	})

	return controllers
}

func generateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}