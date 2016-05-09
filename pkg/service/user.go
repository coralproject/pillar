package service

import (
	"errors"
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// ImportUser imports a new user resource
func ImportUser(context *web.AppContext) (*model.User, *web.AppError) {

	var input model.User
	if err := UnmarshallAndValidate(context, &input); err != nil {
		return nil, err
	}

	var dbEntity model.User
	//upsert if entity exists with same source.id
	context.MDB.DB.C(model.Users).Find(bson.M{"source.id": input.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		input.ID = dbEntity.ID
		if _, err := context.MDB.DB.C(model.Users).UpsertId(dbEntity.ID, &input); err != nil {
			message := fmt.Sprintf("Error updating existing User [%s]", input.Source.ID)
			return nil, &web.AppError{err, message, http.StatusInternalServerError}
		}
		return &input, nil
	}

	return doCreateUser(context.MDB, &input)
}

// CreateUpdateUser creates/updates a user resource
func CreateUpdateUser(context *web.AppContext) (*model.User, *web.AppError) {

	var input model.User
	if err := UnmarshallAndValidate(context, &input); err != nil {
		return nil, err
	}

	if input.ID == "" {
		return createUser(context.MDB, &input)
	}

	return updateUser(context.MDB, &input)
}

// createUser creates a new user resource
func createUser(db *db.MongoDB, input *model.User) (*model.User, *web.AppError) {

	var dbEntity model.User
	//return, if exists
	db.DB.C(model.Users).FindId(input.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("User exists with ID [%s]\n", input.ID)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	return doCreateUser(db, input)
}

// updateUser updates a user
func updateUser(db *db.MongoDB, input *model.User) (*model.User, *web.AppError) {

	var dbEntity *model.User
	//entity not found, return
	db.DB.C(model.Users).FindId(input.ID).One(&dbEntity)
	if dbEntity.ID == "" {
		message := fmt.Sprintf("User not found [%+v]\n", input)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	dbEntity.Tags = input.Tags
	if err := db.DB.C(model.Users).UpdateId(dbEntity.ID, bson.M{"$set": bson.M{"tags": dbEntity.Tags}}); err != nil {
		message := fmt.Sprintf("Error updating user [%+v]\n", input)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	return dbEntity, nil
}

//inserts a new user to the db and any related post-processing
func doCreateUser(db *db.MongoDB, input *model.User) (*model.User, *web.AppError) {
	input.ID = bson.NewObjectId()
	if err := db.DB.C(model.Users).Insert(input); err != nil {
		message := fmt.Sprintf("Error creating user [%s]", err)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	tt := model.TagTarget{Target: model.Users, TargetID: input.ID}
	if err := CreateTagTargets(db, input.Tags, &tt); err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	return input, nil
}

//append action to user's actions array and update stats
func updateUserOnAction(db *db.MongoDB, object *model.Action) error {

	var user model.User
	if db.DB.C(model.Users).FindId(object.TargetID).One(&user); user.ID == "" {
		return errors.New("Cannot update user stats, invalid user " + object.TargetID.String())
	}

	actions := append(user.Actions, *object)
	if user.Stats[object.Type] == nil {
		user.Stats[object.Type] = 0
	}

	user.Stats[object.Type] = user.Stats[object.Type].(int) + 1
	db.DB.C(model.Users).Update(
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
	db.DB.C(model.Users).Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"stats": user.Stats}},
	)
}
