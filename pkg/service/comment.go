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

type reference struct {
	asset *model.Asset
	user  *model.User
}

var ref reference

// ImportComment imports a new comment resource
func ImportComment(context *web.AppContext) (*model.Comment, *web.AppError) {

	var input model.Comment
	if err := UnmarshallAndValidate(context, &input); err != nil {
		return nil, err
	}

	var dbEntity model.Comment
	// Find/Set comment references
	if err := setCommentReferences(context.MDB, &input); err != nil {
		message := fmt.Sprintf("Error setting comment references [%s]", err)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//upsert if entity exists with same source.id
	context.MDB.DB.C(model.Comments).Find(bson.M{"source.id": input.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		input.ID = dbEntity.ID
		if _, err := context.MDB.DB.C(model.Comments).UpsertId(dbEntity.ID, input); err != nil {
			message := fmt.Sprintf("Error updating existing Comment [%s]", input.Source.ID)
			return nil, &web.AppError{err, message, http.StatusInternalServerError}
		}
		return &input, nil
	}

	return doCreateComment(context.MDB, &input)
}

// CreateUpdateComment creates/updates a comment resource
func CreateUpdateComment(context *web.AppContext) (*model.Comment, *web.AppError) {
	var input model.Comment
	if err := UnmarshallAndValidate(context, &input); err != nil {
		return nil, err
	}

	if input.ID == "" {
		return createComment(context.MDB, &input)
	}

	return updateComment(context.MDB, &input)
}

// CreateComment creates a new comment resource
func createComment(db *db.MongoDB, input *model.Comment) (*model.Comment, *web.AppError) {

	var dbEntity model.Comment
	//return, if exists
	db.DB.C(model.Comments).FindId(input.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("Comment exists with ID [%s]\n", input.ID)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	var asset model.Asset
	db.DB.C(model.Assets).FindId(input.AssetID).One(&asset)
	if asset.ID == "" {
		message := fmt.Sprintf("Cannot create Comment, Asset not found [%s]\n", input.AssetID)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}
	ref.asset = &asset

	var user model.User
	db.DB.C(model.Users).FindId(input.UserID).One(&user)
	if user.ID == "" {
		message := fmt.Sprintf("Cannot create Comment, User not found [%s]\n", input.UserID)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}
	ref.user = &user

	input.ID = bson.NewObjectId()
	return doCreateComment(db, input)
}

// updateComment updates a comment
func updateComment(db *db.MongoDB, input *model.Comment) (*model.Comment, *web.AppError) {

	var dbEntity *model.Comment
	//entity not found, return
	db.DB.C(model.Comments).FindId(input.ID).One(&dbEntity)
	if dbEntity.ID == "" {
		message := fmt.Sprintf("Comment not found [%+v]\n", input)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	dbEntity.Tags = input.Tags
	if err := db.DB.C(model.Comments).UpdateId(dbEntity.ID, bson.M{"$set": bson.M{"tags": dbEntity.Tags}}); err != nil {
		message := fmt.Sprintf("Error updating comment [%+v]\n", input)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	return dbEntity, nil
}

//inserts a new comment to the db and any related post-processing
func doCreateComment(db *db.MongoDB, input *model.Comment) (*model.Comment, *web.AppError) {
	if err := db.DB.C(model.Comments).Insert(input); err != nil {
		message := fmt.Sprintf("Error creating comments [%s]", err)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	updateUserOnComment(db, ref.user)
	updateAssetOnComment(db, ref.asset)

	tt := &model.TagTarget{Target: model.Comments, TargetID: input.ID}
	if err := CreateTagTargets(db, input.Tags, tt); err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	return input, nil
}

func setCommentReferences(db *db.MongoDB, input *model.Comment) error {
	input.ID = bson.NewObjectId()

	//find asset and add the reference to it
	var asset model.Asset
	db.DB.C(model.Assets).Find(bson.M{"source.id": input.Source.AssetID}).One(&asset)
	if asset.ID == "" {
		db.DB.C(model.Assets).Find(bson.M{"url": input.Source.AssetID}).One(&asset)
	}
	if asset.ID == "" {
		return errors.New("Cannot find asset from source: " + input.Source.AssetID)
	}
	input.AssetID = asset.ID
	ref.asset = &asset

	//find user and add the reference to it
	var user model.User
	db.DB.C(model.Users).Find(bson.M{"source.id": input.Source.UserID}).One(&user)
	if user.ID == "" {
		err := errors.New("Cannot find user from source: " + input.Source.UserID)
		return err
	}
	input.UserID = user.ID
	ref.user = &user

	//find parent and add the reference to it
	if input.Source.ID != input.Source.ParentID {
		var parent model.Comment
		db.DB.C(model.Comments).Find(bson.M{"source.parent_id": input.Source.ParentID}).One(&parent)
		if parent.ID != "" {
			//prepend immediate parent, followed by parents from hierarchy
			parents := append([]bson.ObjectId{parent.ID}, parent.Parents...)
			input.Parents = parents
			if len(parents) > 0 {
				input.RootID = parents[len(parents)-1]
			}

			//add this as a child for the parent comment
			children := append(parent.Children, input.ID)
			db.DB.C(model.Comments).Update(bson.M{"_id": parent.ID},
				bson.M{"$set": bson.M{"children": children}})
		}
	}

	return nil
}

//append action to comment's actions array and update stats
func updateCommentOnAction(db *db.MongoDB, object *model.Action) error {

	var comment model.Comment
	if db.DB.C(model.Comments).FindId(object.TargetID).One(&comment); comment.ID == "" {
		return errors.New("Cannot update comment stats, invalid comment " + object.TargetID.String())
	}

	actions := append(comment.Actions, *object)

	if comment.Stats == nil {
		comment.Stats = make(map[string]interface{})
	}

	if comment.Stats[object.Type] == nil {
		comment.Stats[object.Type] = 0
	}

	comment.Stats[object.Type] = comment.Stats[object.Type].(int) + 1
	db.DB.C(model.Comments).Update(
		bson.M{"_id": comment.ID},
		bson.M{"$set": bson.M{"actions": actions, "stats": comment.Stats}},
	)

	return nil
}

func getPayloadComment(context *web.AppContext, object interface{}) interface{} {
	comment := object.(*model.Comment)

	var user model.User
	context.MDB.DB.C(model.Users).FindId(comment.UserID).One(&user)
	var asset model.Asset
	context.MDB.DB.C(model.Assets).FindId(comment.AssetID).One(&asset)

	return model.PayloadComment{*comment, asset, user}
}
