package model

import (
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

//Various Constants
const (

	//Various Stats (counts)
	StatsLikes    string = "likes"
	StatsFlags    string = "flags"
	StatsComments string = "comments"

	//Various Collections
	Users      string = "users"
	Assets     string = "assets"
	Actions    string = "actions"
	Comments   string = "comments"
	Tags       string = "tags"
	TagTargets string = "tag_targets"
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

// Metadata denotes a request to add/update Metadata for an entity
type Metadata struct {
	Target   string        `json:"target" bson:"target" validate:"required"`
	TargetID bson.ObjectId `json:"target_id" bson:"target_id" validate:"required"`
	Source   ImportSource  `json:"source" bson:"source"`
	Metadata bson.M        `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// Index denotes a request to add Index to various entities.
type Index struct {
	Target string    `json:"target" bson:"target" validate:"required"`
	Index  mgo.Index `json:"index" bson:"index" validate:"required"`
}
