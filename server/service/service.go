package service

import (
	"os"
	"gopkg.in/mgo.v2"
	"github.com/coralproject/pillar/server/log"
)

const collectionUser string = "user"
const collectionAsset string = "asset"
const collectionComment string = "comment"

var (
	mgoSession *mgo.Session
)

type MongoManager struct {
	Session  *mgo.Session
	Assets   *mgo.Collection
	Users    *mgo.Collection
	Comments *mgo.Collection
}

func (manager *MongoManager) Close() {
	manager.Session.Close()
}

//export MONGODB_URL=mongodb://localhost:27017/coral
func init() {
	uri := os.Getenv("MONGODB_URL")
	if uri == "" {
		log.Logger.Fatal("Error connecting to mongo database: MONGODB_URL not found")
	}

	session, err := mgo.Dial(uri)
	if err != nil {
		log.Logger.Fatalf("Error connecting to mongo database: %s", err)
	}

	mgoSession = session

	//url and src_id on Asset
	mgoSession.DB("").C(collectionAsset).EnsureIndexKey("src_id")
	mgoSession.DB("").C(collectionAsset).EnsureIndexKey("url")

	//src_id on User
	mgoSession.DB("").C(collectionUser).EnsureIndexKey("src_id")

	//source.id on Comment
	mgoSession.DB("").C(collectionComment).EnsureIndexKey("source.id")
}

func GetMongoManager() *MongoManager {

	manager := MongoManager{}

	manager.Session = mgoSession.Clone()
	manager.Assets = manager.Session.DB("").C(collectionAsset)
	manager.Users = manager.Session.DB("").C(collectionUser)
	manager.Comments = manager.Session.DB("").C(collectionComment)

	return &manager
}
