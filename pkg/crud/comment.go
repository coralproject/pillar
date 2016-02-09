package crud

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

var commenter User

// CreateComment creates a new comment resource
func CreateComment(object *Comment) (*Comment, *AppError) {

	// Insert Comment
	manager := GetMongoManager()
	defer manager.Close()

	var dbEntity Comment

	//return, if exists
	if manager.Comments.FindId(object.ID).One(&dbEntity); dbEntity.ID != "" {
		message := fmt.Sprintf("Comment exists with ID [%s]\n", object.ID)
		return nil, &AppError{nil, message, http.StatusConflict}
	}

	//find & return if one exist with the same source.id
	manager.Comments.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("Comment exists with Source.ID [%s]\n", object.Source.ID)
		return nil, &AppError{nil, message, http.StatusConflict}
	}

	//fix all references with ObjectId
	object.ID = bson.NewObjectId()
	if err := setCommentReferences(object, manager); err != nil {
		message := fmt.Sprintf("Error setting comment references [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//	//add actions
	//	if err := setActions(&object, manager); err != nil {
	//		message := fmt.Sprintf("Error setting comment actions [%s]", err)
	//		return nil, &AppError{nil, message, http.StatusInternalServerError}
	//	}

	//	//add notes
	//	if err := setNotes(&object, manager); err != nil {
	//		message := fmt.Sprintf("Error setting comment notes [%s]", err)
	//		return nil, &AppError{nil, message, http.StatusInternalServerError}
	//	}

	if err := manager.Comments.Insert(object); err != nil {
		message := fmt.Sprintf("Error creating comments [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	updateUserOnComment(&commenter, manager)

	err := CreateTagTargets(manager, object.Tags, &TagTarget{Target:Comments, TargetID:object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return object, nil
}

func setCommentReferences(object *Comment, manager *MongoManager) error {
	//find asset and add the reference to it
	var asset Asset
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
		var parent Comment
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
func updateCommentOnAction(comment *Comment, object *Action, manager *MongoManager) {
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
