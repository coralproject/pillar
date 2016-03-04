package handler

import (
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
	"github.com/coralproject/pillar/pkg/model"
	"encoding/json"
)

//ImportUser imports a new user to the system
func ImportUser(w http.ResponseWriter, r *http.Request) {
	var input model.User
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.ImportUser(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}

//ImportAsset imports a new asset to the system
func ImportAsset(w http.ResponseWriter, r *http.Request) {
	var input model.Asset
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.ImportAsset(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}

//ImportComment imports a new comment to the system
func ImportComment(w http.ResponseWriter, r *http.Request) {
	var input model.Comment
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.ImportComment(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}

//ImportAction imports actions into the system
func ImportAction(w http.ResponseWriter, r *http.Request) {
	var input model.Action
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.ImportAction(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}

//ImportNote imports notes into the system
func ImportNote(w http.ResponseWriter, r *http.Request) {
	var input model.Note
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.CreateNote(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}

//ImportMetadata imports metadata to various entities in the system
func ImportMetadata(w http.ResponseWriter, r *http.Request) {
	var input model.Metadata
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.UpdateMetadata(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}
