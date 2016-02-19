package service

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"github.com/coralproject/pillar/pkg/model"
	"errors"
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

	//upsert if entity exists with same source.id
	manager.Users.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		object.ID = dbEntity.ID
		_, err := manager.Users.UpsertId(dbEntity.ID, object)
		if err != nil {
			message := fmt.Sprintf("Error updating existing User [%s], %s", object.Source.ID, err)
			return nil, &AppError{err, message, http.StatusInternalServerError}
		}
		return object, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.Users.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating user [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	err = CreateTagTargets(manager, object.Tags, &model.TagTarget{Target:model.Users, TargetID:object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return object, nil
}

//append action to user's actions array and update stats
func updateUserOnAction(object *model.Action, manager *MongoManager) error {

	var user model.User
	if manager.Users.FindId(object.TargetID).One(&user); user.ID == "" {
		return errors.New("Cannot update user stats, invalid user " + object.TargetID.String())
	}

	actions := append(user.Actions, object.ID)
	if user.Stats[object.Type] == nil {
		user.Stats[object.Type] = 0
	}

	user.Stats[object.Type] = user.Stats[object.Type].(int) + 1
	manager.Comments.Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"actions": actions, "stats": user.Stats}},
	)

	return nil
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