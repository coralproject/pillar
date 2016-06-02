package service

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
)

func galleryContainsSubmissionAnswer(g *model.FormGallery, a *model.FormGalleryAnswer) bool {

	for _, i := range g.Answers {
		if i.SubmissionId == a.SubmissionId && i.AnswerId == a.AnswerId {
			return true
		}
	}

	return false
}

// add an answer to a form gallery
func AddAnswerToFormGallery(context *web.AppContext) (*model.FormGallery, *web.AppError) {

	// get the form gallery in question
	g, err := GetFormGallery(context)
	if err != nil {
		message := fmt.Sprintf("Cannot add answer to form gallery: form not found %s", context.GetValue("id"))
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	// grab the ids of the submission and answer we're adding
	sId := bson.ObjectIdHex(context.GetValue("submission_id"))
	aId := context.GetValue("answer_id") // answer ids are not bson

	// make ourselves a FormGalleryAnswer
	a := model.FormGalleryAnswer{SubmissionId: sId, AnswerId: aId}

	// make sure it's not already there
	if galleryContainsSubmissionAnswer(&g, &a) {
		message := fmt.Sprintf("Cannot add answer to form gallery: already exists", context.GetValue("id"))
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	// append the answer
	g.Answers = append(g.Answers, a)
	if err := context.MDB.DB.C(model.FormGalleries).Update(bson.M{"_id": g.ID}, g); err != nil {
		message := fmt.Sprintf("Cannot add answer to gallery, Error updating FormGallery")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return &g, nil

}

// remove an answer from a form gallery
func RemoveAnswerFromFormGallery(context *web.AppContext) (*model.FormGallery, *web.AppError) {

	// get the form gallery in question
	g, err := GetFormGallery(context)
	if err != nil {
		message := fmt.Sprintf("Cannot remove answer from form gallery: form not found %s", context.GetValue("id"))
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	// grab the ids of the submission and answer we're adding
	sId := bson.ObjectIdHex(context.GetValue("submission_id"))
	aId := context.GetValue("answer_id") // answer ids are not bson

	// make ourselves a FormGalleryAnswer
	a := model.FormGalleryAnswer{SubmissionId: sId, AnswerId: aId}

	// make sure it's not already there
	if !galleryContainsSubmissionAnswer(&g, &a) {
		message := fmt.Sprintf("Cannot remove answer from form gallery: answer not present in %s", context.GetValue("id"))
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	// find the index
	index := 0
	for i, ia := range g.Answers {
		if ia.SubmissionId == a.SubmissionId && ia.AnswerId == a.AnswerId {
			index = i
			break
		}
	}

	// cut the element from the answers by index
	g.Answers = append(g.Answers[:index], g.Answers[index+1:]...)

	// save it
	if err := context.MDB.DB.C(model.FormGalleries).Update(bson.M{"_id": g.ID}, g); err != nil {
		message := fmt.Sprintf("Cannot add answer to gallery, Error updating FormGallery")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return &g, nil

}

//  ** consider implementing this as a method on FormGallery **
func CreateFormGallery(context *web.AppContext) (*model.FormGallery, *web.AppError) {

	// get the form id from the context
	fId := bson.ObjectIdHex(context.GetValue("form_id"))
	if fId == "" {
		message := fmt.Sprintf("Cannot create FormGallery: form_id not provided")
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	// create a new gallery and set it up
	fg := model.FormGallery{
		FormId:      fId,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	// aaaand save it
	fg.ID = bson.NewObjectId()
	if err := context.MDB.DB.C(model.FormGalleries).Insert(fg); err != nil {
		message := fmt.Sprintf("Error inserting FormGallery")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return &fg, nil

}

// embeds the latest version of the FormSubmisison.Answer into
//  a Form Gallery.  Loaded every time to react to Edits/deltes
//  of form submission content
func hydrateFormGallery(g model.FormGallery) model.FormGallery {

	// get a context to load the submissions
	c := web.NewContext(nil, nil)

	// for each answer in the gallery
	for _, a := range g.Answers {

		// load the submission
		c.SetValue("id", a.SubmissionId.Hex())
		s, err := GetFormSubmission(c)
		if err != nil {
			// remove answers from gallery if submission is
			//  deleted?
		}

		// find the answer
		for i, fsa := range s.Answers {
			if fsa.WidgetId == a.AnswerId {

				// and embed it into the form gallery
				g.Answers[i].Answer = fsa
				break
			}
		}
	}

	return g
}

// GetFormGallerys returns an array of FormGallerys
func GetFormGalleriesByForm(c *web.AppContext) ([]model.FormGallery, *web.AppError) {

	idStr := c.GetValue("form_id")
	//we must have an id to delete the search
	if idStr == "" {
		message := fmt.Sprintf("Cannot get FormGalleries. Invalid Id [%s]", idStr)
		return []model.FormGallery{}, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//convert to an ObjectId
	id := bson.ObjectIdHex(idStr)
	fss := make([]model.FormGallery, 0)
	if err := c.MDB.DB.C(model.FormGalleries).Find(bson.M{"form_id": id}).All(&fss); err != nil {
		message := fmt.Sprintf("Error fetching FormGallerys")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	//hydrate them all...
	for i, g := range fss {
		fss[i] = hydrateFormGallery(g)
	}

	return fss, nil
}

// GetFormGallerys returns a single FormGallery by id
func GetFormGallery(c *web.AppContext) (model.FormGallery, *web.AppError) {

	idStr := c.GetValue("id")
	//we must have an id to delete the search
	if idStr == "" {
		message := fmt.Sprintf("Cannot get FormGallery. Invalid Id [%s]", idStr)
		return model.FormGallery{}, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//convert to an ObjectId
	id := bson.ObjectIdHex(idStr)

	f := model.FormGallery{}
	if err := c.MDB.DB.C(model.FormGalleries).Find(bson.M{"_id": id}).One(&f); err != nil {
		message := fmt.Sprintf("Error fetching FormGalleries")
		return model.FormGallery{}, &web.AppError{err, message, http.StatusInternalServerError}
	}

	// hydrate the form gallery
	f = hydrateFormGallery(f)

	return f, nil
}

// DeleteFormGallery deletes a FormGallery
func DeleteFormGallery(c *web.AppContext) *web.AppError {

	idStr := c.GetValue("id")
	//we must have an id to delete the search
	if idStr == "" {
		message := fmt.Sprintf("Cannot delete FormGallery. Invalid Id [%s]", idStr)
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//delete
	if err := c.MDB.DB.C(model.FormGalleries).RemoveId(idStr); err != nil {
		message := fmt.Sprintf("Error deleting FormGallery [%v]", idStr)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
