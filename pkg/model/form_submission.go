package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// this is what we expect for input for a form submission
type FormSubmissionAnswerInput struct {
	WidgetId string      `json:"widget_id"`
	Answer   interface{} `json:"answer"`
}

type FormSubmissionInput struct {
	FormId  string                      `jsont:"form_id"`
	Status  string                      `json:"status" bson:"status"`
	Answers []FormSubmissionAnswerInput `json:"replies"`
}

// here's what a form submission is
type FormSubmissionAnswer struct {
	WidgetId     string      `json:"widget_id" bson:"widget_id"`
	Answer       interface{} `json:"answer" bson:"answer"`
	EditedAnswer interface{} `json:"edited" bson:"edited"`
	Question     interface{} `json:"question" bson:"question"`
	Props        interface{} `json:"props" bson:"props"`
}

type FormSubmission struct {
	ID             bson.ObjectId          `json:"id" bson:"_id"`
	FormId         bson.ObjectId          `json:"form_id" bson:"form_id"`
	Status         string                 `json:"status" bson:"status"`
	Answers        []FormSubmissionAnswer `json:"replies" bson:"replies"`
	Header         interface{}            `json:"header" bson:"header"`
	Footer         interface{}            `json:"footer" bson:"footer"`
	FinishedScreen interface{}            `json:"finishedScreen" bson:"finishedScreen"`
	CreatedBy      interface{}            `json:"created_by" bson:"created_by"` // Todo, decide how to represent ownership here
	UpdatedBy      interface{}            `json:"updated_by" bson:"updated_by"` // Todo, decide how to represent ownership here
	DateCreated    time.Time              `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated    time.Time              `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}

// Id returns the ID for this Model
func (object FormSubmission) Id() string {
	return object.ID.Hex()
}

// ImportSource returns the Source model
func (object FormSubmission) ImportSource() *ImportSource {
	return nil //  Forms have no imports
}

func (object FormSubmission) Validate() error {
	errs := validate.Struct(object)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

// Id returns the ID for this Model
func (object FormSubmissionInput) Id() string {
	return ""
}

// ImportSource returns the Source model
func (object FormSubmissionInput) ImportSource() *ImportSource {
	return nil //  Forms have no imports
}

func (object FormSubmissionInput) Validate() error {
	errs := validate.Struct(object)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
