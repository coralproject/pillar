package service

import (
	"errors"
	"os"

	"github.com/ardanlabs/kit/log"
	"gopkg.in/mgo.v2"
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
		log.Error("start", "init", errors.New("Error connecting - MONGODB_URL not found!"), "Getting MONGODB_URL env variable.")
		os.Exit(1)
	}

	session, err := mgo.Dial(uri)
	if err != nil {
		log.Error("start", "init", err, "Connecting to mongo")
		panic(err) // no, not really <--- do we really need to panic?
	}

	mgoSession = session
}

func GetMongoManager() *MongoManager {

	manager := MongoManager{}

	manager.Session = mgoSession.Clone()
	manager.Assets = manager.Session.DB("").C(collectionAsset)
	manager.Assets.EnsureIndexKey("src_id")
	manager.Assets.EnsureIndexKey("url")

	manager.Users = manager.Session.DB("").C(collectionUser)
	manager.Users.EnsureIndexKey("src_id")

	manager.Comments = manager.Session.DB("").C(collectionComment)

	return &manager
}
