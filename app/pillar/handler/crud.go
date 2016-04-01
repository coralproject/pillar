package handler

import (
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
)

//CreateUpdateAsset end-point allows creation or updation of Asset.
func CreateUpdateAsset(c *web.AppContext) {
	dbObject, err := service.CreateUpdateAsset(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateUser end-point allows creation or updation of User.
func CreateUpdateUser(c *web.AppContext) {
	dbObject, err := service.CreateUpdateUser(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateComment end-point allows creation or updation of Comment.
func CreateUpdateComment(c *web.AppContext) {
	dbObject, err := service.CreateUpdateComment(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateSection end-point allows creation or updation of Section.
func CreateUpdateSection(c *web.AppContext) {
	dbObject, err := service.CreateUpdateSection(c)
	doRespond(c, dbObject, err)
}

//CreateUpdateAuthor end-point allows creation or updation of Author.
func CreateUpdateAuthor(c *web.AppContext) {
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
	dbObject, err := service.CreateUpdateSearch(c)
	doRespond(c, dbObject, err)
}

//DeleteSearch end-point allows deletion of Search.
func DeleteSearch(c *web.AppContext) {
	err := service.DeleteSearch(c)
	doRespond(c, nil, err)
}