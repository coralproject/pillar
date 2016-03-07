package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
)

//CreateUpdateUser creates/updates an user in the system
func CreateUpdateUser(w http.ResponseWriter, r *http.Request) {
	var input model.User
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.CreateUpdateUser(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}

//CreateUpdateAsset creates/updates an asset in the system
func CreateUpdateAsset(w http.ResponseWriter, r *http.Request) {
	var input model.Asset
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.CreateUpdateAsset(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}

//CreateUpdateComment creates/updates a comment in the system
func CreateUpdateComment(w http.ResponseWriter, r *http.Request) {
	var input model.Comment
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.CreateUpdateComment(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}
