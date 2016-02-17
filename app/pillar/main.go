package main

import (
	"github.com/codegangsta/negroni"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/app/pillar/route"
	"github.com/rs/cors"
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding"},
		AllowCredentials: true,
	})

	router := route.NewRouter()
	n := negroni.Classic()
	n.Use(c)
	n.UseHandler(router)
	n.Run(config.GetAddress())
}
