package service

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
)

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

			fmt.Println(message, err)

			return nil, &web.AppError{err, message, http.StatusInternalServerError}
		}
	}

	// update
	if _, err := context.MDB.DB.C(model.Forms).UpsertId(input.ID, input); err != nil {
		message := fmt.Sprintf("Error creating/updating Form")

		return nil, &web.AppError{err, message, http.StatusInternalServerError}
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
