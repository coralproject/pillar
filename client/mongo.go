package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

const collectionComment string = "comments"

var (
	mgoSession *mgo.Session
)

type mongoManager struct {
	session  *mgo.Session
	comments *mgo.Collection
}

func (manager *mongoManager) close() {
	manager.session.Close()
}

//export MONGODB_URL=mongodb://myuser:mypass@localhost:27017/echo
func init() {
	uri := os.Getenv("MONGODB_URL")
	if uri == "" {
		fmt.Println("Error connecting - MONGODB_URL not found!")
		os.Exit(1)
	}

	session, error := mgo.Dial(uri)
	if error != nil {
		panic(error) // no, not really
	}

	mgoSession = session
}

func getMongoManager() *mongoManager {

	manager := mongoManager{}

	manager.session = mgoSession.Clone()
	manager.comments = manager.session.DB("").C(collectionComment)

	return &manager
}
