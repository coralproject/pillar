package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Comment denotes a comment by a user in the system.
type Comment struct {
	ID           bson.ObjectId   `json:"id" bson:"_id"`
	UserID       bson.ObjectId   `json:"user_id" bson:"user_id"`
	AssetID      bson.ObjectId   `json:"asset_id" bson:"asset_id"`
	ParentID     bson.ObjectId   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	Children     []bson.ObjectId `json:"children,omitempty" bson:"children,omitempty"`
	Body         string          `json:"body" bson:"body" validate:"required"`
	Status       string          `json:"status" bson:"status"`
	DateCreated  time.Time       `json:"date_created" bson:"date_created"`
	DateUpdated  time.Time       `json:"date_updated" bson:"date_updated"`
	DateApproved time.Time       `json:"date_approved,omitempty" bson:"date_approved,omitempty"`
	Actions      []bson.ObjectId `json:"actions,omitempty" bson:"actions,omitempty"`
	Notes        []Note          `json:"notes,omitempty" bson:"notes,omitempty"`
	Tags         []string        `json:"tags,omitempty" bson:"tags,omitempty"`
	Stats        bson.M          `json:"stats,omitempty" bson:"stats,omitempty"`
	Metadata     bson.M          `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Source       ImportSource    `json:"source,omitempty" bson:"source,omitempty"`
}

// Validate performs validation on a Comment value before it is processed.
func (com Comment) Validate() error {
	errs := validate.Struct(com)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
