package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

const collectionActor string = "actors"
const collectionComment string = "comments"

var (
	mgoSession *mgo.Session
)

// MongoManager encapsulates a mongo session with all relevant collections
type MongoManager struct {
	Session  *mgo.Session
	Actors   *mgo.Collection
	Comments *mgo.Collection
}

//Close closes the mongodb session; must be called, else the session remain open
func (manager *MongoManager) Close() {
	manager.Session.Close()
}

//export MONGODB_URL=mongodb://myuser:mypass@localhost:27017/echo
func init() {
	uri := os.Getenv("MONGODB_URL")
	if uri == "" {
		fmt.Println("Error connecting - MONGODB_URL not found!")
		os.Exit(1)
	}
	fmt.Printf("Connected to %s\n\n", uri)

	session, error := mgo.Dial(uri)
	if error != nil {
		panic(error) // no, not really
	}

	mgoSession = session
}

//GetMongoManager returns a cloned MongoManager
func GetMongoManager() *MongoManager {

	manager := MongoManager{}

	manager.Session = mgoSession.Clone()
	manager.Actors = manager.Session.DB("").C(collectionActor)
	manager.Comments = manager.Session.DB("").C(collectionComment)

	return &manager
}
