package main

import (
	"github.com/codegangsta/negroni"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/app/pillar/route"
	"github.com/coralproject/pillar/app/pillar/handler"
)

func main() {
	router := route.NewRouter()
	n := negroni.Classic()
	n.Use(handler.CORS())
	n.Use(handler.AppHandler())
	n.UseHandler(router)
	n.Run(config.GetAddress())
}
