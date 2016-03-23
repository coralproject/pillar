package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// GetUserGroups returns the list of all UserGroups in the system
func GetUserGroups(context *web.AppContext) ([]model.UserGroup, *web.AppError) {

	all := make([]model.UserGroup, 0)
	if err := context.DB.UserGroups.Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching tags")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}

// CreateUpdateUserGroup adds or updates a UserGroup
func CreateUpdateUserGroup(context *web.AppContext) (*model.UserGroup, *web.AppError) {
	var input model.UserGroup
	context.Unmarshall(&input)

	fmt.Printf("UserGroup: %v", input)
	var dbEntity model.UserGroup
	//Upsert if entity exists with same ID
	context.DB.UserGroups.Find(input.ID).One(&dbEntity)
	if dbEntity.ID == "" { //new
		input.ID = bson.NewObjectId()
		input.DateCreated = time.Now()
	} else { //existing
		input.DateUpdated = time.Now()
	}

	if _, err := context.DB.UserGroups.UpsertId(input.ID, &input); err != nil {
		fmt.Printf("Error: %s", err)
		message := fmt.Sprintf("Error updating existing UserGroup [%v]", input)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return &input, nil
}

// DeleteUserGroup deletes a UserGroup
func DeleteUserGroup(context *web.AppContext) *web.AppError {
	var input model.UserGroup
	context.Unmarshall(&input)

	//we must have the tag name for deletion
	if input.ID == "" {
		message := fmt.Sprintf("Cannot delete UserGroup with an empty ID [%v]", input)
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//delete
	if err := context.DB.UserGroups.RemoveId(input.ID); err != nil {
		message := fmt.Sprintf("Error deleting UserGroup [%v]", input)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
