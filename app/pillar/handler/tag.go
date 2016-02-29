package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
)

func GetTags(w http.ResponseWriter, r *http.Request) {
	jsonObject := model.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.GetTags()
	doRespond(w, dbObject, err)
}

func CreateUpdateTag(w http.ResponseWriter, r *http.Request) {
	//Get the tag from request
	jsonObject := model.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateUpdateTag(&jsonObject)
	doRespond(w, dbObject, err)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	//Get the tag from request
	jsonObject := model.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	err := service.DeleteTag(&jsonObject)
	doRespond(w, nil, err)
}

