package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/model"
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
	//Get the Index from request
	jsonObject := model.Index{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	err := service.CreateIndex(&jsonObject)
	doRespond(w, nil, err)
}

func HandleUserAction(w http.ResponseWriter, r *http.Request) {
	//Get the UserAction from request
	jsonObject := model.UserAction{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	err := service.CreateUserAction(&jsonObject)
	doRespond(w, nil, err)
}
