package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
)

// Note denotes a note by a user in the system.
type Note struct {
	UserID   bson.ObjectId `json:"user_id" bson:"user_id"`
	Body     string        `json:"body" bson:"body" validate:"required"`
	Date     time.Time     `json:"date" bson:"date" validate:"required"`
	TargetID bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	Target   string        `json:"target" bson:"target" validate:"required"`
	Source   ImportSource  `json:"source,omitempty" bson:"source,omitempty"`
}

// Validate validates this Model
func (object Note) Validate() error {
	errs := validate.Struct(object)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
