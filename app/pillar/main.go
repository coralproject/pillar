package main

import (
	"github.com/coralproject/pillar/app/pillar/route"
	"log"
	"net/http"
)

func main() {

	router := route.NewRouter()

	log.Printf(http.ListenAndServe(":8080", router).Error())
}
