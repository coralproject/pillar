package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Action denotes an action taken by an actor (User) on someone/something.
// TargetType and TargetID may be zero value if data is a sub-document of the Target
// UserID may be zero value if the data is a subdocument of the actor
type Action struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Type     string        `json:"type" bson:"type" validate:"required"`
	UserID   bson.ObjectId `json:"user_id" bson:"user_id" validate:"required"`
	Target   string        `json:"target" bson:"target" validate:"required"`
	TargetID bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	Date     time.Time     `json:"date" bson:"date" validate:"required"`
	Value    string        `json:"value,omitempty" bson:"value,omitempty"`
	Metadata bson.M        `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Source   ImportSource  `json:"source,omitempty" bson:"source,omitempty"`
}
