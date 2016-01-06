package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/coralproject/pillar/server/model"
)

// About - about this application
type About struct {
	app     string
	version string
}

var about About

func init() {
	about.app = "Coral Pillar Web Service"
	about.version = "Version - 0.0.1"
}

//AboutThisApp displays the about page
func AboutThisApp(w http.ResponseWriter, r *http.Request) {
	doRespond(w, about, nil)
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
