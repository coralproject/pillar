package route

import (
	"github.com/gorilla/mux"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/robfig/cron"
	"os"
	"strconv"
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

	scheduleCronJobs()

	return router
}

func scheduleCronJobs() {
	env := os.Getenv("PILLAR_CRON")
	b, err := strconv.ParseBool(env);
	if err != nil || !b {
		return
	}

	sched := os.Getenv("PILLAR_CRON_SEARCH")
	c := cron.New()
	c.AddFunc(sched, service.UpdateSearch)
	c.Start()
}
