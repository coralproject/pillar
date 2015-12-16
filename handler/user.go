package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/model"
	"net/http"
	"github.com/coralproject/pillar/service"
)

//AddUser function adds a new user to the system
func AddUser(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.User{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateUser(jsonObject)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	payload, err := json.Marshal(dbObject)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(payload)
}
