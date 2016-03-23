package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// UserGroup denotes a group of users bound by the filters and tags here.
type UserGroup struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name" validate:"required"`
	Description string        `json:"description" bson:"description" validate:"required"`
	Filters     []bson.M      `json:"filters,omitempty" bson:"filters,omitempty"`
	IncludeTags []string      `json:"include_tags,omitempty" bson:"include_tags,omitempty"`
	ExcludeTags []string      `json:"exclude_tags,omitempty" bson:"exclude_tags,omitempty"`
	DateCreated time.Time     `json:"date_created" bson:"date_created" validate:"required"`
	DateUpdated time.Time     `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
	UserCreated string        `json:"user_created,omitempty" bson:"user_created,omitempty"`
	UserUpdated string        `json:"user_updated,omitempty" bson:"user_updated,omitempty"`
}
