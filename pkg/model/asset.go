package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

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

// Validate performs validation on an Asset value before it is processed.
func (a Asset) Validate() error {
	errs := validate.Struct(a)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
