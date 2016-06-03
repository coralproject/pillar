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

	schedse := os.Getenv("PILLAR_CRON_SEARCH")
	cse := cron.New()
	cse.AddFunc(schedse, service.UpdateSearch)
	cse.Start()

	schedst := os.Getenv("PILLAR_CRON_STATS")
	cst := cron.New()
	cst.AddFunc(schedst, service.CalculateStats)
	cst.Start()
}
