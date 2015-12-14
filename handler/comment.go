package handler

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/model"
	"net/http"
)

//AddComment function adds a new comment to the system
func AddComment(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	comment := model.Comment{}
	json.NewDecoder(r.Body).Decode(&comment)
	fmt.Println("Comment: ", comment)

	w.Header().Set("Content-Type", "application/json")

	dbUser, err := model.CreateComment(comment)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	js, err := json.Marshal(dbUser)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Write content-type, statuscode, payload
	w.WriteHeader(200)
	w.Write(js)
}
