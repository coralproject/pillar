package main

import (
	"github.com/codegangsta/negroni"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/app/pillar/route"
	"github.com/coralproject/pillar/app/pillar/handler"
)

func main() {
	//new Negroni Middleware
	n := negroni.Classic()

	//Add CORS and custom AppHandler
	n.Use(handler.CORS())
	n.Use(handler.AppHandler())

	//Router at the end
	n.UseHandler(route.NewRouter())

	//run server
	n.Run(config.GetAddress())
}
