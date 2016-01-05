package service

import (
	"reflect"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
)

//	"fmt"
//	"reflect"

// CreateAsset creates a new asset resource
func CreateAsset(object model.Asset) (*model.Asset, error) {

	// Insert Comment
	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := model.Asset{}

	//return, if exists
	manager.Assets.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		log.Dev("service", "CreateAsset", "%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return &dbEntity, nil
	}

	//return, if exists
	manager.Assets.Find(bson.M{"url": object.URL}).One(&dbEntity)
	if dbEntity.ID != "" {
		log.Dev("service", "CreateAsset", "%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.URL)
		return &dbEntity, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.Assets.Insert(object)
	if err != nil {
		log.Error("service", "CreateAsset", err, "Insert assets")
		return nil, err
	}

	return &object, nil
}
