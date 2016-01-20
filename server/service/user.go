package service

import (
	"fmt"
	"github.com/coralproject/pillar/server/dto"
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
	manager.Users.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.Source.ID)
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

//append action to user's actions array and update stats
func updateUserOnAction(user *model.User, object *model.Action, manager *MongoManager) {
	actions := append(user.Actions, object.ID)
	if user.Stats[object.Type] == nil {
		user.Stats[object.Type] = 0
	}

	user.Stats[object.Type] = user.Stats[object.Type].(int) + 1
	manager.Comments.Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"actions": actions, "stats": user.Stats}},
	)
}

//update stats on this user for #comments
func updateUserOnComment(user *model.User, manager *MongoManager) {
	if user.Stats == nil {
		user.Stats = make(map[string]interface{})
	}

	if user.Stats[model.StatsComments] == nil {
		user.Stats[model.StatsComments] = 0
	}

	user.Stats[model.StatsComments] = user.Stats[model.StatsComments].(int) + 1
	manager.Users.Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"stats": user.Stats}},
	)
}

func updateUserMetadata(object *dto.Metadata) (interface{}, *AppError) {
	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := model.User{}
	if object.TargetID != "" {
		manager.Users.FindId(object.TargetID).One(&dbEntity)
	} else {
		manager.Users.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	}

	if dbEntity.ID == "" {
		message := fmt.Sprintf("Cannot update metadata for [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	manager.Users.Update(
		bson.M{"_id": dbEntity.ID},
		bson.M{"$set": bson.M{"metadata": object.Metadata}},
	)

	return dbEntity, nil
}
