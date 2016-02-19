package service

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
)

const (
	DefaultMongoUrl string = "mongodb://localhost:27017/coral"
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
	Session    *mgo.Session
	Assets     *mgo.Collection
	Users      *mgo.Collection
	Actions    *mgo.Collection
	Comments   *mgo.Collection
	Tags       *mgo.Collection
	TagTargets *mgo.Collection
}

//Close closes the mongodb session; must be called, else the session remain open
func (manager *MongoManager) Close() {
	manager.Session.Close()
}

func initDB() {
	url := os.Getenv("MONGODB_URL")
	if url == "" {
		log.Printf("$MONGODB_URL not found, trying to connect locally [%s]", DefaultMongoUrl)
	}

	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("Error connecting to mongo database: %s", err)
	}

	mgoSession = session

	//source.id and target+type+user_id+target_id
	mgoSession.DB("").C(model.Actions).EnsureIndexKey("source.id")
	mgoSession.DB("").C(model.Actions).EnsureIndexKey("user_id", "target_id", "target", "type")

	//url and source.id on Asset
	mgoSession.DB("").C(model.Assets).EnsureIndexKey("source.id")
	mgoSession.DB("").C(model.Assets).EnsureIndexKey("url")

	//source.id on User
	mgoSession.DB("").C(model.Users).EnsureIndexKey("source.id")

	//source.id on Comment
	mgoSession.DB("").C(model.Comments).EnsureIndexKey("source.id")

	//name on Tag
	mgoSession.DB("").C(model.Tags).EnsureIndexKey("name")

	//target_id, name and target,
	mgoSession.DB("").C(model.Tags).EnsureIndexKey("target_id", "name", "target")
}

func initIndex() {
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
	if mgoSession == nil {
		initDB()
	}

	manager := MongoManager{}
	manager.Session = mgoSession.Clone()
	manager.Users = manager.Session.DB("").C(model.Users)
	manager.Assets = manager.Session.DB("").C(model.Assets)
	manager.Actions = manager.Session.DB("").C(model.Actions)
	manager.Comments = manager.Session.DB("").C(model.Comments)
	manager.Tags = manager.Session.DB("").C(model.Tags)
	manager.TagTargets = manager.Session.DB("").C(model.TagTargets)

	return &manager
}

// UpdateMetadata updates metadata for an entity
func UpdateMetadata(object *model.Metadata) (interface{}, *AppError) {

	manager := GetMongoManager()
	defer manager.Close()

	collection := manager.Session.DB("").C(object.Target)
	var dbEntity bson.M
	collection.FindId(object.TargetID).One(&dbEntity)
	if len(dbEntity) == 0 {
		collection.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	}

	if len(dbEntity) == 0 {
		message := fmt.Sprintf("Cannot update metadata for [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	collection.Update(
		bson.M{"_id": dbEntity["_id"]},
		bson.M{"$set": bson.M{"metadata": object.Metadata}},
	)

	return dbEntity, nil
}

// CreateIndex creates indexes to various entities
func CreateIndex(object *model.Index) *AppError {
	manager := GetMongoManager()
	defer manager.Close()

	err := manager.Session.DB("").C(object.Target).EnsureIndex(object.Index)
	if err != nil {
		message := fmt.Sprintf("Error creating index [%+v]", object)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
