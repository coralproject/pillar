package service

import (
	"errors"
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/coralproject/pillar/pkg/db"
)

// CreateAction creates a new action resource
func CreateAction(context *AppContext) (*model.Action, *AppError) {

	db := context.DB
	object := context.Input.(model.Action)

	dbEntity := model.Action{}
	//return, if exists
	db.Actions.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("Action exists with source ID [%s]\n", object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//find & return if one exist with the same source.id
	if object.Source.ID != "" {
		db.Actions.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
		if dbEntity.ID != "" {
			message := fmt.Sprintf("Action exists with source [%s]\n", object.Source.ID)
			return nil, &AppError{nil, message, http.StatusInternalServerError}
		}
	}

	if err := setReferences(db, &object); err != nil {
		message := fmt.Sprintf("Error setting action references [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//do not allow duplicate action from this user on the same target
	db.Actions.Find(bson.M{"user_id": object.UserID, "target_id": object.TargetID,
		"target": object.Target, "type": object.Type}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("Duplicate %s action detected by user [%s] on target [%s: %s]\n",
			object.Type, object.UserID, object.Target, object.TargetID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	if err := db.Actions.Insert(object); err != nil {
		message := fmt.Sprintf("Error creating action [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	if err := updateTargetOnAction(db, &object); err != nil {
		message := fmt.Sprintf("Error updating stats on target [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return &object, nil
}

func setReferences(db *db.MongoDB, object *model.Action) error {

	//set _id
	object.ID = bson.NewObjectId()

	//set user_id
	if object.UserID == "" {
		var user model.User
		db.Users.Find(bson.M{"source.id": object.Source.UserID}).One(&user)
		if user.ID == "" {
			err := errors.New("Cannot find user from source: " + object.Source.UserID)
			return err
		}
		object.UserID = user.ID
	}

	//set target_id
	if object.TargetID == "" {
		if err := setTarget(db, object); err != nil {
			return err
		}
	}

	return nil
}

func setTarget(db *db.MongoDB, object *model.Action) error {

	//find target and set the reference
	switch object.Target {
	case model.Users:
		var user model.User
		db.Users.Find(bson.M{"source.id": object.Source.TargetID}).One(&user)
		if user.ID == "" {
			return errors.New("Cannot find user from source: " + object.Source.TargetID)
		}
		//set the reference
		object.TargetID = user.ID
		break

	case model.Comments:
		var comment model.Comment
		db.Comments.Find(bson.M{"source.id": object.Source.TargetID}).One(&comment)
		if comment.ID == "" {
			return errors.New("Cannot find comment from source: " + object.Source.TargetID)
		}
		//set the reference
		object.TargetID = comment.ID
		break
	}

	return nil
}

func updateTargetOnAction(db *db.MongoDB, object *model.Action) error {

	//find target and set the reference
	switch object.Target {
	case model.Users:
		//update comment with this action
		return updateUserOnAction(db, object)

	case model.Comments:
		//update comment with this action
		return updateCommentOnAction(db, object)
	}

	return nil
}
