package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// AppContext encapsulates application specific runtime information
type AppContext struct {
	DB      *db.MongoDB
	Input   interface{}
}

func (c *AppContext) Close() {
	c.DB.Close()
}

func NewContext() *AppContext {
	c := AppContext{}
	c.DB = db.NewMongoDB()
	return &c
}

// AppError encapsulates application specific error
type AppError struct {
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// UpdateMetadata updates metadata for an entity
func UpdateMetadata(context *AppContext) (interface{}, *AppError) {

	db := context.DB
	object := context.Input.(model.Metadata)

	collection := db.Session.DB("").C(object.Target)
	var dbEntity bson.M
	collection.FindId(object.TargetID).One(&dbEntity)
	if len(dbEntity) == 0 {
		collection.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	}

	if len(dbEntity) == 0 {
		message := fmt.Sprintf("Cannot update metadata for [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	collection.Update(
		bson.M{"_id": dbEntity["_id"]},
		bson.M{"$set": bson.M{"metadata": object.Metadata}},
	)

	return dbEntity, nil
}

// CreateIndex creates indexes to various entities
func CreateIndex(context *AppContext) *AppError {

	db := context.DB
	object := context.Input.(model.Index)

	err := db.Session.DB("").C(object.Target).EnsureIndex(object.Index)
	if err != nil {
		message := fmt.Sprintf("Error creating index [%+v]", object)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}

// CreateUserAction inserts an activity by the user
func CreateUserAction(context *AppContext) *AppError {

	db := context.DB
	object := context.Input.(model.CayUserAction)

	object.ID = bson.NewObjectId()
	object.Date = time.Now()
	if object.Release == "" {
		object.Release = "0.1.0"
	}
	err := db.CayUserActions.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating user-action [%s]", err)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
