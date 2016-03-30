package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Search denotes a search bound by a query and tag.
type Search struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name" validate:"required"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`
	Query    	string        `json:"query" bson:"query" validate:"required"`
	Tag     	string        `json:"tag" bson:"tag" validate:"required"`
	Filters     []interface{} `json:"filters,omitempty" bson:"filters,omitempty"`
	Results     []interface{} `json:"results,omitempty" bson:"results,omitempty"`
	DateCreated time.Time     `json:"date_created" bson:"date_created" validate:"required"`
	DateUpdated time.Time     `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
	UserCreated string        `json:"user_created,omitempty" bson:"user_created,omitempty"`
	UserUpdated string        `json:"user_updated,omitempty" bson:"user_updated,omitempty"`
}
