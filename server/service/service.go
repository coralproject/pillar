package service

import (
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2"
	"os"
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
		fmt.Println("Error connecting - MONGODB_URL not found!")
		os.Exit(1)
	}

	session, error := mgo.Dial(uri)
	if error != nil {
		panic(error) // no, not really
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

func findOne(collection *mgo.Collection, query interface{}) model.DBType {

	var one model.DBType

	fmt.Printf("Session [%+v]\n", collection.Database.Session)

	collection.Find(query).One(&one)
	fmt.Printf("One [%+v]\n", one)

	return one
	//
	//	if one == nil {
	//		message := "{collection: " + collection.Name + ", query: " + query + "}"
	//		return errors.New("Object not found - " + message)
	//	}
	//
	//	return nil
}
