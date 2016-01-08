package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"
	"net/http"
)

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
