package handler

import (
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
)

//ImportUser imports a new user to the system
func ImportUser(c *web.AppContext) {
	dbObject, err := service.ImportUser(c)
	doRespond(c, dbObject, err)
}

//ImportAsset imports a new asset to the system
func ImportAsset(c *web.AppContext) {
	dbObject, err := service.ImportAsset(c)
	doRespond(c, dbObject, err)
}

//ImportComment imports a new comment to the system
func ImportComment(c *web.AppContext) {
	dbObject, err := service.ImportComment(c)
	doRespond(c, dbObject, err)
}

//ImportAction imports actions into the system
func ImportAction(c *web.AppContext) {
	dbObject, err := service.ImportAction(c)
	doRespond(c, dbObject, err)
}

//ImportNote imports notes into the system
func ImportNote(c *web.AppContext) {
	dbObject, err := service.CreateNote(c)
	doRespond(c, dbObject, err)
}
