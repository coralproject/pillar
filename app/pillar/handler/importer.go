package handler

import (
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
	"net/http"
)

//ImportUser imports a new user to the system
func ImportUser(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.ImportUser(c)
	doRespond(w, dbObject, err)
}

//ImportAsset imports a new asset to the system
func ImportAsset(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.ImportAsset(c)
	doRespond(w, dbObject, err)
}

//ImportComment imports a new comment to the system
func ImportComment(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.ImportComment(c)
	doRespond(w, dbObject, err)
}

//ImportAction imports actions into the system
func ImportAction(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.ImportAction(c)
	doRespond(w, dbObject, err)
}

//ImportNote imports notes into the system
func ImportNote(w http.ResponseWriter, r *http.Request, c *web.AppContext) {
	dbObject, err := service.CreateNote(c)
	doRespond(w, dbObject, err)
}
