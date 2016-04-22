package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
)

// Tag denotes a unique tag in the system
type Tag struct {
	Name        string    `json:"name" bson:"_id" validate:"required"`
	OldName     string    `json:"old_name,omitempty" bson:"old_name,omitempty"`
	Description string    `json:"description" bson:"description" validate:"required"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}

// Id returns the ID for this Model
func (object Tag) Id() string {
	return object.Name
}

// Validate validates this Model
func (object Tag) Validate() error {
	errs := validate.Struct(object)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

// TagTarget denotes relationship between an entity and its tags
type TagTarget struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Target      string        `json:"target" bson:"target" validate:"required"`
	TargetID    bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	Name        string        `json:"name" bson:"name" validate:"required"`
	DateCreated time.Time     `json:"date_created" bson:"date_created"`
}
