package handler

import (
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
	"net/http"
)

//CreateUpdateAsset end-point allows creation or updation of Asset.
func CreateUpdateAsset(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.CreateUpdateAsset(c)
	doRespond(w, dbObject, err)
}

//CreateUpdateUser end-point allows creation or updation of User.
func CreateUpdateUser(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.CreateUpdateUser(c)
	doRespond(w, dbObject, err)
}

//CreateUpdateComment end-point allows creation or updation of Comment.
func CreateUpdateComment(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.CreateUpdateComment(c)
	doRespond(w, dbObject, err)
}

//CreateUpdateSection end-point allows creation or updation of Section.
func CreateUpdateSection(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.CreateUpdateSection(c)
	doRespond(w, dbObject, err)
}

//CreateUpdateAuthor end-point allows creation or updation of Author.
func CreateUpdateAuthor(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.CreateUpdateAuthor(c)
	doRespond(w, dbObject, err)
}

//CreateUserAction end-point allows registration of all events from front-end.
func CreateUserAction(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	err := service.CreateUserAction(c)
	doRespond(w, nil, err)
}

//CreateIndex end-point allows creation of new database index.
func CreateIndex(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	err := service.CreateIndex(c)
	doRespond(w, nil, err)
}

//UpdateMetadata end-point allows updation of Metadata within an entity.
func UpdateMetadata(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.UpdateMetadata(c)
	doRespond(w, dbObject, err)
}

//GetTags end-point returns all available tags in the system.
func GetTags(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.GetTags(c)
	doRespond(w, dbObject, err)
}

//CreateUpdateTag end-point allows creation or updation of Tag.
func CreateUpdateTag(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.CreateUpdateTag(c)
	doRespond(w, dbObject, err)
}

//DeleteTag end-point allows deletion of Tag.
func DeleteTag(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	err := service.DeleteTag(c)
	doRespond(w, nil, err)
}
