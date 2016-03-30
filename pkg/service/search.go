package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// GetSearches returns the list of all Search items in the system
func GetSearches(context *web.AppContext) ([]model.Search, *web.AppError) {

	all := make([]model.Search, 0)
	if err := context.DB.Searches.Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching searches")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}

// CreateUpdateSearch upserts a Search
func CreateUpdateSearch(context *web.AppContext) (*model.Search, *web.AppError) {
	var input model.Search
	context.Unmarshall(&input)

	fmt.Printf("Search: %v", input)
	var dbEntity model.Search
	//Upsert if entity exists with same ID
	context.DB.Searches.Find(input.ID).One(&dbEntity)
	if dbEntity.ID == "" { //new
		input.ID = bson.NewObjectId()
		input.DateCreated = time.Now()
	} else { //existing
		input.DateUpdated = time.Now()
	}

	if _, err := context.DB.Searches.UpsertId(input.ID, &input); err != nil {
		fmt.Printf("Error: %s", err)
		message := fmt.Sprintf("Error updating existing Search [%v]", input)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return &input, nil
}

// DeleteSearch deletes a Search
func DeleteSearch(context *web.AppContext) *web.AppError {
	var input model.Search
	context.Unmarshall(&input)

	//we must have the tag name for deletion
	if input.ID == "" {
		message := fmt.Sprintf("Cannot delete Search with an empty ID [%v]", input)
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//delete
	if err := context.DB.Searches.RemoveId(input.ID); err != nil {
		message := fmt.Sprintf("Error deleting Search [%v]", input)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
