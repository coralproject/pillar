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
