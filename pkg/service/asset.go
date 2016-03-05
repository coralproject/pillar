package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// ImportAsset imports a new asset into Coral
func ImportAsset(context *AppContext) (*model.Asset, *AppError) {

	db := context.DB
	input := context.Input.(model.Asset)
	var dbEntity model.Asset

	//Upsert if entity exists with same source.id
	db.Assets.Find(bson.M{"source.id": input.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		input.ID = dbEntity.ID
		if _, err := db.Assets.UpsertId(dbEntity.ID, &input); err != nil {
			message := fmt.Sprintf("Error updating existing Asset [%s]", input.Source.ID)
			return nil, &AppError{err, message, http.StatusInternalServerError}
		}
		return &input, nil
	}

	//return, if entity exists
	if dbEntity := assetExists(db, &input); dbEntity != nil {
		message := fmt.Sprintf("Asset exists, id [%s] and url [%s] must be unique.", input.ID, input.URL)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return doCreateAsset(db, &input)
}

// CreateUpdateAsset creates/updates an asset resource
func CreateUpdateAsset(context *AppContext) (*model.Asset, *AppError) {
	input := context.Input.(model.Asset)
	if input.ID == "" {
		return createAsset(context)
	}

	return updateAsset(context)
}

// createAsset creates a new asset resource
func createAsset(context *AppContext) (*model.Asset, *AppError) {

	db := context.DB
	input := context.Input.(model.Asset)

	//return, if entity exists
	if dbEntity := assetExists(db, &input); dbEntity != nil {
		message := fmt.Sprintf("Asset exists, id [%s] and url [%s] must be unique.", input.ID, input.URL)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return doCreateAsset(db, &input)
}

// UpdateAsset updates an existing asset
func updateAsset(context *AppContext) (*model.Asset, *AppError) {
	db := context.DB
	input := context.Input.(model.Asset)
	var dbEntity model.Asset

	//entity not found, return
	db.Assets.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID == "" {
		message := fmt.Sprintf("Asset not found [%s]\n", input.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	dbEntity.Tags = input.Tags
	if err := db.Assets.UpdateId(dbEntity.ID, bson.M{"$set": bson.M{"tags": dbEntity.Tags}}); err != nil {
		message := fmt.Sprintf("Error updating asset [%+v]\n", input)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return &dbEntity, nil
}

//finds and returns an asset if exists, else nil
func assetExists(db *db.MongoDB, input *model.Asset) *model.Asset {
	var dbEntity model.Asset

	//return, if exists
	db.Assets.FindId(input.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		return &dbEntity
	}

	//return, if entity exists with same url
	db.Assets.Find(bson.M{"url": input.URL}).One(&dbEntity)
	if dbEntity.ID != "" {
		return &dbEntity
	}

	return nil
}

//inserts an asset to the db and any related post-processing
func doCreateAsset(db *db.MongoDB, input *model.Asset) (*model.Asset, *AppError) {
	//assign a new ObjectId
	input.ID = bson.NewObjectId()

	if err := db.Assets.Insert(input); err != nil {
		message := fmt.Sprintf("Error creating asset [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	//create/update authors, if any
	for _, one := range input.Authors {
		if _, err := CreateUpdateAuthor(&AppContext{db, one}); err != nil {
			//return nil, err
		}
	}

	c := AppContext{db, model.Section{Name: input.Section}}
	if _, err := CreateSection(&c); err != nil {
		//return nil, err
	}

	tt := &model.TagTarget{Target: model.Assets, TargetID: input.ID}
	if err := CreateTagTargets(db, input.Tags, tt); err != nil {
		//message := fmt.Sprintf("Error creating TagStat [%s]", err)
		//return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return input, nil
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
