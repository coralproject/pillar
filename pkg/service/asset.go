package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"github.com/coralproject/pillar/pkg/db"
)

// CreateAsset creates a new asset resource
func CreateAsset(context *AppContext) (*model.Asset, *AppError) {

	db := context.DB
	object := context.Input.(model.Asset)

	// Insert Asset
	dbEntity := model.Asset{}

	//return, if exists
	db.Assets.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//upsert if entity exists with same source.id
	db.Assets.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		object.ID = dbEntity.ID
		_, err := db.Assets.UpsertId(dbEntity.ID, object)
		if err != nil {
			message := fmt.Sprintf("Error updating existing Asset [%s], %s", object.Source.ID, err)
			return nil, &AppError{err, message, http.StatusInternalServerError}
		}
		return &object, nil
	}

	//return, if entity exists with same url
	db.Assets.Find(bson.M{"url": object.URL}).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with url [%s]\n", reflect.TypeOf(object).Name(), object.URL)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	object.ID = bson.NewObjectId()
	err := db.Assets.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating asset [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	err = CreateTagTargets(db, object.Tags, &model.TagTarget{Target: model.Assets, TargetID: object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return &object, nil
}

//update stats on this asset for #comments
func updateAssetOnComment(db *db.MongoDB, asset *model.Asset) {
	if asset.Stats == nil {
		asset.Stats = make(map[string]interface{})
	}

	if asset.Stats[model.StatsComments] == nil {
		asset.Stats[model.StatsComments] = 0
	}

	asset.Stats[model.StatsComments] = asset.Stats[model.StatsComments].(int) + 1
	db.Assets.Update(
		bson.M{"_id": asset.ID},
		bson.M{"$set": bson.M{"stats": asset.Stats}},
	)
}

// CreateUpdateAsset creates/updates an asset
func CreateUpdateAsset(context *AppContext) (*model.Asset, *AppError) {
	db := context.DB
	object := context.Input.(model.Asset)

	var dbEntity *model.Asset
	//entity not found, return
	db.Assets.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID == "" {
		message := fmt.Sprintf("Asset not found [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	dbEntity.Tags = object.Tags
	if err := db.Assets.UpdateId(dbEntity.ID, bson.M{"$set": bson.M{"tags": dbEntity.Tags}}); err != nil {
		message := fmt.Sprintf("Error updating asset [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return dbEntity, nil
}

