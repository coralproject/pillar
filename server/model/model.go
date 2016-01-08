package model

import (
	"fmt"
	"time"

	"gopkg.in/bluesuncorp/validator.v6"
	"gopkg.in/mgo.v2/bson"
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

//type DBType interface {
//	Id() bson.ObjectId
//}

//==============================================================================

//Various Constants
const (

	//action types
	ActionTypeLikes string = "Likes"
	ActionTypeFlags string = "Flags"

	//target types
	TargetTypeUser    string = "User"
	TargetTypeComment string = "Comment"
)

// Note denotes a note by a user in the system.
type Note struct {
	UserID     bson.ObjectId `json:"user_id" bson:"user_id"`
	Body       string        `json:"body" bson:"body" validate:"required"`
	Date       time.Time     `json:"date" bson:"date" validate:"required"`
	TargetID   bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	TargetType string        `json:"target_type" bson:"target_type" validate:"required"`
	Source     NoteSource    `json:"source" bson:"source"`
}

// NoteSource encapsulates all original id from the source system
type NoteSource struct {
	ID       string `json:"id,omitempty" bson:"id,omitempty"`
	UserID   string `json:"user_id" bson:"user_id" validate:"required"`
	TargetID string `json:"target_id" bson:"target_id" validate:"required"`
}

//==============================================================================

// Asset denotes an asset in the system e.g. an article or a blog etc.
type Asset struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	SourceID   string        `json:"src_id,omitempty" bson:"src_id,omitempty"` // This is the original ID (in the external source) for the asset
	URL        string        `json:"url" bson:"url" validate:"required"`
	Taxonomies []Taxonomy    `json:"taxonomies,omitempty" bson:"taxonomies,omitempty"`
}

// Taxonomy holds all name-value pairs.
type Taxonomy struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

//func (object Asset) Id() bson.ObjectId {
//	return object.ID
//}

// Validate performs validation on an Asset value before it is processed.
func (a Asset) Validate() error {
	errs := validate.Struct(a)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

//==============================================================================

// Action denotes an action taken by an actor (User) on someone/something.
// TargetType and TargetID may be zero value if data is a sub-document of the Target
// UserID may be zero value if the data is a subdocument of the actor
type Action struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Type       string        `json:"type" bson:"type" validate:"required"`
	UserID     bson.ObjectId `json:"user_id" bson:"user_id" validate:"required"`
	TargetID   bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	TargetType string        `json:"target_type" bson:"target_type" validate:"required"`
	Date       time.Time     `json:"date" bson:"date" validate:"required"`
	Value      string        `json:"value,omitempty" bson:"value,omitempty"`
	Source     ActionSource  `json:"source" bson:"source"`
}

// ActionSource encapsulates all original id from the source system
type ActionSource struct {
	ID       string `json:"id,omitempty" bson:"id,omitempty"`
	UserID   string `json:"user_id" bson:"user_id" validate:"required"`
	TargetID string `json:"target_id" bson:"target_id"`
}

//==============================================================================

// User denotes a user in the system.
type User struct {
	ID          bson.ObjectId   `json:"id" bson:"_id"`
	SourceID    string          `json:"src_id" bson:"src_id" validate:"required"`
	UserName    string          `json:"user_name" bson:"user_name" validate:"required"`
	Avatar      string          `json:"avatar" bson:"avatar" validate:"required"`
	Status      string          `json:"status" bson:"status" validate:"required"`
	LastLogin   time.Time       `json:"last_login,omitempty" bson:"last_login,omitempty"`
	MemberSince time.Time       `json:"member_since,omitempty" bson:"member_since,omitempty"`
	Actions     []bson.ObjectId `json:"actions" bson:"actions"`
	Notes       []Note          `json:"notes,omitempty" bson:"notes,omitempty"`
	//	Stats       map[string]interface{} `json:"stats" bson:"stats"`
	//	Source      map[string]interface{} `json:"source" bson:"source"` // source document if imported
}

//func (object User) Id() bson.ObjectId {
//	return object.ID
//}

// Validate performs validation on a User value before it is processed.
func (u User) Validate() error {
	errs := validate.Struct(u)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

//==============================================================================

// Comment denotes a comment by a user in the system.
type Comment struct {
	ID           bson.ObjectId   `json:"id" bson:"_id"`
	UserID       bson.ObjectId   `json:"user_id" bson:"user_id"`
	AssetID      bson.ObjectId   `json:"asset_id" bson:"asset_id"`
	ParentID     bson.ObjectId   `json:"parent_id,omitempty" bson:"parent_d,omitempty"`
	Children     []bson.ObjectId `json:"children,omitempty" bson:"children,omitempty"`
	Body         string          `json:"body" bson:"body" validate:"required"`
	Status       string          `json:"status" bson:"status"`
	DateCreated  time.Time       `json:"date_created" bson:"date_created"`
	DateUpdated  time.Time       `json:"date_updated" bson:"date_updated"`
	DateApproved time.Time       `json:"date_approved,omitempty" bson:"date_approved,omitempty"`
	Actions      []bson.ObjectId `json:"actions" bson:"actions"`
	ActionCounts map[string]int  `json:"actionCounts" bson:"actionCounts"`
	Notes        []Note          `json:"notes" bson:"notes"`
	Source       CommentSource   `json:"source" bson:"source"`
	//	Stats        map[string]interface{} `json:"stats" bson:"stats"`
}

// CommentSource encapsulates all original id from the source system
type CommentSource struct {
	ID       string `json:"id" bson:"id" validate:"required"`
	AssetID  string `json:"asset_id" bson:"asset_id" validate:"required"`
	UserID   string `json:"user_id" bson:"user_id" validate:"required"`
	ParentID string `json:"parent_id" bson:"parent_id"`
}

//func (object Comment) Id() bson.ObjectId {
//	return object.ID
//}

// Validate performs validation on a Comment value before it is processed.
func (com Comment) Validate() error {
	errs := validate.Struct(com)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

//==============================================================================
