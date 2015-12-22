package service

import (
	"reflect"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

// CreateUser creates a new user resource
func CreateUser(object model.User) (*model.User, error) {

	// get a mongo connection
	manager := getMongoManager()
	defer manager.close()

	dbEntity := model.User{}

	//return, if exists
	manager.users.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprint("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		fmt.Printf(message)
		return &dbEntity, nil
	}

	//return, if exists
	manager.users.Find(bson.M{"src_id": object.SourceID}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprint("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.SourceID)
		fmt.Printf(message)
		return &dbEntity, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.users.Insert(object)
	if err != nil {
		log.Error("service", "CreateUser", err, "Inserting users")
		return nil, err
	}

	return &object, nil
}
