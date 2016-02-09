package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/pkg/crud"
	"net/http"
)

func doRespond(w http.ResponseWriter, object interface{}, appErr *crud.AppError) {
	if appErr != nil {
		config.Logger.Printf("Call failed [%+v]", appErr)
		payload, err := json.Marshal(appErr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(appErr.Code)
		w.Write(payload)
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
	//Get the user from request
	jsonObject := crud.Index{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	err := crud.CreateIndex(&jsonObject)
	doRespond(w, nil, err)
}

