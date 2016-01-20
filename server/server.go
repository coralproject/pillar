package main

import (
	"github.com/coralproject/pillar/server/web"
	"log"
	"net/http"
)

func main() {

	router := web.NewRouter()

	log.Printf(http.ListenAndServe(":8080", router).Error())
}
