package service

import (
	"errors"
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
)

// CreateComment creates a new comment resource
func CreateComment(object model.Comment) (*model.Comment, error) {

	// Insert Comment
	manager := getMongoManager()
	defer manager.close()

	var dbEntity model.Comment

	//return, if exists
	manager.comments.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with ID [%s]", object.ID)
		return &dbEntity, nil
	}

	//find & return if one exist with the same source.id
	one := findOne(manager.comments, bson.M{"source.id": object.Source.ID})
	if one != nil {
		fmt.Printf("Entity exists with source [%s]", object.Source.ID)
		comment := one.(model.Comment)
		return &comment, nil
	}

	//fix all references with ObjectId
	object.ID = bson.NewObjectId()
	err := fixReferences(&object, manager)
	if err != nil {
		return nil, err
	}

	err = manager.comments.Insert(object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

func fixReferences(object *model.Comment, manager *mongoManager) error {
	//find asset and add the reference to it
	var asset model.Asset
	manager.assets.Find(bson.M{"src_id": object.Source.AssetID}).One(&asset)
	if asset.ID == "" {
		//asset = findOne(manager.assets, bson.M{"url": object.Source.AssetID})
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
	var parent model.Comment
	manager.comments.Find(bson.M{"source.id": object.Source.UserID}).One(&parent)
	if parent.ID != "" {
		object.ParentID = parent.ID
	}

	return nil
}
