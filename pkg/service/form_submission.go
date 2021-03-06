package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
)

func EditFormSubmissionAnswer(c *web.AppContext) (*model.FormSubmission, *web.AppError) {

	// get our tasty form submission
	s, err := GetFormSubmission(c)
	if err != nil {
		return nil, &web.AppError{nil, "Could not edit submission answer: form submission not found", http.StatusInternalServerError}

	}

	// look for the answer in question
	for i, a := range s.Answers {

		if a.WidgetId == c.GetValue("answer_id") {

			body := model.FormSubmissionEditInput{}
			_ = c.Unmarshall(&body)

			s.Answers[i].EditedAnswer = body.EditedAnswer

		}
	}

	// do the update
	q := bson.M{"_id": s.ID}
	appErr := c.MDB.DB.C(model.FormSubmissions).Update(q, s)
	if appErr != nil {
		message := fmt.Sprintf("Error updating Form Submission after edit")
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	return &s, nil

}

func buildSubmissionFromForm(f model.Form) model.FormSubmission {

	// cook up a new form submission
	fs := model.FormSubmission{}

	// grab the header info from the form
	fs.FormId = f.ID
	fs.Header = f.Header
	fs.Footer = f.Footer

	// for each widget in each step
	for _, s := range f.Steps {
		for _, w := range s.Widgets {

			// make an answer
			a := model.FormSubmissionAnswer{}

			// get the question/title and props for posterity
			a.WidgetId = w.ID
			a.Identity = w.Identity
			a.Question = w.Title
			a.Props = w.Props

			// and slam them into the answers
			fs.Answers = append(fs.Answers, a)

		}
	}

	// toss that fresh submission back
	return fs

}

//  ** consider implementing this as a method on FormSubmission **
// it's a little peculiar:
// each submission to a Form will have a record for every answer no
// matter what the fe sends
// these are prepopulated by buildSubmissionFromForm above
// so..
func setAnswersToFormSubmission(fs model.FormSubmission, fsi model.FormSubmissionInput) model.FormSubmission {

	// for each answer inputted
	for _, ai := range fsi.Answers {

		// look for the answer
		for x, a := range fs.Answers {

			// add the answer to the appropriate spot
			if a.WidgetId == ai.WidgetId {
				fs.Answers[x].Answer = ai.Answer
			}

		}

	}

	return fs

}

func CreateFormSubmission(context *web.AppContext) (*model.FormSubmission, *web.AppError) {

	// we take an input type here as what's passed needs some work
	//  before it's a proper submission
	var input model.FormSubmissionInput
	if err := UnmarshallAndValidate(context, &input); err != nil {
		return nil, err
	}

	/* Todo, custom validation
	if input.Name == "" {
		message := fmt.Sprintf("Invalid Section Name [%s]", input.Name)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}
	*/

	// get the form id from the context
	fId := bson.ObjectIdHex(context.GetValue("form_id"))

	// create a context to get the form
	fc := web.NewContext(nil, nil)
	defer fc.Close()
	fc.SetValue("id", fId.Hex())

	// get the form in question
	f, err := GetForm(fc)
	if err != nil {
		return nil, err
	}

	// build a form submission from the input
	fs := buildSubmissionFromForm(f)

	// set the answers into the submission
	fs = setAnswersToFormSubmission(fs, input)

	// set miscellenia
	fs.DateCreated = time.Now()
	fs.DateUpdated = time.Now()

	// set the number
	fs.Number = getSubmissionCountByForm(fc) + 1

	// aaaand save it
	fs.ID = bson.NewObjectId()
	if err := context.MDB.DB.C(model.FormSubmissions).Insert(fs); err != nil {
		message := fmt.Sprintf("Error inserting FormSubmission")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	// update the stats using the Form Context
	err = updateStats(fc)
	if err != nil {
		return nil, err
	}

	return &fs, nil

}

// GetFormSubmissions returns an array of FormSubmissions
func GetFormSubmissionsByForm(c *web.AppContext) (map[string]interface{}, *web.AppError) {

	result := make(map[string]interface{}, 2)

	/* Get form submissions */

	idStr := c.GetValue("form_id")
	//we must have an id to delete the search
	if idStr == "" {
		message := fmt.Sprintf("Cannot get FormSubmissions. Invalid Id [%s]", idStr)
		return result, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	limit, err := strconv.Atoi(c.GetValue("limit"))
	if err != nil {
		limit = 0
	}

	skip, err := strconv.Atoi(c.GetValue("skip"))
	if err != nil {
		skip = 0
	}

	// always order by date, asc or desc
	orderby := c.GetValue("orderby")
	if orderby == "dsc" {
		orderby = "-date_created"
	} else {
		orderby = "date_created"
	}

	// convert to an ObjectId
	id := bson.ObjectIdHex(idStr)

	// create find query
	find := bson.M{"form_id": id}

	// get what we are searching for
	search := c.GetValue("search")

	if search != "" {
		// ensure that the text index is created
		index := mgo.Index{
			Key:        []string{"$text:replies.answer"},
			Unique:     false,
			DropDups:   false,
			Background: false,
			Sparse:     true,
			Name:       "replies.answer.text",
		}
		_ = c.MDB.DB.C(model.FormSubmissions).EnsureIndex(index)

		// add a text search to the query
		find["$text"] = bson.M{"$search": search}
	}

	searchquery := c.MDB.DB.C(model.FormSubmissions).Find(find)

	/* Calculate totals */
	var onlysearchfss []model.FormSubmission
	if e := searchquery.All(&onlysearchfss); e != nil {
		message := fmt.Sprintf("Error fetching FormSubmissions")
		return result, &web.AppError{e, message, http.StatusInternalServerError}
	}

	// get totals for the specific search
	if result["counts"], err = getTotals(onlysearchfss); err != nil {
		message := fmt.Sprintf("Error calculating totals for FormSubmissions")
		return result, &web.AppError{err, message, http.StatusInternalServerError}
	}

	// get all the form submissions in the system
	if result["counts"].(map[string]interface{})["total_submissions"], err = c.MDB.DB.C(model.FormSubmissions).Find(bson.M{"form_id": id}).Count(); err != nil {
		message := fmt.Sprintf("Error count all FormSubmissions")
		return result, &web.AppError{err, message, http.StatusInternalServerError}
	}

	/* Get the form submissions filter by flag */

	// filter by flagged, bookmarked or any other flags (or not)
	flag := c.GetValue("filterby")

	// build the query to filter by flag
	if flag != "" {
		// we are using -flag to bring all the submissions that do not contain that flag
		if strings.HasPrefix(flag, "-") {
			notflag := strings.TrimLeft(flag, "-")
			find["flags"] = bson.M{"$not": bson.M{"$elemMatch": bson.M{"$in": []string{notflag}}}}
		} else {
			find["flags"] = bson.M{"$regex": flag} // filterby flag
		}
	}

	query := c.MDB.DB.C(model.FormSubmissions).Find(find).Skip(skip).Limit(limit).Sort(orderby)

	var fss []model.FormSubmission
	if e := query.All(&fss); e != nil {
		message := fmt.Sprintf("Error fetching FormSubmissions")
		return result, &web.AppError{e, message, http.StatusInternalServerError}
	}

	result["submissions"] = fss
	return result, nil
}

func getTotals(fss []model.FormSubmission) (map[string]interface{}, error) {
	totals := make(map[string]interface{}, 0)
	var err error

	totalsearch := 0
	totalsperflag := make(map[string]int, 0)
	// loop through the submissions to get total per flag
	for _, f := range fss {
		// count on the search per flag
		for _, flag := range f.Flags {
			totalsperflag[flag] = totalsperflag[flag] + 1
		}
		totalsearch = totalsearch + 1
	}

	totals["search_by_flag"] = totalsperflag
	totals["total_search"] = totalsearch

	return totals, err
}

// GetFormSubmissions returns a single FormSubmission by id
func GetFormSubmission(c *web.AppContext) (model.FormSubmission, *web.AppError) {

	idStr := c.GetValue("id")
	//we must have an id to delete the search

	if idStr == "" {
		message := fmt.Sprintf("Cannot get FormSubmission. Invalid Id [%s]", idStr)
		return model.FormSubmission{}, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//convert to an ObjectId
	id := bson.ObjectIdHex(idStr)

	f := model.FormSubmission{}
	if err := c.MDB.DB.C(model.FormSubmissions).Find(bson.M{"_id": id}).One(&f); err != nil {
		message := fmt.Sprintf("Error fetching FormSubmissions")
		return model.FormSubmission{}, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return f, nil
}

// DeleteFormSubmission deletes a FormSubmission
func DeleteFormSubmission(c *web.AppContext) *web.AppError {

	idStr := c.GetValue("id")
	//we must have an id to delete the search
	if idStr == "" {
		message := fmt.Sprintf("Cannot delete FormSubmission. Invalid Id [%s]", idStr)
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	id := bson.ObjectIdHex(idStr)

	//delete
	if err := c.MDB.DB.C(model.FormSubmissions).RemoveId(id); err != nil {
		message := fmt.Sprintf("Error deleting FormSubmission [%v]", idStr)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}

// given a form's id and a stats, update the form with the status
func UpdateFormSubmissionStatus(context *web.AppContext) (*model.FormSubmission, *web.AppError) {

	// todo, gracefully message invalid ids
	id := bson.ObjectIdHex(context.GetValue("id"))
	status := context.GetValue("status")

	// let's make sure we don't update all of them..
	q := bson.M{"_id": id}
	s := bson.M{"$set": bson.M{"status": status, "date_updated": time.Now()}}

	// do the update
	err := context.MDB.DB.C(model.FormSubmissions).Update(q, s)
	if err != nil {
		message := fmt.Sprintf("Error updating Form Submission status")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	var f *model.FormSubmission
	err = context.MDB.DB.C(model.FormSubmissions).FindId(id).One(&f)
	if err != nil {
		message := fmt.Sprintf("Could not find Form Submission ", id)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return f, nil

}

/*  Flag functionality specified here can be abstracted as Flaggabe behavior */
func RemoveFlagFromFormSubmission(context *web.AppContext) (*model.FormSubmission, *web.AppError) {

	// get our tasty form submission
	s, err := GetFormSubmission(context)
	if err != nil {
		return nil, &web.AppError{nil, "Could not edit submission answer: form submission not found", http.StatusInternalServerError}
	}

	fi := -1 // a var to store the flag's index
	f := context.GetValue("flag")

	// find the flag
	for i, tf := range s.Flags {
		if tf == f {
			fi = i
			break
		}
	}

	// slice that flag out
	if fi != -1 {
		s.Flags = append(s.Flags[:fi], s.Flags[fi+1:]...)
	}

	// let's make sure we don't update all of them..
	q := bson.M{"_id": s.ID}
	u := bson.M{"$set": bson.M{"flags": s.Flags, "date_updated": time.Now()}}

	// do the update
	err2 := context.MDB.DB.C(model.FormSubmissions).Update(q, u)
	if err2 != nil {
		message := fmt.Sprintf("Error updating Form Submission after removing flag")
		return nil, &web.AppError{err2, message, http.StatusInternalServerError}
	}

	return &s, nil

}

func AddFlagToFormSubmission(context *web.AppContext) (*model.FormSubmission, *web.AppError) {

	// get our tasty form submission
	s, err := GetFormSubmission(context)
	if err != nil {
		return nil, &web.AppError{nil, "Could not edit submission answer: form submission not found", http.StatusInternalServerError}
	}

	fi := -1 // a var to store the flag's index
	f := context.GetValue("flag")

	// find the flag
	for i, tf := range s.Flags {
		if tf == f {
			fi = i
			break
		}
	}

	// if it's not there, add it
	if fi == -1 {
		s.Flags = append(s.Flags, f)
	}

	// let's make sure we don't update all of them..
	q := bson.M{"_id": s.ID}
	u := bson.M{"$set": bson.M{"flags": s.Flags, "date_updated": time.Now()}}

	// do the update
	err2 := context.MDB.DB.C(model.FormSubmissions).Update(q, u)
	if err2 != nil {
		message := fmt.Sprintf("Error updating Form Submission after removing flag")
		return nil, &web.AppError{err2, message, http.StatusInternalServerError}
	}

	return &s, nil

}

func SearchFormSubmissions(c *web.AppContext) ([]model.FormSubmission, *web.AppError) {

	// ensure that the text index is created
	index := mgo.Index{
		Key:        []string{"$text:replies.answer"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     true,
		Name:       "replies.answer.text",
	}
	_ = c.MDB.DB.C(model.FormSubmissions).EnsureIndex(index)

	// get what we are searching for
	s := c.GetValue("search")

	// we are building a text search query
	// go find it!
	var fss []model.FormSubmission
	q := c.MDB.DB.C(model.FormSubmissions).Find(bson.M{"$text": bson.M{"$search": s}})

	if err := q.All(&fss); err != nil {
		message := fmt.Sprintf("Error searching Form Submissions by %s", s)
		werr := web.AppError{err, message, http.StatusInternalServerError}
		return nil, &werr
	}
	return fss, nil
}
