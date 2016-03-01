package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
)

func doRespond(w http.ResponseWriter, object interface{}, appErr *service.AppError) {
	if appErr != nil {
		config.Logger.Printf("Call failed [%+v]", appErr)
		payload, err := json.Marshal(appErr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(appErr.Code)
		w.Write(payload)
		return
	}

	payload, err := json.Marshal(object)
	if err != nil {
		config.Logger.Printf("Call failed [%+v]", err)
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
	jsonObject := model.CayUserAction{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	err := service.CreateUserAction(&jsonObject)
	doRespond(w, nil, err)
}
