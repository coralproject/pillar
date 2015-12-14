package model

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"gopkg.in/bluesuncorp/validator.v6"
)

const collectionUser string 	= "user"
const collectionAsset string 	= "asset"
const collectionComment string 	= "comment"

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	config := validator.Config{
		TagName:         "validate",
		ValidationFuncs: validator.BakedInValidators,
	}

	validate = validator.New(config)
}

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

func getMongoManager(collectionName string) *mongoManager {
	session := mgoSession.Clone()
	collection := session.DB("").C(collectionName)
	return &mongoManager{session, collection}
}

//==============================================================================
