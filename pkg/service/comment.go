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

type reference struct {
	asset *model.Asset
	user  *model.User
}
var ref reference

// CreateComment creates a new comment resource
func CreateComment(context *AppContext) (*model.Comment, *AppError) {

	db := context.DB
	object := context.Input.(model.Comment)

	// Insert Comment
	var dbEntity model.Comment

	//return, if exists
	if db.Comments.FindId(object.ID).One(&dbEntity); dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	object.ID = bson.NewObjectId()
	if err := setCommentReferences(db, &object); err != nil {
		message := fmt.Sprintf("Error setting comment references [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//upsert if entity exists with same source.id
	db.Comments.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		object.ID = dbEntity.ID
		_, err := db.Comments.UpsertId(dbEntity.ID, object)
		if err != nil {
			message := fmt.Sprintf("Error updating existing Comment [%s], %s", object.Source.ID, err)
			return nil, &AppError{err, message, http.StatusInternalServerError}
		}
		return &object, nil
	}

	//Insert new comment
	if err := db.Comments.Insert(object); err != nil {
		message := fmt.Sprintf("Error creating comments [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	updateUserOnComment(db, ref.user)
	updateAssetOnComment(db, ref.asset)

	err := CreateTagTargets(db, object.Tags, &model.TagTarget{Target: model.Comments, TargetID: object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return &object, nil
}

func setCommentReferences(db *db.MongoDB, object *model.Comment) error {
	//find asset and add the reference to it
	var asset model.Asset
	db.Assets.Find(bson.M{"source.id": object.Source.AssetID}).One(&asset)
	if asset.ID == "" {
		db.Assets.Find(bson.M{"url": object.Source.AssetID}).One(&asset)
	}
	if asset.ID == "" {
		return errors.New("Cannot find asset from source: " + object.Source.AssetID)
	}
	object.AssetID = asset.ID
	ref.asset = &asset

	//find user and add the reference to it
	var user model.User
	db.Users.Find(bson.M{"source.id": object.Source.UserID}).One(&user)
	if user.ID == "" {
		err := errors.New("Cannot find user from source: " + object.Source.UserID)
		return err
	}
	object.UserID = user.ID
	ref.user = &user

	//find parent and add the reference to it
	if object.Source.ID != object.Source.ParentID {
		var parent model.Comment
		db.Comments.Find(bson.M{"source.parent_id": object.Source.ParentID}).One(&parent)
		if parent.ID != "" {
			object.ParentID = parent.ID
			//add this as a child for the parent comment
			children := append(parent.Children, object.ID)
			db.Comments.Update(bson.M{"_id": parent.ID},
				bson.M{"$set": bson.M{"children": children}})
		}
	}

	return nil
}

//append action to comment's actions array and update stats
func updateCommentOnAction(db *db.MongoDB, object *model.Action) error {

	var comment model.Comment
	if db.Comments.FindId(object.TargetID).One(&comment); comment.ID == "" {
		return errors.New("Cannot update comment stats, invalid comment " + object.TargetID.String())
	}

	actions := append(comment.Actions, object.ID)

	if comment.Stats == nil {
		comment.Stats = make(map[string]interface{})
	}

	if comment.Stats[object.Type] == nil {
		comment.Stats[object.Type] = 0
	}

	comment.Stats[object.Type] = comment.Stats[object.Type].(int) + 1
	db.Comments.Update(
		bson.M{"_id": comment.ID},
		bson.M{"$set": bson.M{"actions": actions, "stats": comment.Stats}},
	)

	return nil
}

// CreateUpdateComment creates/updates a comment
func CreateUpdateComment(context *AppContext) (*model.Comment, *AppError) {
	db := context.DB
	object := context.Input.(model.Comment)

	var dbEntity *model.Comment
	//entity not found, return
	db.Comments.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID == "" {
		message := fmt.Sprintf("Comment not found [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	dbEntity.Tags = object.Tags
	if err := db.Comments.UpdateId(dbEntity.ID, bson.M{"$set": bson.M{"tags": dbEntity.Tags}}); err != nil {
		message := fmt.Sprintf("Error updating comment [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return dbEntity, nil
}
