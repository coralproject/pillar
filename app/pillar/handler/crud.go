package handler

import (
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
)

//CreateUpdateAsset end-point allows creation or updation of Asset.
func CreateUpdateAsset(c *web.AppContext) {
	c.Event = model.EventAssetAddUpdate
	dbObject, err := service.CreateUpdateAsset(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateUser end-point allows creation or updation of User.
func CreateUpdateUser(c *web.AppContext) {
	c.Event = model.EventUserAddUpdate
	dbObject, err := service.CreateUpdateUser(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateComment end-point allows creation or updation of Comment.
func CreateUpdateComment(c *web.AppContext) {
	c.Event = model.EventCommentAddUpdate
	dbObject, err := service.CreateUpdateComment(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateSection end-point allows creation or updation of Section.
func CreateUpdateSection(c *web.AppContext) {
	c.Event = model.EventSectionAddUpdate
	dbObject, err := service.CreateUpdateSection(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateAuthor end-point allows creation or updation of Author.
func CreateUpdateAuthor(c *web.AppContext) {
	c.Event = model.EventAuthorAddUpdate
	dbObject, err := service.CreateUpdateAuthor(c)
	doRespond(c, dbObject, err)
}

//CreateUserAction end-point allows registration of all events from front-end.
func CreateUserAction(c *web.AppContext) {
	err := service.CreateUserAction(c)
	doRespond(c, nil, err)
}

//CreateIndex end-point allows creation of new database index.
func CreateIndex(c *web.AppContext) {
	err := service.CreateIndex(c)
	doRespond(c, nil, err)
}

//UpdateMetadata end-point allows updation of Metadata within an entity.
func UpdateMetadata(c *web.AppContext) {
	dbObject, err := service.UpdateMetadata(c)
	doRespond(c, dbObject, err)
}

//GetTags end-point returns all available tags in the system.
func GetTags(c *web.AppContext) {
	dbObject, err := service.GetTags(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateTag end-point allows creation or updation of Tag.
func CreateUpdateTag(c *web.AppContext) {
	dbObject, err := service.CreateUpdateTag(c)
	doRespond(c, dbObject, err)
}

//DeleteTag end-point allows deletion of Tag.
func DeleteTag(c *web.AppContext) {
	err := service.DeleteTag(c)
	doRespond(c, nil, err)
}

//GetSearches end-point returns all available Searches in the system.
func GetSearches(c *web.AppContext) {
	dbObject, err := service.GetSearches(c)
	doRespond(c, dbObject, err)
}

//GetSearch end-point returns a single Search.
func GetSearch(c *web.AppContext) {
	dbObject, err := service.GetSearch(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateSearch end-point allows creation or updation of Search.
func CreateUpdateSearch(c *web.AppContext) {
	c.Event = model.EventSearchAddUpdate
	dbObject, err := service.CreateUpdateSearch(c)
	doRespond(c, dbObject, err)
}

//DeleteSearch end-point allows deletion of Search.
func DeleteSearch(c *web.AppContext) {
	err := service.DeleteSearch(c)
	doRespond(c, nil, err)
}

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

func AddAnswerToFormGallery(c *web.AppContext) {
	dbObject, err := service.AddAnswerToFormGallery(c)
	doRespond(c, dbObject, err)
}

func RemoveAnswerFromFormGallery(c *web.AppContext) {
	dbObject, err := service.AddAnswerToFormGallery(c)
	doRespond(c, dbObject, err)
}
