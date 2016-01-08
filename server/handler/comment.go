package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"
	"net/http"
)

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
