package service

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/server/config"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
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

func initDB() {
	file, err := os.Open("dbindex.json")
	if err != nil {
		log.Fatalf("Error opening file %s\n", err.Error())
	}

	objects := []model.Index{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading index information %v\n", err)
	}

	for _, one := range objects {
		if err := CreateIndex(&one); err != nil {
			log.Fatalf("Error creating indexes %v\n", err)
		}
	}
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

func CreateIndex(object *model.Index) *AppError {

	fmt.Printf("Object: %+v\n\n", object)

	manager := GetMongoManager()
	defer manager.Close()

	err := manager.Session.DB("").C(object.Target).EnsureIndex(object.Index)
	if err != nil {
		message := fmt.Sprintf("Error creating index [%+v]", object)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
