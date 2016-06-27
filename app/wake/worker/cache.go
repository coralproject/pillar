package worker

import (
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
)

var assetCache map[string]model.Asset
var userCache  map[string]model.User

func init() {
	assetCache = make(map[string]model.Asset)
	userCache = make(map[string]model.User)
}

func getUser(mdb *db.MongoDB, id bson.ObjectId) model.User {

	o := userCache[id.Hex()]
	if o.ID != "" {
		//log.Printf("Found User from Cache!\n")
		return o
	}

	var object model.User
	mdb.DB.C(CollectionUsers).FindId(id).One(&object)
	userCache[id.Hex()] = object
	return object
}

func getAsset(mdb *db.MongoDB, id bson.ObjectId) model.Asset {

	o := assetCache[id.Hex()]
	if o.ID != "" {
		//log.Printf("Found Asset from Cache!\n")
		return o
	}

	var object model.Asset
	mdb.DB.C(CollectionAssets).FindId(id).One(&object)
	assetCache[id.Hex()] = object
	return object
}
