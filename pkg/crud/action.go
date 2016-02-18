package crud

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// CreateAction creates a new action resource
func CreateAction(object *Action) (*Action, *AppError) {

	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := Action{}

	//return, if exists
	manager.Actions.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("Action exists with source ID [%s]\n", object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//find & return if one exist with the same source.id
	if object.Source.ID != "" {
		manager.Actions.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
		if dbEntity.ID != "" {
			message := fmt.Sprintf("Action exists with source [%s]\n", object.Source.ID)
			return nil, &AppError{nil, message, http.StatusInternalServerError}
		}
	}

	//do not allow duplicate action from this user
	user := findUser(object, manager)
	manager.Actions.Find(bson.M{"user_id": user.ID, "target": object.Target, "type": object.Type}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("Action exists with user [%s], target [%s] and type [%s]\n",
			user.ID, object.Target, object.Type)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	object.ID = bson.NewObjectId()
	if err := setActionReferences(object, user, manager); err != nil {
		message := fmt.Sprintf("Error setting action references [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	err := manager.Actions.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating action [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return object, nil
}

func setActionReferences(object *Action, user *User, manager *MongoManager) error {

	object.UserID = user.ID

	//find target and set the reference
	switch object.Target {
	case Users:
		var user User
		manager.Users.Find(bson.M{"source.id": object.Source.TargetID}).One(&user)
		if user.ID == "" {
			err := errors.New("Cannot find user from source: " + object.Source.TargetID)
			return err
		}
		//set the reference
		object.TargetID = user.ID

		//update comment with this action
		updateUserOnAction(&user, object, manager)
		break

	case Comments:
		var comment Comment
		manager.Comments.Find(bson.M{"source.id": object.Source.TargetID}).One(&comment)
		if comment.ID == "" {
			err := errors.New("Cannot find comment from source: " + object.Source.TargetID)
			return err
		}
		//set the reference
		object.TargetID = comment.ID

		//update comment with this action
		updateCommentOnAction(&comment, object, manager)
		break
	}

	return nil
}

func findUser(object *Action, manager *MongoManager) *User {
	var user User

	if object.UserID != "" {
		manager.Users.FindId(object.UserID).One(&user)
		if user.ID != "" {
			return &user;
		}
	}

	manager.Users.Find(bson.M{"source.id": object.Source.UserID}).One(&user)
	if user.ID != "" {
		return &user
	}

	return nil
}
