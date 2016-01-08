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
