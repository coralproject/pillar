package stats

import (
	//  "errors"

	"github.com/coralproject/pillar/server/model"
	//	"github.com/coralproject/pillar/server/service"

	//"github.com/ardanlabs/kit/log"
	//"gopkg.in/mgo.v2/bson"
)

// onCreateComment takes a comment model and updates appropriate stats
func onCreateComment(c model.Comment) model.Comment {

	// concept: each thing we can measure should have a count and a list
	//    if we query to get the list, we can find the count using go to save queries
	//    doing the counts with go will also ensure that the count matches the array

	// TODO: user.stats.comments - comment count
	// TODO: user.stats.comment_list - array of comment IDs

	// TODO: user.stats.first_comment_date: if this is the first comment for the user set this
	// TODO: user.stats.last_comment_date: if this is the latest, set latest date

	// TODO: user.stats.replies - if this is a reply, count replies and insert into
	// TODO: user.stats.reply_list - if this is a reply, update array of comments user has replied to

	// TODO: user.stats.assets_commented_on_list - array of assets commented on
	// TODO: user.stats.assets_commented_on - count how many distinct assets this user has commented on

	// TODO later: user.stats.median_words - median words of all user comments
	// TODO later: user.stats.mean_words - mean word count of all user comments

	// TODO: user.stats.users_replied_to - count how many distinct users this user has applied to
	// TODO: user.stats.users_replied_to_list - array of users replied to

	// TODO: asset.stats.comments - comment count
	// TODO: asset.stats.comment_list - array of comment IDs

	// TODO: asset.stats.first_comment_date: if this is the first comment for the user set this
	// TODO: asset.stats.last_comment_date: if this is the latest, set latest date

	// TODO: asset.stats.replies - if this is a reply, count replies and insert into
	// TODO: asset.stats.reply_list - if this is a reply, update array of comments user has replied to

	// TODO: asset.stats.assets_commented_on_list - array of assets commented on
	// TODO: asset.stats.assets_commented_on - count how many distinct assets this user has commented on

	// TODO later: asset.stats.median_words - median words of all user comments
	// TODO later: asset.stats.mean_words - mean word count of all user comments

	//	db.DB("coral").C("comment").FindId(bson.ObjectId("567b0850e19ac8852dd2bb5c")).One(&c)

	return c

}

/*
func updateCommentOnAssetCount(_id string) error {

	return nil

}
*/
