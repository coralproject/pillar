package service

import (
	"github.com/coralproject/pillar/server/log"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

// CreateUser creates a new user resource
func CreateUser(object model.User) (*model.User, error) {

	// get a mongo connection
	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := model.User{}

	//return, if exists
	manager.Users.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		log.Logger.Printf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.Users.Find(bson.M{"src_id": object.SourceID}).One(&dbEntity)
	if dbEntity.ID != "" {
		log.Logger.Printf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.SourceID)
		return &dbEntity, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.Users.Insert(object)
	if err != nil {
		log.Logger.Printf("Error creating user [%s]", err);
		return nil, err
	}

	return &object, nil
}
