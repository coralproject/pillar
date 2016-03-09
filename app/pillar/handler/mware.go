package handler

import (
	"github.com/codegangsta/negroni"
	"github.com/coralproject/pillar/pkg/web"
	"github.com/gorilla/context"
	"github.com/rs/cors"
	"net/http"
)

const appContext string = "app-context"

//CORS is a handler to take care of CORS
func CORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding"},
		AllowCredentials: true,
	})
}

//AppHandler is a handler to inject gorilla context to the request.
func AppHandler() negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		// set header
		rw.Header().Set("Content-Type", "application/json")

		//Create/Inject DB
		c := web.NewContext(r.Header, r.Body)
		defer c.Close()

		//Create/Inject AppContext for the Handlers to use
		SetAppContext(r, c)
		next(rw, r)
	})
}

//GetAppContext returns an AppContext from the global gorilla context.
func GetAppContext(r *http.Request, input interface{}) *web.AppContext {

	rc := context.Get(r, appContext)
	if rc == nil {
		return nil
	}

	//inject input data if any
	c := rc.(*web.AppContext)
	c.Marshall(input)
	return c
}

//SetAppContext sets an AppContext to the global gorilla context.
func SetAppContext(r *http.Request, val *web.AppContext) {
	context.Set(r, appContext, val)
}
