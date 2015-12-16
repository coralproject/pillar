package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"
	"net/http"
)

//AddComment function adds a new comment to the system
func AddComment(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Comment{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateComment(jsonObject)
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
