package handler

import (
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
)

//ImportUser imports a new user to the system
func ImportUser(c *web.AppContext) {
	c.Event = model.EventUserImport
	dbObject, err := service.ImportUser(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Import_User", 1, 1.0)
}

//ImportAsset imports a new asset to the system
func ImportAsset(c *web.AppContext) {
	c.Event = model.EventAssetImport
	dbObject, err := service.ImportAsset(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Import_Asset", 1, 1.0)
}

//ImportComment imports a new comment to the system
func ImportComment(c *web.AppContext) {
	c.Event = model.EventCommentImport
	dbObject, err := service.ImportComment(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Import_Comment", 1, 1.0)
}

//ImportAction imports actions into the system
func ImportAction(c *web.AppContext) {
	c.Event = model.EventActionImport
	dbObject, err := service.ImportAction(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Import_Action", 1, 1.0)
}

//ImportNote imports notes into the system
func ImportNote(c *web.AppContext) {
	c.Event = model.EventNoteImport
	dbObject, err := service.CreateNote(c)
	doRespond(c, dbObject, err)
	c.SD.Client.Inc("Import_Note", 1, 1.0)
}
