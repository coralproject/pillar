package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/crud"
	"net/http"
)

func doRespond(w http.ResponseWriter, object interface{}, appErr *crud.AppError) {
	if appErr != nil {
		Logger.Printf("Call failed [%s]", appErr.Message)
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	payload, err := json.Marshal(object)
	if err != nil {
		Logger.Printf("Call failed [%s]", err.Error())
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

func CreateTag(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := crud.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.UpsertTag(&jsonObject)
	doRespond(w, dbObject, err)
}
