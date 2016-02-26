package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
)

//CreateUpdateUser creates/updates an user in the system
func CreateUpdateUser(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.User{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateUpdateUser(&jsonObject)
	doRespond(w, dbObject, err)
}

//CreateUpdateAsset creates/updates an asset in the system
func CreateUpdateAsset(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Asset{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateUpdateAsset(&jsonObject)
	doRespond(w, dbObject, err)
}

//CreateUpdateComment creates/updates a comment in the system
func CreateUpdateComment(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Comment{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateUpdateComment(&jsonObject)
	doRespond(w, dbObject, err)
}
