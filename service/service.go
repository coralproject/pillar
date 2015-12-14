package service

import (
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/coralproject/pillar/model"
)

const collectionUser string 	= "user"
const collectionAsset string 	= "asset"
const collectionComment string 	= "comment"

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
func FindCommentByID(id string) (*model.Comment, error) {

	comment := model.Comment{}

	// Fetch user
	manager := getMongoManager(collectionComment)
	defer manager.close()
	err := manager.collection.FindId(id).One(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// CreateComment creates a new comment resource
func CreateComment(comment model.Comment) (*model.Comment, error) {

	// Write the user to mongo
	manager := getMongoManager(collectionComment)
	defer manager.close()


	dbItem, _ := FindCommentByID(comment.CommentID)
	if dbItem != nil {
		fmt.Printf("Comment found:", dbItem.CommentID)
		return dbItem, nil
	}

	err := manager.collection.Insert(comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// FindUserByID retrieves an individual user resource
func FindUserByID(id string) (*model.User, error) {

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

// FindUserByEmail retrieves an individual user resource
func FindUserByEmail(email string) (*model.User, error) {

	user := model.User{}
	// Fetch user
	manager := getMongoManager(collectionUser)
	defer manager.close()
	err := manager.collection.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user resource
func CreateUser(user model.User) (*model.User, error) {
	manager := getMongoManager(collectionUser)
	defer manager.close()

	dbItem, _ := FindUserByID(user.UserID)
	if dbItem != nil {
		fmt.Printf("User found:", dbItem.UserID)
		return dbItem, nil
	}

	// Write the user to mongo
	err := manager.collection.Insert(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
