package model

import "gopkg.in/mgo.v2"

//Indicies defines all the indicies for Coral mongo database.
var Indicies = []Index{

	//Actions Indexes
	{
		Actions,
		mgo.Index{
			Key: []string{"source.id"},
		},
	},
	{
		Actions, mgo.Index{
			Key: []string{"user_id", "target_id", "target", "type"},
		},
	},

	//Assets Indexes
	{
		Assets, mgo.Index{
			Key: []string{"source.id"},
		},
	},
	{
		Assets, mgo.Index{
			Key: []string{"url"},
		},
	},

	//Users Indexes
	{
		Users, mgo.Index{
			Key: []string{"source.id"},
		},
	},

	//Comments Indexes
	{
		Comments, mgo.Index{
			Key: []string{"source.id"},
		},
	},
	{
		Comments, mgo.Index{
			Key: []string{"user_id"},
		},
	},
	{
		Comments, mgo.Index{
			Key: []string{"source.parent_id"},
		},
	},

	//Tags Indexes
	{
		Tags, mgo.Index{
			Key: []string{"name"},
		},
	},

	//TagTargets Indexes
	{
		TagTargets, mgo.Index{
			Key: []string{"target_id", "name", "target"},
		},
	},
}
