package service

import (
	"fmt"
	"github.com/coralproject/pillar/model"
	"gopkg.in/mgo.v2/bson"
)

// CreateAsset creates a new asset resource
func CreateAsset(input model.Asset) (*model.Asset, error) {

	// Insert Comment
	manager := getMongoManager()
	defer manager.close()

	dbEntity := model.Asset{}

	//return, if exists
	manager.assets.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with ID [%s]", input.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.assets.Find(bson.M{"url": input.URL}).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("Entity exists with URL [%s]", input.URL)
		return &dbEntity, nil
	}

	input.ID = bson.NewObjectId()
	err := manager.assets.Insert(input)
	if err != nil {
		return nil, err
	}

	return &dbEntity, nil
}
