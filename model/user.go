package model

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

// User denotes a user in the system.
type User struct {
	UserID      string                 `json:"user_id" bson:"user_id"`
	UserName    string                 `json:"user_name" bson:"user_name" validate:"required"`
	Avatar      string                 `json:"avatar" bson:"avatar" validate:"omitempty,url"`
	LastLogin   time.Time              `json:"last_login" bson:"last_login"`
	MemberSince time.Time              `json:"member_since" bson:"member_since"`
	ActionsBy   []Action               `json:"actions_by" bson:"actions_by"`
	ActionsOn   []Action               `json:"actions_on" bson:"actions_on"`
	Notes       []Note                 `json:"notes" bson:"notes"`
	Stats       map[string]interface{} `json:"stats" bson:"stats"`
	Source      map[string]interface{} `json:"source" bson:"source"` // source document if imported
}

// FindUserByID retrieves an individual user resource
func FindUserByID(id string) (*User, error) {

	user := User{}

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
func FindUserByEmail(email string) (*User, error) {

	user := User{}
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
func CreateUser(user User) (*User, error) {
	manager := getMongoManager(collectionUser)
	defer manager.close()

	//	dbUser, err := FindUserByEmail(user.Email)
	//	if dbUser != nil {
	//		fmt.Printf("User found:", dbUser.Email)
	//		return dbUser, nil
	//	}
	//
	//	uuid, err := util.NewUUID()
	//	if err != nil {
	//		fmt.Printf("Error getting a new UUID: %v\n", err)
	//		log.Fatal(err)
	//	}
	//
	//	user.ID = uuid
	//	user.MemberSince = time.Now()

	// Write the user to mongo
	err1 := manager.collection.Insert(user)
	if err1 != nil {
		return nil, err1
	}

	return &user, nil
}
