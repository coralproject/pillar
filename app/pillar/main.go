package main

import (
	"github.com/codegangsta/negroni"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/app/pillar/handler"
	"github.com/coralproject/pillar/app/pillar/route"
)

const (
	// VersionNumber is the version for sponge
	VersionNumber = 0.1
)

func main() {
	//new Negroni Middleware
	n := negroni.Classic()
	n.Use(negroni.NewLogger())

	//Add CORS
	n.Use(handler.CORS())

	//Router at the end
	n.UseHandler(route.NewRouter())

	//run server
	n.Run(config.Address())
}
