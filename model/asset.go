package model

import (
	"fmt"
)

// Taxonomy holds all name-value pairs.
type Taxonomy struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

// Asset denotes an asset in the system e.g. an article or a blog etc.
type Asset struct {
	AssetID    string     `json:"asset_id" bson:"_id"`
	SourceID   string     `json:"src_id" bson:"src_id"`
	URL        string     `json:"url" bson:"url" validate:"url"`
	Taxonomies []Taxonomy `json:"taxonomies" bson:"taxonomies"`
}

// Validate performs validation on an Asset value before it is processed.
func (a *Asset) Validate() error {
	errs := validate.Struct(a)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
