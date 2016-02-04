package main

import (
	"github.com/coralproject/pillar/app/pillar/route"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {

	router := route.NewRouter()

	log.Fatalf(http.ListenAndServe(":8080", handlers.CORS()(router)).Error())
}
