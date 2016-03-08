package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
	"io"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

// AppContext encapsulates application specific runtime information
type AppContext struct {
	DB      *db.MongoDB
	Body    io.ReadCloser
}

func (c *AppContext) Close() {
	c.DB.Close()
}

//Unmarshall unmarshalls the Body to the passed object
func (c *AppContext) Unmarshall(input interface{}) error {
	bytez, _ := ioutil.ReadAll(c.Body)
	if err := json.Unmarshal(bytez, input); err != nil {
		return err
	}
	return nil
}

//Marshall marshalls an incoming object and sets it to the Body
func (c *AppContext) Marshall(j interface{}) {
	bytez, _ := json.Marshal(j)
	c.Body = ioutil.NopCloser(bytes.NewReader(bytez))
}

func NewContext(body io.ReadCloser) *AppContext {
	return &AppContext{db.NewMongoDB(), body}
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
	var input model.Metadata
	json.NewDecoder(context.Body).Decode(&input)

	collection := db.Session.DB("").C(input.Target)
	var dbEntity bson.M
	collection.FindId(input.TargetID).One(&dbEntity)
	if len(dbEntity) == 0 {
		collection.Find(bson.M{"source.id": input.Source.ID}).One(&dbEntity)
	}

	if len(dbEntity) == 0 {
		message := fmt.Sprintf("Cannot update metadata for [%+v]\n", input)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	collection.Update(
		bson.M{"_id": dbEntity["_id"]},
		bson.M{"$set": bson.M{"metadata": input.Metadata}},
	)

	return dbEntity, nil
}

// CreateIndex creates indexes to various entities
func CreateIndex(context *AppContext) *AppError {

	db := context.DB
	var input model.Index
	json.NewDecoder(context.Body).Decode(&input)

	err := db.Session.DB("").C(input.Target).EnsureIndex(input.Index)
	if err != nil {
		message := fmt.Sprintf("Error creating index [%+v]", input)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}

// CreateUserAction inserts an activity by the user
func CreateUserAction(context *AppContext) *AppError {

	db := context.DB
	var input model.CayUserAction
	json.NewDecoder(context.Body).Decode(&input)

	input.ID = bson.NewObjectId()
	input.Date = time.Now()
	if input.Release == "" {
		input.Release = "0.1.0"
	}
	err := db.CayUserActions.Insert(input)
	if err != nil {
		message := fmt.Sprintf("Error creating user-action [%s]", err)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
