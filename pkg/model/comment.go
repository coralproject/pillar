package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Comment denotes a comment by a user in the system.
type Comment struct {
	ID          bson.ObjectId   `json:"id" bson:"_id"`
	UserID      bson.ObjectId   `json:"user_id" bson:"user_id"`
	AssetID     bson.ObjectId   `json:"asset_id" bson:"asset_id"`
	RootID      bson.ObjectId   `json:"root_id,omitempty" bson:"root_id,omitempty"`
	Parents     []bson.ObjectId `json:"parents,omitempty" bson:"parents,omitempty"`
	Children    []bson.ObjectId `json:"children,omitempty" bson:"children,omitempty"`
	Body        string          `json:"body" bson:"body" validate:"required"`
	Status      string          `json:"status" bson:"status"`
	DateCreated time.Time       `json:"date_created" bson:"date_created"`
	DateUpdated time.Time       `json:"date_updated" bson:"date_updated"`
	Actions     []Action        `json:"actions,omitempty" bson:"actions,omitempty"`
	Notes       []Note          `json:"notes,omitempty" bson:"notes,omitempty"`
	Tags        []string        `json:"tags,omitempty" bson:"tags,omitempty"`
	Roles       []string        `json:"roles,omitempty" bson:"roles,omitempty"`
	Stats       bson.M          `json:"stats,omitempty" bson:"stats,omitempty"`
	Metadata    bson.M          `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Source      ImportSource    `json:"source,omitempty" bson:"source,omitempty"`
}

// Content encapsulates content-type and its data.
type Content struct {
	MimeType string `json:"type" bson:"type" validate:"required"`
	Body     string `json:"body" bson:"body" validate:"required"`
}

// CommentHistory encapsulates a snapshot of comment when it was last updated.
type CommentHistory struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	UserID  bson.ObjectId `json:"user_id" bson:"user_id" validate:"required"`
	Comment Comment       `json:"comment" bson:"comment" validate:"required"`
	Event   string        `json:"event" bson:"event" validate:"required"`
	Date    time.Time     `json:"date" bson:"date" validate:"required"`
}

// Id returns the ID for this Model
func (object Comment) Id() string {
	return object.ID.Hex()
}

// Validate validates this Model
func (object Comment) Validate() error {
	errs := validate.Struct(object)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
