package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"
	"net/http"
	"fmt"
)

//AddAsset function adds a new user to the system
func AddAsset(w http.ResponseWriter, r *http.Request) {
	fmt.Print("In AddAsset")
	//Get the user from request
	jsonObject := model.Asset{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, statuscode and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateAsset(jsonObject)
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
