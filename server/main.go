package main

import (
	"log"
	"net/http"

	"github.com/coralproject/pillar/server/web"
)

func main() {

	router := web.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
