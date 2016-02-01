package crud

import (
	"fmt"
	"time"

	"gopkg.in/bluesuncorp/validator.v6"
	"gopkg.in/mgo.v2"
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

	//Various Stats (counts)
	StatsLikes    string = "likes"
	StatsFlags    string = "flags"
	StatsComments string = "comments"

	// Various Collections
	CollectionUser      string = "user"
	CollectionAsset     string = "asset"
	CollectionAction    string = "action"
	CollectionComment   string = "comment"
	CollectionTag       string = "tag"
	CollectionTagTarget string = "tag_target"
)

// ImportSource encapsulates all original id from the source system
// this is a common struct used primarily for import purposes
// client is responsible for passing in valid/correct source data
type ImportSource struct {
	ID       string `json:"id,omitempty" bson:"id,omitempty"`
	UserID   string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	TargetID string `json:"target_id,omitempty" bson:"target_id,omitempty"`
	AssetID  string `json:"asset_id,omitempty" bson:"asset_id,omitempty"`
	ParentID string `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
}

//==============================================================================

// Note denotes a note by a user in the system.
type Note struct {
	UserID   bson.ObjectId `json:"user_id" bson:"user_id"`
	Body     string        `json:"body" bson:"body" validate:"required"`
	Date     time.Time     `json:"date" bson:"date" validate:"required"`
	TargetID bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	Target   string        `json:"target" bson:"target" validate:"required"`
	Source   ImportSource  `json:"source" bson:"source"`
}

//==============================================================================

// Asset denotes an asset in the system e.g. an article or a blog etc.
type Asset struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	URL        string        `json:"url" bson:"url" validate:"required"`
	Tags       []string      `json:"tags,omitempty" bson:"tags,omitempty"`
	Taxonomies []Taxonomy    `json:"taxonomies,omitempty" bson:"taxonomies,omitempty"`
	Source     ImportSource  `json:"source" bson:"source"`
	Metadata   bson.M        `json:"metadata,omitempty" bson:"metadata,omitempty"`
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
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Type     string        `json:"type" bson:"type" validate:"required"`
	UserID   bson.ObjectId `json:"user_id" bson:"user_id" validate:"required"`
	Target   string        `json:"target" bson:"target" validate:"required"`
	TargetID bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	Date     time.Time     `json:"date" bson:"date" validate:"required"`
	Value    string        `json:"value,omitempty" bson:"value,omitempty"`
	Source   ImportSource  `json:"source" bson:"source"`
	Metadata bson.M        `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

//==============================================================================

// User denotes a user in the system.
type User struct {
	ID          bson.ObjectId   `json:"id" bson:"_id"`
	Name        string          `json:"name" bson:"name" validate:"required"`
	Avatar      string          `json:"avatar" bson:"avatar" validate:"required"`
	Status      string          `json:"status" bson:"status" validate:"required"`
	LastLogin   time.Time       `json:"last_login,omitempty" bson:"last_login,omitempty"`
	MemberSince time.Time       `json:"member_since,omitempty" bson:"member_since,omitempty"`
	Actions     []bson.ObjectId `json:"actions,omitempty" bson:"actions,omitempty"`
	Notes       []Note          `json:"notes,omitempty" bson:"notes,omitempty"`
	Tags        []string        `json:"tags,omitempty" bson:"tags,omitempty"`
	Source      ImportSource    `json:"source" bson:"source"`
	Stats       bson.M          `json:"stats,omitempty" bson:"stats,omitempty"`
	Metadata    bson.M          `json:"metadata,omitempty" bson:"metadata,omitempty"`
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
	Source       ImportSource    `json:"source" bson:"source"`
	Stats        bson.M          `json:"stats,omitempty" bson:"stats,omitempty"`
	Metadata     bson.M          `json:"metadata,omitempty" bson:"metadata,omitempty"`
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

// Metadata denotes a request to add/update Metadata for an entity
type Metadata struct {
	Target   string        `json:"target" bson:"target" validate:"required"`
	TargetID bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	Source   ImportSource  `json:"source" bson:"source"`
	Metadata bson.M        `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

//==============================================================================

// Index denotes a request to add Index to various entities.
type Index struct {
	Target string    `json:"target" bson:"target" validate:"required"`
	Index  mgo.Index `json:"index" bson:"index" validate:"required"`
}

//==============================================================================

// Tag denotes a unique tag in the system
type Tag struct {
	Name        string    `json:"name" bson:"_id" validate:"required"`
	Description string    `json:"description" bson:"description" validate:"required"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
}

// TagTarget denotes relationship between an entity and its tags
type TagTarget struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Target      string        `json:"target" bson:"target" validate:"required"`
	TargetID    bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	Name        string        `json:"name" bson:"name" validate:"required"`
	DateCreated time.Time     `json:"date_created" bson:"date_created"`
}
