package service

import (
	"fmt"
	"github.com/coralproject/pillar/model"
	"gopkg.in/mgo.v2/bson"
)

// CreateUser creates a new user resource
func CreateUser(input model.User) (*model.User, error) {

	// Insert Comment
	manager := getMongoManager()
	defer manager.close()

	dbEntity := model.User{}

	//return, if exists
	manager.users.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with ID [%s]", input.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.users.Find(bson.M{"src_id": input.SourceID}).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with URL [%s]", input.SourceID)
		return &dbEntity, nil
	}

	input.ID = bson.NewObjectId()
	err := manager.users.Insert(input)
	if err != nil {
		return nil, err
	}

	return &dbEntity, nil
}
