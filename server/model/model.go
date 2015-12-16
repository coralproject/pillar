package model

import (
	"fmt"
	"gopkg.in/bluesuncorp/validator.v6"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	config := validator.Config{
		TagName:         "validate",
		ValidationFuncs: validator.BakedInValidators,
	}

	validate = validator.New(config)
}

//==============================================================================

// Action denotes an action taken by an actor (User) on someone/something.
//   TargetType and Target id may be zero value if data is a subdocument of the Target
//   UserID may be zero value if the data is a subdocument of the actor
type Action struct {
	UserID     bson.ObjectId `json:"user_id" bson:"user_id" validate:"required"`
	Type       string        `json:"type" bson:"type"`
	Value      string        `json:"value" bson:"value"`
	TargetType string        `json:"target_type" bson:"target_type"` // eg: comment, "" for actions existing within target documents
	TargetID   string        `json:"target_type" bson:"target_type"` // eg: 23423
	Date       time.Time     `json:"date" bson:"date"`
}

// Note denotes a note by a user in the system.
type Note struct {
	UserID bson.ObjectId `json:"user_id" bson:"user_id"`
	Body   string        `json:"body" bson:"body" validate:"required"`
	Date   time.Time     `json:"date" bson:"date"`
}

// CommentSource encapsulates all original id from the source system
type CommentSource struct {
	ID       string `json:"id" bson:"id" validate:"required"`
	AssetID  string `json:"asset_id" bson:"asset_id" validate:"required"`
	UserID   string `json:"user_id" bson:"user_id" validate:"required"`
	ParentID string `json:"parent_id" bson:"parent_id"`
}

// Comment denotes a comment by a user in the system.
type Comment struct {
	ID           bson.ObjectId   `json:"id" bson:"_id"`
	UserID       bson.ObjectId   `json:"user_id" bson:"user_id"`
	ParentID     bson.ObjectId   `json:"parent_id" bson:"parent_d"`
	AssetID      bson.ObjectId   `json:"asset_id" bson:"asset_id"`
	Children     []bson.ObjectId `json:"children" bson:"children"`
	Body         string          `json:"body" bson:"body" validate:"required"`
	Status       string          `json:"status" bson:"status"`
	DateCreated  time.Time       `json:"date_created" bson:"date_created"`
	DateUpdated  time.Time       `json:"date_updated" bson:"date_updated"`
	DateApproved time.Time       `json:"date_approved" bson:"date_approved"`
	Actions      []Action        `json:"actions" bson:"actions"`
	ActionCounts map[string]int  `json:"actionCounts" bson:"actionCounts"`
	Notes        []Note          `json:"notes" bson:"notes"`
	Source       CommentSource   `json:"source" bson:"source"`
	//	Stats        map[string]interface{} `json:"stats" bson:"stats"`
}

// Validate performs validation on a Comment value before it is processed.
func (com *Comment) Validate() error {
	errs := validate.Struct(com)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

//==============================================================================

// Taxonomy holds all name-value pairs.
type Taxonomy struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

// Asset denotes an asset in the system e.g. an article or a blog etc.
type Asset struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	SourceID   string        `json:"src_id,omitempty" bson:"src_id,omitempty"`
	URL        string        `json:"url" bson:"url" validate:"required"`
	Taxonomies []Taxonomy    `json:"taxonomies,omitempty" bson:"taxonomies,omitempty"`
}

// Validate performs validation on an Asset value before it is processed.
func (a *Asset) Validate() error {
	errs := validate.Struct(a)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

//==============================================================================

// User denotes a user in the system.
type User struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	SourceID    string        `json:"src_id" bson:"src_id" validate:"required"`
	UserName    string        `json:"user_name" bson:"user_name" validate:"required"`
	Avatar      string        `json:"avatar" bson:"avatar" validate:"required"`
	Status      string        `json:"status" bson:"status" validate:"required"`
	LastLogin   time.Time     `json:"last_login,omitempty" bson:"last_login,omitempty"`
	MemberSince time.Time     `json:"member_since,omitempty" bson:"member_since,omitempty"`
	ActionsBy   []Action      `json:"actions_by,omitempty" bson:"actions_by,omitempty"`
	ActionsOn   []Action      `json:"actions_on,omitempty" bson:"actions_on,omitempty"`
	Notes       []Note        `json:"notes,omitempty" bson:"notes,omitempty"`
	//	Stats       map[string]interface{} `json:"stats" bson:"stats"`
	//	Source      map[string]interface{} `json:"source" bson:"source"` // source document if imported
}

// Validate performs validation on a User value before it is processed.
func (u *User) Validate() error {
	errs := validate.Struct(u)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
