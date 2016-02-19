package service

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"github.com/coralproject/pillar/pkg/model"
)

var commenter model.User

// CreateComment creates a new comment resource
func CreateComment(object *model.Comment) (*model.Comment, *AppError) {

	// Insert Comment
	manager := GetMongoManager()
	defer manager.Close()

	var dbEntity model.Comment

	//return, if exists
	if manager.Comments.FindId(object.ID).One(&dbEntity); dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	if err := setCommentReferences(object, manager); err != nil {
		message := fmt.Sprintf("Error setting comment references [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//upsert if entity exists with same source.id
	manager.Comments.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		object.ID = dbEntity.ID
		_, err := manager.Users.UpsertId(dbEntity.ID, object)
		if err != nil {
			message := fmt.Sprintf("Error updating existing Comment [%s], %s", object.Source.ID, err)
			return nil, &AppError{err, message, http.StatusInternalServerError}
		}
		return object, nil
	}

	//Insert new comment
	object.ID = bson.NewObjectId()
	if err := manager.Comments.Insert(object); err != nil {
		message := fmt.Sprintf("Error creating comments [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	updateUserOnComment(&commenter, manager)

	err := CreateTagTargets(manager, object.Tags, &model.TagTarget{Target:model.Comments, TargetID:object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return object, nil
}

func setCommentReferences(object *model.Comment, manager *MongoManager) error {
	//find asset and add the reference to it
	var asset model.Asset
	manager.Assets.Find(bson.M{"source.id": object.Source.AssetID}).One(&asset)
	if asset.ID == "" {
		manager.Assets.Find(bson.M{"url": object.Source.AssetID}).One(&asset)
	}
	if asset.ID == "" {
		return errors.New("Cannot find asset from source: " + object.Source.AssetID)
	}
	object.AssetID = asset.ID

	//find user and add the reference to it
	manager.Users.Find(bson.M{"source.id": object.Source.UserID}).One(&commenter)
	if commenter.ID == "" {
		err := errors.New("Cannot find user from source: " + object.Source.UserID)
		return err
	}
	object.UserID = commenter.ID

	//find parent and add the reference to it
	if object.Source.ID != object.Source.ParentID {
		var parent model.Comment
		manager.Comments.Find(bson.M{"source.parent_id": object.Source.ParentID}).One(&parent)
		if parent.ID != "" {
			object.ParentID = parent.ID
			//add this as a child for the parent comment
			children := append(parent.Children, object.ID)
			manager.Comments.Update(bson.M{"_id": parent.ID},
				bson.M{"$set": bson.M{"children": children}})
		}
	}

	return nil
}

//append action to comment's actions array and update stats
func updateCommentOnAction(object *model.Action, manager *MongoManager) error {

	var comment model.Comment
	if manager.Comments.FindId(object.TargetID).One(&comment); comment.ID == "" {
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
	manager.Comments.Update(
		bson.M{"_id": comment.ID},
		bson.M{"$set": bson.M{"actions": actions, "stats": comment.Stats}},
	)

	return nil
}

//func setActions(object *Comment, manager *MongoManager) error {
//	var user User
//	var invalidUsers []string
//
//	for i := 0; i < len(object.Actions); i++ {
//		one := &object.Actions[i]
//		manager.Users.Find(bson.M{"source.id": one.Source.UserID}).One(&user)
//		if user.ID == "" {
//			invalidUsers = append(invalidUsers, one.Source.UserID)
//			continue
//		}
//
//		one.UserID = user.ID
//	}
//
//	if len(invalidUsers) > 0 {
//		return errors.New("Error setting comment actions - Cannot find users")
//	}
//
//	return nil
//}

//func setNotes(object *Comment, manager *MongoManager) error {
//	var user User
//	var invalidUsers []string
//
//	for i := 0; i < len(object.Notes); i++ {
//		one := &object.Notes[i]
//		manager.Users.Find(bson.M{"source.id": one.SourceUserID}).One(&user)
//		if user.ID == "" {
//			invalidUsers = append(invalidUsers, one.SourceUserID)
//			continue
//		}
//
//		one.UserID = user.ID
//	}
//
//	if len(invalidUsers) > 0 {
//		return errors.New("Error setting comment notes - Cannot find users")
//	}
//
//	return nil
//}
