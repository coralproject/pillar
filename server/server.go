package main

import (
	"github.com/coralproject/pillar/server/log"
	"github.com/coralproject/pillar/server/web"
	"net/http"
)

func main() {

	router := web.NewRouter()

	log.Logger.Print(http.ListenAndServe(":8080", router))
}
