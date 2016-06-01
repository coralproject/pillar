package service

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
)

// calculate stats for Forms
func updateStats(context *web.AppContext) *web.AppError {

	// get the form in question
	f, err := GetForm(context)
	if err != nil {
		message := fmt.Sprintf("Could not load form to update stats")
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	// do some counting
	responses, err2 := context.MDB.DB.C(model.FormSubmissions).Find(bson.M{"form_id": f.ID}).Count()
	if err2 != nil {
		message := fmt.Sprintf("Could not perform count of form submissions")
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	// update the stats subdoc

	s := model.FormStats{}
	s.Responses = responses
	err2 = context.MDB.DB.C(model.Forms).Update(bson.M{"_id": f.ID}, bson.M{"$set": bson.M{"stats": s}})
	if err2 != nil {
		message := fmt.Sprintf("Error updating form stats")
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	return nil

}

// given a form's id and a stats, update the form with the status
func UpdateFormStatus(context *web.AppContext) (*model.Form, *web.AppError) {

	// todo, gracefully message invalid ids
	id := bson.ObjectIdHex(context.GetValue("id"))
	status := context.GetValue("status")

	// let's make sure we don't update all of them..
	q := bson.M{"_id": id}
	s := bson.M{"$set": bson.M{"status": status, "date_updated": time.Now()}}

	// do the update
	err := context.MDB.DB.C(model.Forms).Update(q, s)
	if err != nil {
		message := fmt.Sprintf("Error updating Form status")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	var f *model.Form
	err = context.MDB.DB.C(model.Forms).FindId(id).One(&f)
	if err != nil {
		message := fmt.Sprintf("Could not find Form ", id)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return f, nil

}

func CreateUpdateForm(context *web.AppContext) (*model.Form, *web.AppError) {

	var input model.Form
	if err := UnmarshallAndValidate(context, &input); err != nil {
		return nil, err
	}

	/* Todo, custom validation
	if input.Name == "" {
		message := fmt.Sprintf("Invalid Section Name [%s]", input.Name)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}
	*/

	var dbEntity model.Form
	if context.MDB.DB.C(model.Forms).FindId(input.ID).One(&dbEntity); dbEntity.ID == "" {
		input.DateCreated = time.Now()
	}
	input.DateUpdated = time.Now()

	// create
	if input.ID == "" {

		input.ID = bson.NewObjectId()

		if err := context.MDB.DB.C(model.Forms).Insert(input); err != nil {
			message := fmt.Sprintf("Error inserting Form")
			return nil, &web.AppError{err, message, http.StatusInternalServerError}
		}

		// store the id into the context as a hex
		//  to match up with what we expect from web params
		context.SetValue("id", input.ID.Hex())

	}

	// do the update
	if _, err := context.MDB.DB.C(model.Forms).UpsertId(input.ID, input); err != nil {
		message := fmt.Sprintf("Error creating/updating Form")

		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	// always update the stats
	err := updateStats(context)
	if err != nil {
		return nil, err
	}

	return &input, nil

}

// GetForms returns an array of Forms
func GetForms(context *web.AppContext) ([]model.Form, *web.AppError) {

	//set created-date for the new ones
	all := make([]model.Form, 0)
	if err := context.MDB.DB.C(model.Forms).Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching Forms")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}

// GetForms returns a single form by id
func GetForm(c *web.AppContext) (model.Form, *web.AppError) {

	// which one do they want?
	idStr := c.GetValue("id")
	if idStr == "" {
		message := fmt.Sprintf("Cannot get Form. Invalid Id [%s]", idStr)
		return model.Form{}, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//convert to an ObjectId
	id := bson.ObjectIdHex(idStr)

	f := model.Form{}
	if err := c.MDB.DB.C(model.Forms).Find(bson.M{"_id": id}).One(&f); err != nil {
		message := fmt.Sprintf("Error fetching Forms")
		return model.Form{}, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return f, nil
}

// DeleteForm deletes a Form
func DeleteForm(c *web.AppContext) *web.AppError {

	idStr := c.GetValue("id")
	//we must have an id to delete
	if idStr == "" {
		message := fmt.Sprintf("Cannot delete Form. Invalid Id [%s]", idStr)
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//delete
	if err := c.MDB.DB.C(model.Forms).RemoveId(idStr); err != nil {
		message := fmt.Sprintf("Error deleting Form [%v]", idStr)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
