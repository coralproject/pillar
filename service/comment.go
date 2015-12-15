package service

import (
	"fmt"
	"github.com/coralproject/pillar/model"
	"gopkg.in/mgo.v2/bson"
)

// CreateComment creates a new comment resource
func CreateComment(input model.Comment) (*model.Comment, error) {

	// Insert Comment
	manager := getMongoManager()
	defer manager.close()

	dbEntity := model.Comment{}

	//return, if exists
	manager.comments.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with ID [%s]", input.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.comments.Find(bson.M{"source.id": input.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with URL [%s]", input.Source.ID)
		return &dbEntity, nil
	}

	input.ID = bson.NewObjectId()
	err := manager.comments.Insert(input)
	if err != nil {
		return nil, err
	}

	return &dbEntity, nil
}
