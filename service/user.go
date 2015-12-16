package service

import (
	"fmt"
	"github.com/coralproject/pillar/model"
	"gopkg.in/mgo.v2/bson"
)

// CreateUser creates a new user resource
func CreateUser(object model.User) (*model.User, error) {

	// Insert Comment
	manager := getMongoManager()
	defer manager.close()

	dbEntity := model.User{}

	//return, if exists
	manager.users.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with ID [%s]", object.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.users.Find(bson.M{"src_id": object.SourceID}).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with URL [%s]", object.SourceID)
		return &dbEntity, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.users.Insert(object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}
