package crud

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
)

// CreateAsset creates a new asset resource
func CreateAsset(object *Asset) (*Asset, *AppError) {

	// Insert Asset
	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := Asset{}

	//return, if exists
	manager.Assets.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//return, if entity exists with same source.id
	manager.Assets.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with source [%s]\n", reflect.TypeOf(object).Name(), object.Source.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//return, if entity exists with same url
	manager.Assets.Find(bson.M{"url": object.URL}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with url [%s]\n", reflect.TypeOf(object).Name(), object.URL)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	object.ID = bson.NewObjectId()
	err := manager.Assets.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating asset [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	err = CreateTagStats(manager, object.Tags, &TagTarget{Target: CollectionAsset, TargetID: object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return object, nil
}
