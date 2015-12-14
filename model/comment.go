package model

import (
	"time"
)

// Action denotes an action taken by an actor (User) on someone/something.
//   TargetType and Target id may be zero value if data is a subdocument of the Target
//   UserID may be zero value if the data is a subdocument of the actor
type Action struct {
	UserID     string    `json:"user_id" bson:"user_id" validate:"required"`
	Type       string    `json:"type" bson:"type"`
	Value      string    `json:"value" bson:"value"`
	TargetType string    `json:"target_type" bson:"target_type"` // eg: comment, "" for actions existing within target documents
	TargetID   string    `json:"target_type" bson:"target_type"` // eg: 23423
	Date       time.Time `json:"date" bson:"date"`
}

// Note denotes a note by a user in the system.
type Note struct {
	UserID string    `json:"user_id" bson:"user_id"`
	Body   string    `json:"body" bson:"body" validate:"required"`
	Date   time.Time `json:"date" bson:"date"`
}

// Comment denotes a comment by a user in the system.
type Comment struct {
	CommentID    string                 `json:"comment_id" bson:"comment_id"`
	UserID       string                 `json:"user_id" bson:"user_id" validate:"required"`
	ParentID     string                 `json:"parent_id" bson:"parent_d"`
	AssetID      string                 `json:"asset_id" bson:"asset_id"`
	Children     []string               `json:"children" bson:"children"` // experimental
	Path         string                 `json:"path" bson:"path"`
	Body         string                 `json:"body" bson:"body" validate:"required"`
	Status       string                 `json:"status" bson:"status"`
	DateCreated  time.Time              `json:"date_created" bson:"date_created"`
	DateUpdated  time.Time              `json:"date_updated" bson:"date_updated"`
	DateApproved time.Time              `json:"date_approved" bson:"date_approved"`
	Actions      []Action               `json:"actions" bson:"actions"`
	ActionCounts map[string]int         `json:"actionCounts" bson:"actionCounts"`
	Notes        []Note                 `json:"notes" bson:"notes"`
	Stats        map[string]interface{} `json:"stats" bson:"stats"`
	Source       map[string]interface{} `json:"source" bson:"source"` // source document if imported
}

// CreateComment creates a new comment resource
func CreateComment(comment Comment) (*Comment, error) {

	// Write the user to mongo
	manager := getMongoManager(collectionComment)
	defer manager.close()

	err1 := manager.collection.Insert(comment)
	if err1 != nil {
		return nil, err1
	}

	return &comment, nil
}
