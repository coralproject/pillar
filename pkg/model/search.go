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
	Query       string        `json:"query" bson:"query" validate:"required"`
	Tag         string        `json:"tag" bson:"tag" validate:"required"`
	Filters     bson.M        `json:"filters,omitempty" bson:"filters,omitempty"`
	Result      SearchResult  `json:"result,omitempty" bson:"result,omitempty"`
	DateCreated time.Time     `json:"date_created" bson:"date_created" validate:"required"`
	DateUpdated time.Time     `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
	UserCreated string        `json:"user_created,omitempty" bson:"user_created,omitempty"`
	UserUpdated string        `json:"user_updated,omitempty" bson:"user_updated,omitempty"`
}

type SearchHistory struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Action string        `json:"action" bson:"action" validate:"required"`
	Date   time.Time     `json:"date" bson:"date" validate:"required"`
	Search Search        `json:"search" bson:"search" validate:"required"`
}

type SearchResult struct {
	Count     int           `json:"count,omitempty" bson:"count,omitempty"`
	Histogram []interface{} `json:"histogram,omitempty" bson:"histogram,omitempty"`
	Users     []User        `json:"users,omitempty" bson:"users,omitempty"`
}
