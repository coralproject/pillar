package model

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

const collectionAsset string = "asset"
const collectionUser string = "user"
const collectionComment string = "comment"

var (
	mgoSession *mgo.Session
)

type mongoManager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func (manager *mongoManager) close() {
	manager.session.Close()
}

//export MONGODB_URL=mongodb://localhost:27017/coral
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

func getSession() *mgo.Session {
	return mgoSession.Clone()
}

func getMongoManager(collectionName string) *mongoManager {
	session := getSession()
	collection := session.DB("").C(collectionName)
	return &mongoManager{session, collection}
}

//==============================================================================
