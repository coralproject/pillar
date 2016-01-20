package dto

import (
	"github.com/coralproject/pillar/server/model"
	"gopkg.in/mgo.v2/bson"
)

// Metadata denotes a request to add/update Metadata for an entity
type Metadata struct {
	Target   string                 `json:"target" bson:"target" validate:"required"`
	TargetID bson.ObjectId          `json:"target_id" bson:"target_id" validate:"required"`
	Source   model.ImportSource     `json:"source" bson:"source"`
	Metadata map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}
