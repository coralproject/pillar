package main

import (
	"github.com/coralproject/pillar/web/route"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {

	router := route.NewRouter()

	log.Printf(http.ListenAndServe(":8080", handlers.CORS()(router)).Error())
}
