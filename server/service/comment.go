package service

import (
	"errors"
	"reflect"
	"github.com/coralproject/pillar/server/log"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
)

// CreateComment creates a new comment resource
func CreateComment(object model.Comment) (*model.Comment, error) {

	// Insert Comment
	manager := GetMongoManager()
	defer manager.Close()

	var dbEntity model.Comment

	//return, if exists
	if manager.Comments.FindId(object.ID).One(&dbEntity); dbEntity.ID != "" {
		log.Logger.Printf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return &dbEntity, nil
	}

	//find & return if one exist with the same source.id
	manager.Comments.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		log.Logger.Printf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.Source.ID)
		return &dbEntity, nil
	}

	//fix all references with ObjectId
	object.ID = bson.NewObjectId()
	if err := setReferences(&object, manager); err != nil {
		log.Logger.Printf("Error setting references [%s]", err);
		return nil, err
	}

	//add actions
	if err := setActions(&object, manager); err != nil {
		log.Logger.Printf("Error setting actions [%s]", err);
		return nil, err
	}

	//add notes
	if err := setNotes(&object, manager); err != nil {
		log.Logger.Printf("Error setting notes [%s]", err);
		return nil, err
	}

	//fmt.Printf("Comment: %+v\n\n", object)
	if err := manager.Comments.Insert(object); err != nil {
		log.Logger.Printf("Error creating comments [%s]", err);
		return nil, err
	}

	return &object, nil
}

func setReferences(object *model.Comment, manager *MongoManager) error {
	//find asset and add the reference to it
	var asset model.Asset
	manager.Assets.Find(bson.M{"src_id": object.Source.AssetID}).One(&asset)
	if asset.ID == "" {
		manager.Assets.Find(bson.M{"url": object.Source.AssetID}).One(&asset)
	}
	if asset.ID == "" {
		return errors.New("Cannot find asset from source: " + object.Source.AssetID)
	}
	object.AssetID = asset.ID

	//find user and add the reference to it
	var user model.User
	manager.Users.Find(bson.M{"src_id": object.Source.UserID}).One(&user)
	if user.ID == "" {
		err := errors.New("Cannot find user from source: " + object.Source.UserID)
		return err
	}
	object.UserID = user.ID

	//find parent and add the reference to it
	if object.Source.ID != object.Source.ParentID {
		var parent model.Comment
		manager.Comments.Find(bson.M{"source.parent_id": object.Source.ParentID}).One(&parent)
		if parent.ID != "" {
			object.ParentID = parent.ID
			//add this as a child for the parent comment
			//parent.Children = make([]bson.ObjectId, 10)
			children := append(parent.Children, object.ID)
			manager.Comments.Update(bson.M{"_id": parent.ID},
				bson.M{"$set": bson.M{"children": children}})
		}
	}

	return nil
}

func setActions(object *model.Comment, manager *MongoManager) error {
	var user model.User
	var invalid_users []string

	for i := 0; i < len(object.Actions); i++ {
		one := &object.Actions[i]
		manager.Users.Find(bson.M{"src_id": one.SourceUserID}).One(&user)
		if user.ID == "" {
			invalid_users = append(invalid_users, one.SourceUserID)
			continue
		}

		one.UserID = user.ID
	}

	if len(invalid_users) > 0 {
		return errors.New("Error setting Actions - Cannot find users")
	}

	return nil
}

func setNotes(object *model.Comment, manager *MongoManager) error {
	var user model.User
	var invalid_users []string

	for i := 0; i < len(object.Notes); i++ {
		one := &object.Notes[i]
		manager.Users.Find(bson.M{"src_id": one.SourceUserID}).One(&user)
		if user.ID == "" {
			invalid_users = append(invalid_users, one.SourceUserID)
			continue
		}

		one.UserID = user.ID
	}

	if len(invalid_users) > 0 {
		return errors.New("Error setting Actions - Cannot find users")
	}

	return nil
}
