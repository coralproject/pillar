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
	c.SD.Client.Inc("Create_Update_Asset", 1, 1.0)
}

//CreateUpdateUser end-point allows creation or updation of User.
func CreateUpdateUser(c *web.AppContext) {
	c.Event = model.EventUserAddUpdate
	dbObject, err := service.CreateUpdateUser(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Create_Update_User", 1, 1.0)
}

//CreateUpdateComment end-point allows creation or updation of Comment.
func CreateUpdateComment(c *web.AppContext) {
	c.Event = model.EventCommentAddUpdate
	dbObject, err := service.CreateUpdateComment(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Create_Update_Comment", 1, 1.0)
}

//CreateUpdateSection end-point allows creation or updation of Section.
func CreateUpdateSection(c *web.AppContext) {
	c.Event = model.EventSectionAddUpdate
	dbObject, err := service.CreateUpdateSection(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Create_Update_Section", 1, 1.0)
}

//CreateUpdateAuthor end-point allows creation or updation of Author.
func CreateUpdateAuthor(c *web.AppContext) {
	c.Event = model.EventAuthorAddUpdate
	dbObject, err := service.CreateUpdateAuthor(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Create_Update_Author", 1, 1.0)
}

//CreateUserAction end-point allows registration of all events from front-end.
func CreateUserAction(c *web.AppContext) {
	err := service.CreateUserAction(c)
	doRespond(c, nil, err)
	c.SD.Client.Inc("Create_User_Action", 1, 1.0)
}

//CreateIndex end-point allows creation of new database index.
func CreateIndex(c *web.AppContext) {
	err := service.CreateIndex(c)
	doRespond(c, nil, err)
	c.SD.Client.Inc("Create_Index", 1, 1.0)
}

//UpdateMetadata end-point allows updation of Metadata within an entity.
func UpdateMetadata(c *web.AppContext) {
	dbObject, err := service.UpdateMetadata(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Update_Metadata", 1, 1.0)
}

//GetTags end-point returns all available tags in the system.
func GetTags(c *web.AppContext) {
	dbObject, err := service.GetTags(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Get_Tags", 1, 1.0)
}

//CreateUpdateTag end-point allows creation or updation of Tag.
func CreateUpdateTag(c *web.AppContext) {
	dbObject, err := service.CreateUpdateTag(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Create_Update_Tag", 1, 1.0)
}

//DeleteTag end-point allows deletion of Tag.
func DeleteTag(c *web.AppContext) {
	err := service.DeleteTag(c)
	doRespond(c, nil, err)
	c.SD.Client.Inc("Delete_Tag", 1, 1.0)
}

//GetSearches end-point returns all available Searches in the system.
func GetSearches(c *web.AppContext) {
	dbObject, err := service.GetSearches(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Get_Searches", 1, 1.0)
}

//GetSearch end-point returns a single Search.
func GetSearch(c *web.AppContext) {
	dbObject, err := service.GetSearch(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Get_Search", 1, 1.0)
}

//CreateUpdateSearch end-point allows creation or updation of Search.
func CreateUpdateSearch(c *web.AppContext) {
	c.Event = model.EventSearchAddUpdate
	dbObject, err := service.CreateUpdateSearch(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Create_Update_Search", 1, 1.0)
}

//DeleteSearch end-point allows deletion of Search.
func DeleteSearch(c *web.AppContext) {
	err := service.DeleteSearch(c)
	doRespond(c, nil, err)
	c.SD.Client.Inc("Delete_Search", 1, 1.0)
}
