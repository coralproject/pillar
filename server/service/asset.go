package service

import (
	"reflect"
	"github.com/coralproject/pillar/server/log"
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
)


// CreateAsset creates a new asset resource
func CreateAsset(object model.Asset) (*model.Asset, error) {

	// Insert Comment
	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := model.Asset{}

	//return, if exists
	manager.Assets.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		log.Logger.Printf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return &dbEntity, nil
	}

	//return, if entity exists with same src_id
	manager.Assets.Find(bson.M{"src_id": object.SourceID}).One(&dbEntity)
	if dbEntity.ID != "" {
		log.Logger.Printf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.SourceID)
		return &dbEntity, nil
	}

	//return, if entity exists with same url
	manager.Assets.Find(bson.M{"url": object.URL}).One(&dbEntity)
	if dbEntity.ID != "" {
		log.Logger.Printf("%s exists with url [%s]\n", reflect.TypeOf(object).Name(), object.URL)
		return &dbEntity, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.Assets.Insert(object)
	if err != nil {
		log.Logger.Printf("Error creating asset [%s]", err);
		return nil, err
	}

	return &object, nil
}
