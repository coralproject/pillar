package main

import (
	"github.com/coralproject/pillar/server/web"
	"net/http"
	"log"
)

func main() {

	router := web.NewRouter()

	log.Printf(http.ListenAndServe(":8080", router).Error())
}
