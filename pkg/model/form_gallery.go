package model

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type FormGalleryAnswer struct {
	SubmissionId    bson.ObjectId          `json:"submission_id" bson:"submission_id"`
	AnswerId        string                 `json:"answer_id" bson:"answer_id"`
	Answer          FormSubmissionAnswer   `json:"answer,omitempty" bson:"answer,omitempty"`                     // not saved to db, hydrated when reading only!
	IdentityAnswers []FormSubmissionAnswer `json:"identity_answers,omitempty" bson:"identity_answers,omitempty"` // not saved to db, hydrated when reading only!
}

type FormGallery struct {
	ID          bson.ObjectId          `json:"id" bson:"_id"`
	FormId      bson.ObjectId          `json:"form_id" bson:"form_id"`
	Headline    string                 `json:"headline" bson:"headline"`
	Description string                 `json:"description" bson:"description"`
	Config      map[string]interface{} `json:"config" bson:"config"`
	Answers     []FormGalleryAnswer    `json:"answers" bson:"answers"`
	DateCreated time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}

// I am, form_gallery
func (o FormGallery) GetType() string {
	return "form_gallery"
}

// Record all Historical Events for FormGalleries
func (o FormGallery) IsRecordableEvent(e string) bool {
	return true
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
