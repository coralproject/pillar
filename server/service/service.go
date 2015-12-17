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

type mongoManager struct {
	session  *mgo.Session
	assets   *mgo.Collection
	users    *mgo.Collection
	comments *mgo.Collection
}

func (manager *mongoManager) close() {
	manager.session.Close()
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

func getMongoManager() *mongoManager {

	manager := mongoManager{}

	manager.session = mgoSession.Clone()
	manager.assets = manager.session.DB("").C(collectionAsset)
	manager.users = manager.session.DB("").C(collectionUser)
	manager.comments = manager.session.DB("").C(collectionComment)

	return &manager
}
