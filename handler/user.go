package handler

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/model"
	"github.com/coralproject/pillar/service"
	"net/http"
)

//AddUser function adds a new user to the system
func AddUser(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	user := model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println("User: ", user)

	w.Header().Set("Content-Type", "application/json")

	dbUser, err := service.CreateUser(user)
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
