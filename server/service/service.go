package service

import (
	"fmt"
	"github.com/coralproject/pillar/server/config"
	"github.com/coralproject/pillar/server/dto"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
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
	session, err := mgo.Dial(config.GetContext().MongoURL)
	if err != nil {
		log.Fatalf("Error connecting to mongo database: %s", err)
	}

	mgoSession = session

	//url and source.id on Asset
	mgoSession.DB("").C(model.CollectionAction).EnsureIndexKey("source.id")

	//url and source.id on Asset
	mgoSession.DB("").C(model.CollectionAsset).EnsureIndexKey("source.id")
	mgoSession.DB("").C(model.CollectionAsset).EnsureIndexKey("url")

	//source.id on User
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

// UpdateMetadata updates metadata to an entity
func UpdateMetadata(object *dto.Metadata) (interface{}, *AppError) {

	switch object.Target {

	case model.CollectionAsset:
		return updateAssetMetadata(object)
	case model.CollectionAction:
		return updateActionMetadata(object)
	case model.CollectionUser:
		return updateUserMetadata(object)
	case model.CollectionComment:
		return updateCommentMetadata(object)
	}

	message := fmt.Sprintf("Invalid metadata [%+v]\n", object)
	return nil, &AppError{nil, message, http.StatusInternalServerError}
}
