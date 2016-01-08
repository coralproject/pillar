package service

import (
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"errors"
)

// CreateAction creates a new action resource
func CreateAction(object model.Action) (*model.Action, *AppError) {

	// Insert Comment
	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := model.Action{}

	//return, if exists
	manager.Actions.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//find & return if one exist with the same source.id
	if object.Source.ID != "" {
		manager.Actions.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
		if dbEntity.ID != "" {
			message := fmt.Sprintf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.Source.ID)
			return nil, &AppError{nil, message, http.StatusInternalServerError}
		}
	}

	object.ID = bson.NewObjectId()
	if err := setActionReferences(&object, manager); err != nil {
		message := fmt.Sprintf("Error setting action references [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	err := manager.Actions.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating action [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return &object, nil
}

func setActionReferences(object *model.Action, manager *MongoManager) error {

	//find user and set the reference
	var user model.User
	manager.Users.Find(bson.M{"src_id": object.Source.UserID}).One(&user)
	if user.ID == "" {
		err := errors.New("Cannot find user from source: " + object.Source.UserID)
		return err
	}
	object.UserID = user.ID

	//find target and set the reference
	switch object.TargetType {
		case model.TargetTypeUser:
			var user model.User
			manager.Users.Find(bson.M{"src_id": object.Source.TargetID}).One(&user)
			if user.ID == "" {
				err := errors.New("Cannot find user from source: " + object.Source.TargetID)
				return err
			}
			//set the reference
			object.TargetID = user.ID

			//also append this action to user's actions array
			actions := append(user.Actions, object.ID)
			manager.Users.Update(bson.M{"_id": user.ID},
				bson.M{"$set": bson.M{"actions": actions}})
			break

		case model.TargetTypeComment:
			var comment model.Comment
			manager.Comments.Find(bson.M{"source.id": object.Source.TargetID}).One(&comment)
			if comment.ID == "" {
				err := errors.New("Cannot find comment from source: " + object.Source.TargetID)
				return err
			}
			//set the reference
			object.TargetID = comment.ID

			//also append this action to comment's actions array
			actions := append(comment.Actions, object.ID)
			manager.Comments.Update(bson.M{"_id": comment.ID},
				bson.M{"$set": bson.M{"actions": actions}})
			break
	}

	return nil
}

