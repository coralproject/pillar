package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"
	"net/http"
)

//ImportUser imports a new user to the system
func ImportUser(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.User{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateUser(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportAsset imports a new asset to the system
func ImportAsset(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Asset{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateAsset(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportComment imports a new comment to the system
func ImportComment(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Comment{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateComment(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportAction imports actions into the system
func ImportAction(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Action{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateAction(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportNote imports notes into the system
func ImportNote(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Note{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateNote(&jsonObject)
	doRespond(w, dbObject, err)
}
