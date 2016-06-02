package handler

/*

	This file contains handlers for web endpoint
	invocations of Ask services:

	* Forms
	* Form Submissions
	* Form Galleries

*/

import (
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
)

func CreateUpdateForm(c *web.AppContext) {
	dbObject, err := service.CreateUpdateForm(c)
	doRespond(c, dbObject, err)
}

func UpdateFormStatus(c *web.AppContext) {
	dbObject, err := service.UpdateFormStatus(c)
	doRespond(c, dbObject, err)
}

func GetForms(c *web.AppContext) {
	dbObject, err := service.GetForms(c)
	doRespond(c, dbObject, err)
}

func GetForm(c *web.AppContext) {
	dbObject, err := service.GetForm(c)
	doRespond(c, dbObject, err)
}

func DeleteForm(c *web.AppContext) {
	err := service.DeleteForm(c)
	doRespond(c, nil, err)
}

func CreateFormSubmission(c *web.AppContext) {
	dbObject, err := service.CreateFormSubmission(c)
	doRespond(c, dbObject, err)
}

func UpdateFormSubmissionStatus(c *web.AppContext) {
	dbObject, err := service.UpdateFormSubmissionStatus(c)
	doRespond(c, dbObject, err)
}

func EditFormSubmissionAnswer(c *web.AppContext) {
	dbObject, err := service.EditFormSubmissionAnswer(c)
	doRespond(c, dbObject, err)
}

func GetFormSubmissionsByForm(c *web.AppContext) {
	dbObject, err := service.GetFormSubmissionsByForm(c)
	doRespond(c, dbObject, err)
}

func GetFormSubmission(c *web.AppContext) {
	dbObject, err := service.GetFormSubmission(c)
	doRespond(c, dbObject, err)
}

func DeleteFormSubmission(c *web.AppContext) {
	err := service.DeleteSearch(c)
	doRespond(c, nil, err)
}

func AddFlagToFormSubmission(c *web.AppContext) {
	dbObject, err := service.AddFlagToFormSubmission(c)
	doRespond(c, dbObject, err)
}

func RemoveFlagFromFormSubmission(c *web.AppContext) {
	dbObject, err := service.RemoveFlagFromFormSubmission(c)
	doRespond(c, dbObject, err)
}

func AddAnswerToFormGallery(c *web.AppContext) {
	dbObject, err := service.AddAnswerToFormGallery(c)
	doRespond(c, dbObject, err)
}

func RemoveAnswerFromFormGallery(c *web.AppContext) {
	dbObject, err := service.AddAnswerToFormGallery(c)
	doRespond(c, dbObject, err)
}

func GetFormGalleriesByForm(c *web.AppContext) {
	dbObject, err := service.GetFormGalleriesByForm(c)
	doRespond(c, dbObject, err)
}

func GetFormGallery(c *web.AppContext) {
	dbObject, err := service.GetFormGallery(c)
	doRespond(c, dbObject, err)
}
