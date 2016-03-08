package handler

import (
	"github.com/codegangsta/negroni"
	"github.com/coralproject/pillar/pkg/web"
	"github.com/gorilla/context"
	"github.com/rs/cors"
	"net/http"
)

const appContext string = "app-context"

func CORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding"},
		AllowCredentials: true,
	})
}

func AppHandler() negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		// set header
		rw.Header().Set("Content-Type", "application/json")

		//Create/Inject DB
		c := web.NewContext(nil)
		defer c.Close()

		//Create/Inject AppContext for the Handlers to use
		SetAppContext(r, c)
		next(rw, r)
	})
}

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

func SetAppContext(r *http.Request, val *web.AppContext) {
	context.Set(r, appContext, val)
}
