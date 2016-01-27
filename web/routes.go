package web

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/coralproject/pillar/handler"
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
		handler.AboutThisApp,
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
		handler.ImportAsset,
	},
	Route{
		"User",
		"POST",
		"/api/import/user",
		handler.ImportUser,
	},
	Route{
		"Comment",
		"POST",
		"/api/import/comment",
		handler.ImportComment,
	},
	Route{
		"Action",
		"POST",
		"/api/import/action",
		handler.ImportAction,
	},
	Route{
		"Note",
		"POST",
		"/api/import/note",
		handler.ImportNote,
	},
	Route{
		"Metadata",
		"POST",
		"/api/import/metadata",
		handler.ImportMetadata,
	},
	Route{
		"Index",
		"POST",
		"/api/import/index",
		handler.CreateIndex,
	},
}

func getRoutes() Routes {
	file, err := os.Open("routes.json")
	if err != nil {
		log.Fatal("Error opening routes.json: %s", err)
	}

	var routes Routes
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&routes); err != nil {
		log.Fatal("Error parsing routes.json: %s", err)
	}

	return routes
}
