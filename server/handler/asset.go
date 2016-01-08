package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"
	"net/http"
)

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
