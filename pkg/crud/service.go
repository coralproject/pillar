package crud

import (
	"encoding/json"
	"fmt"
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
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
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

func init() {
	url := os.Getenv("MONGODB_URL")
	if url == "" {
		log.Fatal("Error initializing Service: MONGODB_URL not found.")
	}

	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("Error connecting to mongo database: %s", err)
	}

	mgoSession = session

	//url and source.id on Asset
	mgoSession.DB("").C(Actions).EnsureIndexKey("source.id")

	//url and source.id on Asset
	mgoSession.DB("").C(Assets).EnsureIndexKey("source.id")
	mgoSession.DB("").C(Assets).EnsureIndexKey("url")

	//source.id on User
	mgoSession.DB("").C(Users).EnsureIndexKey("source.id")

	//source.id on Comment
	mgoSession.DB("").C(Comments).EnsureIndexKey("source.id")

	//name on Tag
	mgoSession.DB("").C(Tags).EnsureIndexKey("name")

	//target_id, name and target,
	mgoSession.DB("").C(Tags).EnsureIndexKey("target_id", "name", "target")
}

func initDB() {
	file, err := os.Open("dbindex.json")
	if err != nil {
		log.Fatalf("Error opening file %s\n", err.Error())
	}

	objects := []Index{}
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
	manager.Users = manager.Session.DB("").C(Users)
	manager.Assets = manager.Session.DB("").C(Assets)
	manager.Actions = manager.Session.DB("").C(Actions)
	manager.Comments = manager.Session.DB("").C(Comments)
	manager.Tags = manager.Session.DB("").C(Tags)
	manager.TagTargets = manager.Session.DB("").C(TagTargets)

	return &manager
}

// UpdateMetadata updates metadata for an entity
func UpdateMetadata(object *Metadata) (interface{}, *AppError) {

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
func CreateIndex(object *Index) *AppError {
	manager := GetMongoManager()
	defer manager.Close()

	err := manager.Session.DB("").C(object.Target).EnsureIndex(object.Index)
	if err != nil {
		message := fmt.Sprintf("Error creating index [%+v]", object)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
