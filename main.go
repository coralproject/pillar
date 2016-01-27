package main

import (
	"github.com/coralproject/pillar/web"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {

	router := web.NewRouter()

	log.Printf(http.ListenAndServe(":8080", handlers.CORS()(router)).Error())
}
