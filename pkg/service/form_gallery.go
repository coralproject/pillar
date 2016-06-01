package service

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
)

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
