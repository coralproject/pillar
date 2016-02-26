package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Section denotes a media section
type Section struct {
	Name        string    `json:"name" bson:"_id" validate:"required"`
	Description string    `json:"description" bson:"description" validate:"required"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
	Stats       bson.M    `json:"stats,omitempty" bson:"stats,omitempty"`
	Metadata    bson.M    `json:"metadata,omitempty" bson:"metadata,omitempty"`
}
