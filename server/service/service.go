package service

import (
	"github.com/coralproject/pillar/server/cfg"
	"gopkg.in/mgo.v2"
	"log"
	"github.com/coralproject/pillar/server/model"
)

// AppError encapsulates application specific error
type AppError struct {
	Error   error
	Message string
	Code    int
}

var (
	mgoSession *mgo.Session
)

// MongoManager encapsulates a mongo session with all relevant collections
type MongoManager struct {
	Session  *mgo.Session
	Assets   *mgo.Collection
	Users    *mgo.Collection
	Actions  *mgo.Collection
	Comments *mgo.Collection
}

//Close closes the mongodb session; must be called, else the session remain open
func (manager *MongoManager) Close() {
	manager.Session.Close()
}

//export MONGODB_URL=mongodb://localhost:27017/coral
func init() {
	session, err := mgo.Dial(cfg.GetContext().MongoURL)
	if err != nil {
		log.Fatalf("Error connecting to mongo database: %s", err)
	}

	mgoSession = session

	//url and src_id on Asset
	mgoSession.DB("").C(model.CollectionAction).EnsureIndexKey("source.id")

	//url and src_id on Asset
	mgoSession.DB("").C(model.CollectionAsset).EnsureIndexKey("source.id")
	mgoSession.DB("").C(model.CollectionAsset).EnsureIndexKey("url")

	//src_id on User
	mgoSession.DB("").C(model.CollectionUser).EnsureIndexKey("source.id")

	//source.id on Comment
	mgoSession.DB("").C(model.CollectionComment).EnsureIndexKey("source.id")
}

//GetMongoManager returns a cloned MongoManager
func GetMongoManager() *MongoManager {

	manager := MongoManager{}
	manager.Session = mgoSession.Clone()
	manager.Users = manager.Session.DB("").C(model.CollectionUser)
	manager.Assets = manager.Session.DB("").C(model.CollectionAsset)
	manager.Actions = manager.Session.DB("").C(model.CollectionAction)
	manager.Comments = manager.Session.DB("").C(model.CollectionComment)

	return &manager
}
