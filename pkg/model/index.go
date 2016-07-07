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

	//Form Submission Indexes
	{
		FormSubmissions, mgo.Index{
			Key:      []string{"$text:$**"},
			Unique:   false,
			DropDups: false,
		},
	},
	{
		FormSubmissions, mgo.Index{
			Key:      []string{"form_id"},
			Unique:   false,
			DropDups: false,
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

	//Users Indexes
	{
		Users, mgo.Index{
			Key: []string{"source.id"},
		},
	},

	// User Statistics indexes
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.all.count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.all.replied_count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.all.replied_ratio"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.all.reply_count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.all.reply_ratio"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.all.word_count_average"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.ModeratorDeleted.count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.ratios.ModeratorDeleted"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.ModeratorApproved.count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.ratios.ModeratorApproved"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.CommunityFlagged.count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.comments.all.ratios.Communityflagged"},
		},
	},

	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.actions.received.likes.count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.actions.performed.likes.count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.actions.received.flags.count"},
		},
	},
	{
		UserStatistics, mgo.Index{
			Key: []string{"statistics.actions.performed.flags.count"},
		},
	},
}
