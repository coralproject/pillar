package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/server/log"
	"github.com/coralproject/pillar/server/service"
	"net/http"
)

func doRespond(w http.ResponseWriter, object interface{}, appErr *service.AppError) {
	if appErr != nil {
		log.Logger.Printf("Call failed [%s]", appErr.Message)
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	payload, err := json.Marshal(object)
	if err != nil {
		log.Logger.Printf("Call failed [%s]", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
