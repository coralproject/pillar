package handler

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"net/http"
)

//About shows the about page
func About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Comment Demo App, Version - 0.0.1")
}

//Logout logs the user out of the system
func Logout(w http.ResponseWriter, r *http.Request) {
}

//Login logs the user into the system
func Login(w http.ResponseWriter, r *http.Request) {

	//Get the user from request
	user := model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println("User: ", user)

	w.Header().Set("Content-Type", "application/json")

	//	dbUser, err := model.Login(user)
	//	if err != nil {
	//		w.WriteHeader(401)
	//		return
	//	}
	//
	//	js, err := json.Marshal(dbUser)
	//	if err != nil {
	//		w.WriteHeader(500)
	//		return
	//	}

	// Write content-type, statuscode, payload
	w.WriteHeader(200)
	//w.Write(js)
}
