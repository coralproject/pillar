package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/crud"
	"net/http"
)

func TagsPreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
  w.Header().Set("Access-Control-Allow-Headers",
      "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-CoralCustom")
	doRespond(w, nil, nil)
}


func GetTags(w http.ResponseWriter, r *http.Request) {
	jsonObject := crud.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dbObject, err := crud.GetTags()
	doRespond(w, dbObject, err)
}

func UpsertTag(w http.ResponseWriter, r *http.Request) {
	//Get the tag from request
	jsonObject := crud.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dbObject, err := crud.UpsertTag(&jsonObject)
	doRespond(w, dbObject, err)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	//Get the tag from request
	jsonObject := crud.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := crud.DeleteTag(&jsonObject)
	doRespond(w, nil, err)
}

