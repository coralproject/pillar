package route

import (
	"github.com/coralproject/pillar/app/pillar/handler"
)

//Route defines mappings of end-points to handler methods
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc handler.AppHandlerFunc
}

var routes = []Route{
	//Generic or Common ones
	Route{"GET", "/about", handler.AboutThisApp},

	//Import Handlers
	Route{"POST", "/api/import/asset", handler.ImportAsset},
	Route{"POST", "/api/import/user", handler.ImportUser},
	Route{"POST", "/api/import/comment", handler.ImportComment},
	Route{"POST", "/api/import/action", handler.ImportAction},
	Route{"POST", "/api/import/note", handler.ImportNote},

	//Tag Handlers
	Route{"GET", "/api/tags", handler.GetTags},
	Route{"POST", "/api/tag", handler.CreateUpdateTag},
	Route{"DELETE", "/api/tag", handler.DeleteTag},

	//Manage User Activities
	Route{"POST", "/api/cay/useraction", handler.HandleUserAction},

	//Update Handlers
	Route{"POST", "/api/asset", handler.CreateUpdateAsset},
	Route{"POST", "/api/user", handler.CreateUpdateUser},
	Route{"POST", "/api/comment", handler.CreateUpdateComment},
	Route{"POST", "/api/index", handler.CreateIndex},
	Route{"POST", "/api/metadata", handler.UpdateMetadata},
}
