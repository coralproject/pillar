package web

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/pillar/server/handler"
	"github.com/gorilla/mux"
)

//Route defines mappings of end-points to handler methods
type Route struct {
	Name        string           `json:"name" bson:"name"`
	Method      string           `json:"method" bson:"method"`
	Pattern     string           `json:"pattern" bson:"pattern"`
	HandlerFunc http.HandlerFunc `json:"handler" bson:"handler"`
}

//Routes is an array of Route
type Routes []Route

//NewRouter returns a new mux.Router
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"About",
		"GET",
		"/about",
		handler.About,
	},
	Route{
		"Login",
		"POST",
		"/login",
		handler.Login,
	},
	Route{
		"Logout",
		"POST",
		"/logout",
		handler.Logout,
	},
	Route{
		"User",
		"POST",
		"/api/import/asset",
		handler.AddAsset,
	},
	Route{
		"User",
		"POST",
		"/api/import/user",
		handler.AddUser,
	},
	Route{
		"Comment",
		"POST",
		"/api/import/comment",
		handler.AddComment,
	},
}

func getRoutes() Routes {
	file, err := os.Open("routes.json")
	if err != nil {
		log.Error("routes", "getroutes", err, "Opening routes.json")
	}

	var routes Routes
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&routes); err != nil {
		log.Error("routes", "getroutes", err, "Parsing routes.json")
	}

	return routes
}
