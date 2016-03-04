package service

import (
	"errors"
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/coralproject/pillar/pkg/db"
)

// ImportAction imports a new action resource
func ImportAction(context *AppContext) (*model.Action, *AppError) {
	db := context.DB
	input := context.Input.(model.Action)

	if err := setReferences(db, &input); err != nil {
		message := fmt.Sprintf("Error setting action references [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//return, if entity exists
	if dbEntity := actionExists(db, &input); dbEntity != nil {
		message := fmt.Sprintf("Action exists [%v]", input)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return doCreateAction(db, &input)
}

// CreateUpdateAction creates/updates an action
func CreateUpdateAction(context *AppContext) (*model.Action, *AppError) {
	input := context.Input.(model.Action)
	if input.ID == "" {
		return createAction(context)
	}

	return updateAction(context)
}

// createAction creates a new action resource
func createAction(context *AppContext) (*model.Action, *AppError) {

	db := context.DB
	input := context.Input.(model.Action)

	//return, if entity exists
	if dbEntity := actionExists(db, &input); dbEntity != nil {
		message := fmt.Sprintf("Action exists [%v]", input)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return doCreateAction(db, &input)
}

// updateAction updates an action resource
func updateAction(context *AppContext) (*model.Action, *AppError) {
	db := context.DB
	input := context.Input.(model.Action)
	var dbEntity model.Action

	//entity not found, return
	db.Actions.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID == "" {
		message := fmt.Sprintf("Action not found [%s]\n", input.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//do we really need to update actions?
	//code here
	return &dbEntity, nil
}

func doCreateAction(db *db.MongoDB, input *model.Action) (*model.Action, *AppError) {
	if err := db.Actions.Insert(input); err != nil {
		message := fmt.Sprintf("Error creating action [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	if err := updateTargetOnAction(db, input); err != nil {
		message := fmt.Sprintf("Error updating stats on target [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return input, nil
}

//finds and returns an action if exists, else nil
func actionExists(db *db.MongoDB, input *model.Action) *model.Action {
	var dbEntity model.Action

	//return, if exists
	db.Actions.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		return &dbEntity
	}

	//do not allow duplicate action from this user on the same target
	db.Actions.Find(bson.M{"user_id": input.UserID, "target_id": input.TargetID,
		"target": input.Target, "type": input.Type}).One(&dbEntity)
	if dbEntity.ID != "" {
		return &dbEntity
	}

	return nil
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
