package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
	"github.com/coralproject/pillar/pkg/model"
)

func doRespond(w http.ResponseWriter, object interface{}, appErr *service.AppError) {
	if appErr != nil {
		config.Logger.Printf("Call failed [%v]", appErr)
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
		config.Logger.Printf("Call failed [%v]", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func CreateIndex(w http.ResponseWriter, r *http.Request) {
	var input model.Index
	json.NewDecoder(r.Body).Decode(&input)
	err := service.CreateIndex(GetAppContext(r, input))
	doRespond(w, nil, err)
}

func HandleUserAction(w http.ResponseWriter, r *http.Request) {
	var input model.Action
	json.NewDecoder(r.Body).Decode(&input)
	err := service.CreateUserAction(GetAppContext(r, input))
	doRespond(w, nil, err)
}

