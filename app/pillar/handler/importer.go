package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/crud"
	"net/http"
)

//ImportUser imports a new user to the system
func ImportUser(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := crud.User{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.CreateUser(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportAsset imports a new asset to the system
func ImportAsset(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := crud.Asset{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.CreateAsset(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportComment imports a new comment to the system
func ImportComment(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := crud.Comment{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.CreateComment(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportAction imports actions into the system
func ImportAction(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := crud.Action{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.CreateAction(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportNote imports notes into the system
func ImportNote(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := crud.Note{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.CreateNote(&jsonObject)
	doRespond(w, dbObject, err)
}

//ImportMetadata imports metadata to various entities in the system
func ImportMetadata(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := crud.Metadata{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.UpdateMetadata(&jsonObject)
	doRespond(w, dbObject, err)
}
