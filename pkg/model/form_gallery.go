package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type FormGalleryAnswer struct {
	SubmissionId bson.ObjectId        `json:"submission_id" bson:"submission_id"`
	AnswerId     string               `json:"answer_id" bson:"answer_id"`
	Answer       FormSubmissionAnswer `json:"answer,omitempty" bson:"answer,omitempty"` // not saved to db, hydrated when reading only!
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
