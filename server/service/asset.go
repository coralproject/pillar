package service

import (
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

// CreateAsset creates a new asset resource
func CreateAsset(object model.Asset) (*model.Asset, error) {

	// Insert Comment
	manager := getMongoManager()
	defer manager.close()

	dbEntity := model.Asset{}

	//return, if exists
	manager.assets.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.assets.Find(bson.M{"url": object.URL}).One(&dbEntity)
	if dbEntity.ID != "" {
		fmt.Printf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.URL)
		return &dbEntity, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.assets.Insert(object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}
