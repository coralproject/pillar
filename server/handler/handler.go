package handler

import (
	"encoding/json"
	"net/http"
	"github.com/coralproject/pillar/server/log"
)

// AppError - application specific error
type appError struct {
	Error   error
	Message string
	Code    int
}

func doRespond(w http.ResponseWriter, object interface{}, err error) {
	if err != nil {
		log.Logger.Printf("Error %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	payload, err := json.Marshal(object)
	if err != nil {
		log.Logger.Printf("Error %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
	w.Write(payload)
}

