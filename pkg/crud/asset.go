package crud

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
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
		message := fmt.Sprintf("Asset exists with ID [%s]\n", object.ID)
		return nil, &AppError{nil, message, http.StatusConflict}
	}

	//return, if entity exists with same source.id
	manager.Assets.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("Asset exists with Source.ID [%s]\n", object.Source.ID)
		return nil, &AppError{nil, message, http.StatusConflict}
	}

	//return, if entity exists with same url
	manager.Assets.Find(bson.M{"url": object.URL}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("Asset exists with URL [%s]\n", object.URL)
		return nil, &AppError{nil, message, http.StatusConflict}
	}

	object.ID = bson.NewObjectId()
	err := manager.Assets.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating asset [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	err = CreateTagTargets(manager, object.Tags, &TagTarget{Target: Assets, TargetID: object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return object, nil
}
