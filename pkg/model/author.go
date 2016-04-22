package model

import "gopkg.in/mgo.v2/bson"

// Author denotes a writer who curated an Asset (or story).
type Author struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name" validate:"required"`
	URL      string `json:"url,omitempty" bson:"url,omitempty"`
	Twitter  string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Facebook string `json:"facebook,omitempty" bson:"facebook,omitempty"`
	Stats    bson.M `json:"stats,omitempty" bson:"stats,omitempty"`
	Metadata bson.M `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// Id returns the ID for this Model
func (object Author) Id() string {
	return object.ID
}

// Validate validates this Model
func (object Author) Validate() error {
	//	errs := validate.Struct(object)
	//	if errs != nil {
	//		return fmt.Errorf("%v", errs)
	//	}

	return nil
}
