package service

import (
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
)

// CreateUser creates a new user resource
func CreateUser(object *model.User) (*model.User, *AppError) {

	// get a mongo connection
	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := model.User{}

	//return, if exists
	manager.Users.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//return, if exists
	manager.Users.Find(bson.M{"src_id": object.SourceID}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.SourceID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	object.ID = bson.NewObjectId()
	err := manager.Users.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating user [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return object, nil
}
