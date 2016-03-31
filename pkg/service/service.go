package service

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
	"github.com/coralproject/pillar/pkg/model/event"
)

func GetPayload(context *web.AppContext, object interface{}) interface{} {

	switch object.(type) {
	case *model.Comment:
		comment := object.(*model.Comment)
		var user model.User
		context.DB.Users.FindId(comment.UserID).One(&user);
		var asset model.Asset
		context.DB.Assets.FindId(comment.AssetID).One(&asset);
		return event.PayloadComment{*comment, asset, user}
	default:
		return nil
	}
}

// UpdateMetadata updates metadata for an entity
func UpdateMetadata(context *web.AppContext) (interface{}, *web.AppError) {

	db := context.DB
	var input model.Metadata
	json.NewDecoder(context.Body).Decode(&input)

	collection := db.Session.DB("").C(input.Target)
	var dbEntity bson.M
	collection.FindId(input.TargetID).One(&dbEntity)
	if len(dbEntity) == 0 {
		collection.Find(bson.M{"source.id": input.Source.ID}).One(&dbEntity)
	}

	if len(dbEntity) == 0 {
		message := fmt.Sprintf("Cannot update metadata for [%+v]\n", input)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	collection.Update(
		bson.M{"_id": dbEntity["_id"]},
		bson.M{"$set": bson.M{"metadata": input.Metadata}},
	)

	return dbEntity, nil
}

// CreateIndex creates indexes to various entities
func CreateIndex(context *web.AppContext) *web.AppError {

	db := context.DB
	var input model.Index
	json.NewDecoder(context.Body).Decode(&input)

	err := db.Session.DB("").C(input.Target).EnsureIndex(input.Index)
	if err != nil {
		message := fmt.Sprintf("Error creating index [%+v]", input)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}

// CreateUserAction inserts an activity by the user
func CreateUserAction(context *web.AppContext) *web.AppError {

	db := context.DB
	var input model.CayUserAction
	json.NewDecoder(context.Body).Decode(&input)

	input.ID = bson.NewObjectId()
	input.Date = time.Now()
	if input.Release == "" {
		input.Release = "0.1.0"
	}
	err := db.CayUserActions.Insert(input)
	if err != nil {
		message := fmt.Sprintf("Error creating user-action [%s]", err)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
