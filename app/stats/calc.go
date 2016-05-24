package main

import (

	//	"github.com/ardanlabs/kit/log"

	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Stats structs is the stats we are calculating
type Stats struct {
	Comments          map[string]int `bson:"comments" json:"comments"`
	Sections          map[string]int `bson:"sections" json:"sections"`
	Authors           map[string]int `bson:"authors" json:"authors"`
	AcceptRatio       float32        `bson:"accept_ratio" json:"accept_ratio"`
	Replies           int            `bson:"replies" json:"replied"`
	RepliesPerComment float32        `bson:"replies_per_comment" json:"replies_per_comment"`
	//	Replied           int            `bson:"replied" json:"replied"`
	//	RepliedRatio      float32        `bson:"replied_ratio" json:"replied_ratio"`
}

func calc(cs []model.Comment) Stats {

	d := model.Stats{
		Comments: make(map[string]int),
		Sections: make(map[string]int),
		Authors:  make(map[string]int),
	}

	get(cs)

	for _, comment := range cs {

		//var user model.User
		//db.C("user").Find(bson.M{"_id": comment.UserID}).One(&user)

		// comments status
		// refactor with the Categorical Data Type / translation
		d.Comments["total"]++
		switch comment.Status {
		case "1":
			d.Comments["unmoderated"]++
		case "2":
			d.Comments["accepted"]++
		case "3":
			d.Comments["rejected"]++
		case "4":
			d.Comments["escalated"]++

		}

		d.Replies += len(comment.Children)

		var asset model.Asset
		db.C("asset").Find(bson.M{"_id": comment.AssetID}).One(&asset)

		for _, author := range asset.Metadata.Authors {
			if author.Title_case_name != "" {
				d.Authors[author.Title_case_name]++
			}
		}

		if asset.Metadata.Section.DisplayName != "" {
			d.Sections[asset.Metadata.Section.DisplayName]++
		}

	}

	if d.Comments["total"] > 0 {
		d.AcceptRatio = float32(d.Comments["accepted"]) / float32(d.Comments["total"])
		d.RepliesPerComment = float32(d.Replies) / float32(d.Comments["total"])
	}

	return d

}
