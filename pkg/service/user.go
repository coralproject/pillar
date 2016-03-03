package service

import (
	"errors"
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"github.com/coralproject/pillar/pkg/db"
)

// CreateUser creates a new user resource
func CreateUser(context *AppContext) (*model.User, *AppError) {

	db := context.DB
	object := context.Input.(model.User)

	dbEntity := model.User{}

	//return, if exists
	db.Users.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//upsert if entity exists with same source.id
	db.Users.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		object.ID = dbEntity.ID
		_, err := db.Users.UpsertId(dbEntity.ID, &object)
		if err != nil {
			message := fmt.Sprintf("Error updating existing User [%s], %s", object.Source.ID, err)
			return nil, &AppError{err, message, http.StatusInternalServerError}
		}
		return &object, nil
	}

	object.ID = bson.NewObjectId()
	err := db.Users.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating user [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	err = CreateTagTargets(db, object.Tags, &model.TagTarget{Target: model.Users, TargetID: object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return &object, nil
}

//append action to user's actions array and update stats
func updateUserOnAction(db *db.MongoDB, object *model.Action) error {

	var user model.User
	if db.Users.FindId(object.TargetID).One(&user); user.ID == "" {
		return errors.New("Cannot update user stats, invalid user " + object.TargetID.String())
	}

	actions := append(user.Actions, object.ID)
	if user.Stats[object.Type] == nil {
		user.Stats[object.Type] = 0
	}

	user.Stats[object.Type] = user.Stats[object.Type].(int) + 1
	db.Comments.Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"actions": actions, "stats": user.Stats}},
	)

	return nil
}

//update stats on this user for #comments
func updateUserOnComment(db *db.MongoDB, user *model.User) {
	if user.Stats == nil {
		user.Stats = make(map[string]interface{})
	}

	if user.Stats[model.StatsComments] == nil {
		user.Stats[model.StatsComments] = 0
	}

	user.Stats[model.StatsComments] = user.Stats[model.StatsComments].(int) + 1
	db.Users.Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"stats": user.Stats}},
	)
}

// CreateUpdateUser creates/updates a user
func CreateUpdateUser(context *AppContext) (*model.User, *AppError) {

	db := context.DB
	object := context.Input.(model.User)

	var dbEntity *model.User
	//entity not found, return
	db.Users.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID == "" {
		message := fmt.Sprintf("User not found [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	dbEntity.Tags = object.Tags
	if err := db.Users.UpdateId(dbEntity.ID, bson.M{"$set": bson.M{"tags": dbEntity.Tags}}); err != nil {
		message := fmt.Sprintf("Error updating user [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return dbEntity, nil
}

