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
	Authors    []Author      `json:"authors,omitempty" bson:"authors,omitempty"`
	Section    string        `json:"section,omitempty" bson:"section,omitempty"`
	Subsection string        `json:"subsection,omitempty" bson:"subsection,omitempty"`
	Stats      bson.M        `json:"stats,omitempty" bson:"stats,omitempty"`
	Metadata   bson.M        `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Source     ImportSource  `json:"source,omitempty" bson:"source,omitempty"`
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
