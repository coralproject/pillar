package main

import (
	"github.com/codegangsta/negroni"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/app/pillar/handler"
	"github.com/coralproject/pillar/app/pillar/route"
)

func main() {
	//new Negroni Middleware
	n := negroni.Classic()

	//Add CORS
	n.Use(handler.CORS())

	//Router at the end
	n.UseHandler(route.NewRouter())

	//run server
	n.Run(config.Address())
}
