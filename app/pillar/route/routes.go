package route

import (
	"github.com/coralproject/pillar/app/pillar/handler"
	"net/http"
)

//Route defines mappings of end-points to handler methods
type Route struct {
	Method      string           `json:"method" bson:"method"`
	Pattern     string           `json:"pattern" bson:"pattern"`
	HandlerFunc http.HandlerFunc `json:"handler" bson:"handler"`
}

var routes = []Route{
	//Generic or Common ones
	Route{"GET", "/about", handler.AboutThisApp},

	//Import related end-points
	Route{"POST", "/api/import/asset", handler.ImportAsset},
	Route{"POST", "/api/import/user", handler.ImportUser},
	Route{"POST", "/api/import/comment", handler.ImportComment},
	Route{"POST", "/api/import/action", handler.ImportAction},
	Route{"POST", "/api/import/note", handler.ImportNote},
	Route{"POST", "/api/import/metadata", handler.ImportMetadata},
	Route{"POST", "/api/import/index", handler.CreateIndex},

	//Manage Tags
	Route{"GET", "/api/tags", handler.GetTags},
	Route{"POST", "/api/tag", handler.CreateUpdateTag},
	Route{"DELETE", "/api/tag", handler.DeleteTag},

	//Manage User Activities
	Route{"POST", "/api/cay/useraction", handler.HandleUserAction},

	//Update various entities
	Route{"POST", "/api/asset", handler.CreateUpdateAsset},
	Route{"POST", "/api/user", handler.CreateUpdateUser},
	Route{"POST", "/api/comment", handler.CreateUpdateComment},
	Route{"POST", "/api/index", handler.CreateIndex},
}
