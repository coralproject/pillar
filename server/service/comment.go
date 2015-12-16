package service

import (
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
)

// CreateComment creates a new comment resource
func CreateComment(object model.Comment) (*model.Comment, error) {

	// Insert Comment
	manager := getMongoManager()
	defer manager.close()

	dbEntity := model.Comment{}

	//return, if exists
	manager.comments.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with ID [%s]", object.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.comments.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with URL [%s]", object.Source.ID)
		return &dbEntity, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.comments.Insert(object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}
