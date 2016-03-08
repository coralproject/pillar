package handler

import (
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
)

//CreateUpdateUser creates/updates an user in the system
func CreateUpdateUser(w http.ResponseWriter, r *http.Request, c *service.AppContext) {
	dbObject, err := service.CreateUpdateUser(c)
	doRespond(w, dbObject, err)
}

//CreateUpdateAsset creates/updates an asset in the system
func CreateUpdateAsset(w http.ResponseWriter, r *http.Request, c *service.AppContext) {
	dbObject, err := service.CreateUpdateAsset(c)
	doRespond(w, dbObject, err)
}

//CreateUpdateComment creates/updates a comment in the system
func CreateUpdateComment(w http.ResponseWriter, r *http.Request, c *service.AppContext) {
	dbObject, err := service.CreateUpdateComment(c)
	doRespond(w, dbObject, err)
}
