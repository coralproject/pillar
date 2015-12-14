package main

import (
	"github.com/coralproject/pillar/web"
	"log"
	"net/http"
)

func main() {

	router := web.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
