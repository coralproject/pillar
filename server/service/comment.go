package service

import (
	"errors"
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

// CreateComment creates a new comment resource
func CreateComment(object model.Comment) (*model.Comment, error) {

	// Insert Comment
	manager := getMongoManager()
	defer manager.close()

	var dbEntity model.Comment

	//return, if exists
	if manager.comments.FindId(object.ID).One(&dbEntity); dbEntity.ID != "" {
		fmt.Printf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return &dbEntity, nil
	}

	//find & return if one exist with the same source.id
	manager.comments.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.Source.ID)
		return &dbEntity, nil
	}

	//fix all references with ObjectId
	object.ID = bson.NewObjectId()
	if err := fixReferences(&object, manager); err != nil {
		return nil, err
	}

	if err := manager.comments.Insert(object); err != nil {
		return nil, err
	}

	return &object, nil
}

func fixReferences(object *model.Comment, manager *mongoManager) error {
	//find asset and add the reference to it
	var asset model.Asset
	manager.assets.Find(bson.M{"src_id": object.Source.AssetID}).One(&asset)
	if asset.ID == "" {
		manager.assets.Find(bson.M{"url": object.Source.AssetID}).One(&asset)
	}
	if asset.ID == "" {
		return errors.New("Cannot find asset from source: " + object.Source.AssetID)
	}
	object.AssetID = asset.ID

	//find user and add the reference to it
	var user model.User
	manager.users.Find(bson.M{"src_id": object.Source.UserID}).One(&user)
	if user.ID == "" {
		return errors.New("Cannot find user from source: " + object.Source.UserID)
	}
	object.UserID = user.ID

	//find parent and add the reference to it
	if object.Source.ID != object.Source.ParentID {
		var parent model.Comment
		manager.comments.Find(bson.M{"source.parent_id": object.Source.ParentID}).One(&parent)
		if parent.ID != "" {
			object.ParentID = parent.ID
			//add this as a child for the parent comment
			//parent.Children = make([]bson.ObjectId, 10)
			children := append(parent.Children, object.ID)
			manager.comments.Update(bson.M{"_id": parent.ID},
				bson.M{ "$set" : bson.M{"children" : children} } )
		}
	}
	return nil
}
