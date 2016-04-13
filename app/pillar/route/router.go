package route

import (
	"github.com/gorilla/mux"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/robfig/cron"
)

//NewRouter returns a new mux.Router
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(route.HandlerFunc)
	}

	c := cron.New()
	c.AddFunc("@every 30m", service.UpdateSearch)
	c.Start()

	return router
}
