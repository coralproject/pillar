package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type FormGalleryAnswer struct {
	ResponseId bson.ObjectId `json:"response_id" bson:"response_id"`
	AnswerId   bson.ObjectId `json:"answer_id" bson:"answer_id"`
}

type FormGallery struct {
	ID          bson.ObjectId       `json:"id" bson:"_id"`
	FormId      bson.ObjectId       `json:"form_id" bson:"form_id"`
	Answers     []FormGalleryAnswer `json:"answers" bson:"answers"`
	DateCreated time.Time           `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated time.Time           `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}

// Id returns the ID for this Model
func (object FormGallery) Id() string {
	return object.ID.Hex()
}

// ImportSource returns the Source model
func (object FormGallery) ImportSource() *ImportSource {
	return nil //  Forms have no imports
}

func (object FormGallery) Validate() error {
	errs := validate.Struct(object)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
