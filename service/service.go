package service

import (
	"fmt"
	"github.com/coralproject/pillar/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

const collectionUser string = "user"
const collectionAsset string = "asset"
const collectionComment string = "comment"

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

// FindCommentByID retrieves an individual comment resource
func FindCommentByID(id bson.ObjectId) (*model.Comment, error) {

	comment := model.Comment{}

	// Fetch comment
	manager := getMongoManager(collectionComment)
	defer manager.close()
	err := manager.collection.FindId(id).One(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// FindCommentBySourceID retrieves an individual comment resource
func FindCommentBySourceID(srcId string) (*model.Comment, error) {

	comment := model.Comment{}

	// Fetch comment
	manager := getMongoManager(collectionComment)
	defer manager.close()

	err := manager.collection.Find(bson.M{"src_id": srcId}).One(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// CreateComment creates a new comment resource
func CreateComment(comment model.Comment) (*model.Comment, error) {

	//return, if exists
	dbItem, _ := FindCommentByID(comment.ID)
	if dbItem != nil {
		fmt.Printf("Comment[%s] exists!", comment.ID)
		return dbItem, nil
	}

	//return, if exists
	dbItem, _ = FindCommentBySourceID(comment.SourceID)
	if dbItem != nil {
		fmt.Printf("Comment[%s] exists!", comment.SourceID)
		return dbItem, nil
	}

	// Insert Comment
	manager := getMongoManager(collectionComment)
	defer manager.close()

	comment.ID = bson.NewObjectId()
	err := manager.collection.Insert(comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// FindUserByID retrieves an individual user resource
func FindUserByID(id bson.ObjectId) (*model.User, error) {

	user := model.User{}

	// Fetch user
	manager := getMongoManager(collectionUser)
	defer manager.close()

	err := manager.collection.FindId(id).One(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindUserBySourceID retrieves an individual user resource
func FindUserBySourceID(srcId string) (*model.User, error) {

	user := model.User{}
	// Fetch user
	manager := getMongoManager(collectionUser)
	defer manager.close()
	err := manager.collection.Find(bson.M{"src_id": srcId}).One(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user resource
func CreateUser(user model.User) (*model.User, error) {

	//return, if exists
	dbItem, _ := FindUserByID(user.ID)
	if dbItem != nil {
		fmt.Printf("User[%s] exists!", user.ID)
		return dbItem, nil
	}

	//return, if exists
	dbItem, _ = FindUserBySourceID(user.SourceID)
	if dbItem != nil {
		fmt.Printf("User[%s] exists!", user.SourceID)
		return dbItem, nil
	}

	manager := getMongoManager(collectionUser)
	defer manager.close()
	user.ID = bson.NewObjectId()

	// Write the user to mongo
	err := manager.collection.Insert(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
