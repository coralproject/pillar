package service

import (
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
	//	"fmt"
	//	"reflect"
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
		//fmt.Printf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.users.Find(bson.M{"src_id": object.SourceID}).One(&dbEntity)
	if dbEntity.ID != "" {
		//fmt.Printf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.SourceID)
		return &dbEntity, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.users.Insert(object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}
