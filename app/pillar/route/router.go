package route

import (
	"os"
	"strconv"

	"github.com/coralproject/pillar/pkg/service"
	"github.com/gorilla/mux"
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

	scheduleCronJobs()

	return router
}

// Adds cron tasks and start CRON. For it to work the environment variable PILLAR_CRON has to be set to True
func scheduleCronJobs() {
	env := os.Getenv("PILLAR_CRON")
	b, err := strconv.ParseBool(env)
	if err != nil || !b {
		return
	}

	sched := os.Getenv("PILLAR_CRON_SEARCH")
	c := cron.New()
	c.AddFunc(sched, service.UpdateSearch)
	c.AddFunc(sched, service.CalculateStats)
	c.Start()
}
