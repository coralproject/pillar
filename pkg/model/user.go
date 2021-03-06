package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// User denotes a user in the system.
type User struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name" validate:"required"`
	Avatar      string        `json:"avatar" bson:"avatar"`
	Status      string        `json:"status" bson:"status"`
	LastLogin   time.Time     `json:"last_login,omitempty" bson:"last_login,omitempty"`
	MemberSince time.Time     `json:"member_since,omitempty" bson:"member_since,omitempty"`
	Actions     []Action      `json:"actions,omitempty" bson:"actions,omitempty"`
	Notes       []Note        `json:"notes,omitempty" bson:"notes,omitempty"`
	Tags        []string      `json:"tags,omitempty" bson:"tags,omitempty"`
	Stats       bson.M        `json:"stats,omitempty" bson:"stats,omitempty"`
	Metadata    bson.M        `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Source      *ImportSource  `json:"source,omitempty" bson:"source,omitempty"`
}

// Id returns the ID for this Model
func (object User) Id() string {
	return object.ID.Hex()
}

// ImportSource returns the Source model
func (object User) ImportSource() *ImportSource {
	return object.Source
}

// Validate validates this Model
func (object User) Validate() error {
	errs := validate.Struct(object)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
