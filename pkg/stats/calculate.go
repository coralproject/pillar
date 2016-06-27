package stats

import (
	"github.com/ardanlabs/kit/log"
	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/backend"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/model/statistics"
)

// CalculateUserStatistics calculates User Statistics , creating a collection user_statistics with a document per user
func CalculateUserStatistics(ctx context.Context) error {
	// Look for a backend in the context and return an error if one is not
	// present.
	b, ok := ctx.Value("backend").(backend.Backend)
	if !ok {
		return backend.ErrBackendNotInitializedError
	}

	// Get the users iterator.
	iter, err := b.Find("users", nil) // To Do, only run it on the user that was changed through pillar
	if err != nil {
		return err
	}

	// Pipeline expects a generic input channel.
	in := make(chan interface{})

	go func() {
		defer close(in)
		if err := backend.Each(iter, func(doc interface{}) error {

			// Assert that the document is the type we expect.
			user, ok := doc.(*model.User)
			if !ok {
				return backend.ErrBackendType
			}

			in <- user
			return nil
		}); err != nil {
			log.Error(uid, "stats", err, "Calculating User Statistics.")
			return
		}
	}()

	accumulator := aggregate.Pipeline(ctx, in, func() aggregate.Accumulator {
		return statistics.NewUserAccumulator()
	})

	if userAccumulator, ok := accumulator.(*statistics.UserAccumulator); ok {

		for _, dimension := range userAccumulator.Dimensions() {
			if err := b.Upsert("dimensions", map[string]interface{}{"name": dimension.Name}, dimension); err != nil {
				log.Error(uid, "stats", err, "Calculating User Statistics.")
			}
		}
	}

	return nil
}

//
// import "github.com/coralproject/pillar/pkg/model"
//
// // Stats structs is the stats we are calculating
// type Stats struct {
// 	Comments          map[string]int `bson:"comments" json:"comments"`
// 	Sections          map[string]int `bson:"sections" json:"sections"`
// 	Authors           map[string]int `bson:"authors" json:"authors"`
// 	AcceptRatio       float32        `bson:"accept_ratio" json:"accept_ratio"`
// 	Replies           int            `bson:"replies" json:"replied"`
// 	RepliesPerComment float32        `bson:"replies_per_comment" json:"replies_per_comment"`
// 	//	Replied           int            `bson:"replied" json:"replied"`
// 	//	RepliedRatio      float32        `bson:"replied_ratio" json:"replied_ratio"`
// }
//
// func calc(cs []model.Comment) Stats {
//
// 	d := Stats{
// 		Comments: make(map[string]int),
// 		Sections: make(map[string]int),
// 		Authors:  make(map[string]int),
// 	}
//
// 	get(cs)
//
// 	for _, comment := range cs {
//
// 		//var user model.User
// 		//db.C("user").Find(bson.M{"_id": comment.UserID}).One(&user)
//
// 		// comments status
// 		// refactor with the Categorical Data Type / translation <-- This is very NYT specific - To Do
// 		d.Comments["total"]++
// 		switch comment.Status {
// 		case "1":
// 			d.Comments["unmoderated"]++
// 		case "2":
// 			d.Comments["accepted"]++
// 		case "3":
// 			d.Comments["rejected"]++
// 		case "4":
// 			d.Comments["escalated"]++
//
// 		}
//
// 		d.Replies += len(comment.Children)
//
// 		// var asset model.Asset
// 		// db.C("asset").Find(bson.M{"_id": comment.AssetID}).One(&asset)
// 		//
// 		// for _, author := range asset.Metadata.Authors {
// 		// 	if author.Title_case_name != "" {
// 		// 		d.Authors[author.Title_case_name]++
// 		// 	}
// 		// }
// 		//
// 		// if asset.Metadata.Section.DisplayName != "" {
// 		// 	d.Sections[asset.Metadata.Section.DisplayName]++
// 		// }
// 	}
//
// 	if d.Comments["total"] > 0 {
// 		d.AcceptRatio = float32(d.Comments["accepted"]) / float32(d.Comments["total"])
// 		d.RepliesPerComment = float32(d.Replies) / float32(d.Comments["total"])
// 	}
//
// 	return d
//
// }
