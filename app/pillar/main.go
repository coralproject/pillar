package main

import (
	"github.com/coralproject/pillar/app/pillar/route"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"github.com/coralproject/pillar/app/pillar/config"
)

func main() {

	router := route.NewRouter()

	log.Fatalf(http.ListenAndServe(config.GetAddress(), handlers.CORS()(router)).Error())
}
