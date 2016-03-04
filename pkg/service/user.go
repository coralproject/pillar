package service

import (
	"errors"
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// ImportUser imports a new user resource
func ImportUser(context *AppContext) (*model.User, *AppError) {

	db := context.DB
	input := context.Input.(model.User)
	var dbEntity model.User

	//upsert if entity exists with same source.id
	db.Users.Find(bson.M{"source.id": input.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		input.ID = dbEntity.ID
		if _, err := db.Users.UpsertId(dbEntity.ID, &input); err != nil {
			message := fmt.Sprintf("Error updating existing User [%s]", input.Source.ID)
			return nil, &AppError{err, message, http.StatusInternalServerError}
		}
		return &input, nil
	}

	return doCreateUser(db, &input)
}

// CreateUpdateUser creates/updates a user resource
func CreateUpdateUser(context *AppContext) (*model.User, *AppError) {
	input := context.Input.(model.User)
	if input.ID == "" {
		return createUser(context)
	}

	return updateUser(context)
}

// createUser creates a new user resource
func createUser(context *AppContext) (*model.User, *AppError) {

	db := context.DB
	input := context.Input.(model.User)
	var dbEntity model.User

	//return, if exists
	db.Users.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("User exists with ID [%s]\n", input.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return doCreateUser(db, &input)
}

// updateUser updates a user
func updateUser(context *AppContext) (*model.User, *AppError) {

	db := context.DB
	input := context.Input.(model.User)

	var dbEntity *model.User
	//entity not found, return
	db.Users.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID == "" {
		message := fmt.Sprintf("User not found [%+v]\n", input)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	dbEntity.Tags = input.Tags
	if err := db.Users.UpdateId(dbEntity.ID, bson.M{"$set": bson.M{"tags": dbEntity.Tags}}); err != nil {
		message := fmt.Sprintf("Error updating user [%+v]\n", input)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return dbEntity, nil
}

//inserts a new user to the db and any related post-processing
func doCreateUser(db *db.MongoDB, input *model.User) (*model.User, *AppError) {
	input.ID = bson.NewObjectId()
	if err := db.Users.Insert(input); err != nil {
		message := fmt.Sprintf("Error creating user [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	tt := model.TagTarget{Target: model.Users, TargetID: input.ID}
	if err := CreateTagTargets(db, input.Tags, &tt); err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return input, nil
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
