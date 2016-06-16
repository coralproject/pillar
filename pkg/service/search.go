package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
)

// GetSearches returns the list of all Search items in the system
func GetSearches(context *web.AppContext) ([]model.Search, *web.AppError) {

	all := make([]model.Search, 0)
	if err := context.MDB.DB.C(model.Searches).Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching searches")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}

// GetSearch returns one Search
func GetSearch(context *web.AppContext) (*model.Search, *web.AppError) {

	idStr := context.GetValue("id")
	//we must have an id to delete the search
	if idStr == "" {
		message := fmt.Sprintf("Cannot fetch Search. Invalid Id [%s]", idStr)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//convert to an ObjectId
	id := bson.ObjectIdHex(idStr)

	var search model.Search
	if err := context.MDB.DB.C(model.Searches).FindId(id).One(&search); err != nil {
		message := fmt.Sprintf("Error fetching one Search [%s]", err)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return &search, nil
}

// CreateUpdateSearch upserts a Search
func CreateUpdateSearch(context *web.AppContext) (*model.Search, *web.AppError) {
	var input model.Search
	if err := UnmarshallAndValidate(context, &input); err != nil {
		return nil, err
	}

	var dbEntity model.Search
	//Upsert if entity exists with same ID
	context.MDB.DB.C(model.Searches).FindId(input.ID).One(&dbEntity)

	if dbEntity.ID == "" { //new
		input.ID = bson.NewObjectId()
		input.Query = strings.Join([]string{"trust_search_", input.ID.Hex()}, "")
		input.DateCreated = time.Now()
	} else { //existing
		input.DateUpdated = time.Now()
	}

	if _, err := context.MDB.DB.C(model.Searches).UpsertId(input.ID, &input); err != nil {
		fmt.Printf("Error: %s", err)
		message := fmt.Sprintf("Error updating existing Search [%v]", input)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	//save an entry in history
	createSearchHistory(context, input)

	return &input, nil
}

func createSearchHistory(context *web.AppContext, search model.Search) {
	var sh model.SearchHistory
	sh.ID = bson.NewObjectId()
	if search.DateUpdated.IsZero() {
		sh.Action = "create"
	} else {
		sh.Action = "update"
	}
	sh.Date = time.Now()
	sh.Search = search
	if err := context.MDB.DB.C(model.SrchHistory).Insert(sh); err != nil {
		log.Printf("Error creating SearchHistory [%s]", err)
	}
}

// DeleteSearch deletes a Search
func DeleteSearch(context *web.AppContext) *web.AppError {

	idStr := context.GetValue("id")
	//we must have an id to delete the search
	if idStr == "" {
		message := fmt.Sprintf("Cannot delete Search. Invalid Id [%s]", idStr)
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//convert to an ObjectId
	id := bson.ObjectIdHex(idStr)

	//delete
	if err := context.MDB.DB.C(model.Searches).RemoveId(id); err != nil {
		message := fmt.Sprintf("Error deleting Search [%s]", id)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
