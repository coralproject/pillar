package handler

import (
	"github.com/rs/cors"
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

//func AppHandler() negroni.Handler {
//	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//
//		// set header
//		rw.Header().Set("Content-Type", "application/json")
//
//		//Create/Inject DB
//		c := service.NewContext()
//		defer c.Close()
//
//		//Create/Inject AppContext for the Handlers to use
//		SetAppContext(r, c)
//		next(rw, r)
//	})
//}
//
//func GetAppContext(r *http.Request, input interface{}) *service.AppContext {
//
//	rc := context.Get(r, appContext)
//	if rc == nil {
//		return nil
//	}
//
//	c := rc.(*service.AppContext)
//	//inject input data if any
//	if input != nil {
//		c.Input = input
//	}
//
//	return c
//}
//
//func SetAppContext(r *http.Request, val *service.AppContext) {
//	context.Set(r, appContext, val)
//}
