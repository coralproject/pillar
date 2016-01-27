package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/config"
	"github.com/coralproject/pillar/model"
	"github.com/coralproject/pillar/service"
	"net/http"
)

func doRespond(w http.ResponseWriter, object interface{}, appErr *service.AppError) {
	if appErr != nil {
		config.Logger.Printf("Call failed [%s]", appErr.Message)
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	payload, err := json.Marshal(object)
	if err != nil {
		config.Logger.Printf("Call failed [%s]", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func CreateIndex(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Index{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	err := service.CreateIndex(&jsonObject)
	doRespond(w, nil, err)
}
