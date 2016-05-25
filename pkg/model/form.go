package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type FormWidget struct {
	ID        int64       `json:"id" bson:"_id"`
	Type      string      `json:"type" bson:"type"`
	Component string      `json:"component" bson:"component"`
	Title     string      `json:"title" bson:"title"`
	Wrapper   interface{} `json:"wrapper" bson:"wrapper"`
	Props     interface{} `json:"props" bson:"props"`
}

type FormStep struct {
	ID      int64        `json:"id" bson:"_id"`
	Name    string       `json:"name" bson:"name"`
	Widgets []FormWidget `json:"widgets" bson:"widgets"`
}

type Form struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Settings       interface{}   `json:"settings" bson:"settings"`
	Header         interface{}   `json:"header" bson:"header"`
	Footer         interface{}   `json:"footer" bson:"footer"`
	FinishedScreen interface{}   `json:"finishedScreen" bson:"finishedScreen"`
	Steps          []FormStep    `json:"steps" bson:"steps"`
	DateCreated    time.Time     `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateUpdated    time.Time     `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
}

// Id returns the ID for this Model
func (object Form) Id() string {
	return object.ID.Hex()
}

// ImportSource returns the Source model
func (object Form) ImportSource() *ImportSource {
	return nil //  Forms have no imports
}

func (object Form) Validate() error {
	errs := validate.Struct(object)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
