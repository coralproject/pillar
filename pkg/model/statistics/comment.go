package statistics

import (
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/backend"
	"github.com/coralproject/pillar/pkg/model"
)

type CommentStatistics struct {
	Count int `json:"count" bson:"count,minsize"`

	// Replied concerns the comments of this group.
	RepliedCount      int      `json:"replied_count" bson:"replied_count,minsize"`
	RepliedToComments []string `json:"replied_comments,omitempty" bson:"replied_comments,omitempty"`
	RepliedToUsers    []string `json:"replied_users,omitempty" bson:"replied_users,omitempty"`
	RepliedRatio      float64  `json:"replied_ratio" bson:"replied_ratio"`

	// Reply concerns replies to the comments of this group.
	ReplyCount    int      `json:"reply_count" bson:"reply_count,minsize"`
	ReplyComments []string `json:"reply_comments,omitempty" bson:"reply_comments,omitempty"`
	ReplyUsers    []string `json:"reply_users,omitempty" bson:"reply_users,omitempty"`
	ReplyRatio    float64  `json:"reply_ratio" bson:"reply_ratio"`

	// First and last comments.
	First time.Time `json:"first,omitempty" bson:"first,omitempty"`
	Last  time.Time `json:"last,omitempty" bson:"last,omitempty"`

	// Text analysis.
	WordCountAverage float64 `json:"word_count_average" bson:"word_count_average"`
}

type CommentStatisticsAccumulator struct {
	Counts, RepliedComments, RepliedUsers, ReplyComments, ReplyUsers aggregate.Int
	First                                                            *aggregate.Oldest
	Last                                                             *aggregate.Newest
}

func NewCommentStatisticsAccumulator() *CommentStatisticsAccumulator {
	return &CommentStatisticsAccumulator{
		Counts:          aggregate.NewInt(),
		RepliedComments: aggregate.NewInt(),
		RepliedUsers:    aggregate.NewInt(),
		ReplyComments:   aggregate.NewInt(),
		ReplyUsers:      aggregate.NewInt(),
		First:           aggregate.NewOldest(),
		Last:            aggregate.NewNewest(),
	}
}

func (a *CommentStatisticsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	if comment, ok := object.(*model.Comment); ok {
		a.Counts.Add("count", 1)

		// Word count.
		a.Counts.Add("word_count", len(strings.Split(comment.Body, " ")))

		// Handle replied.
		if comment.RootID.Hex() != "" {
			a.Counts.Add("replied_count", 1)
			a.RepliedComments.Add(comment.RootID.Hex(), 1)
		}

		// Handle reply.
		for _, reply := range comment.Children {
			a.Counts.Add("reply_count", 1)
			a.ReplyComments.Add(reply.Hex(), 1)
		}

		a.First.Check(comment.DateCreated)
		a.Last.Check(comment.DateCreated)
	}
}

func (a *CommentStatisticsAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("CommentStatisticsAccumulator error: unexpected combine type")
	case *CommentStatisticsAccumulator:
		a.Counts.Combine(typedObject.Counts)
		a.RepliedComments.Combine(typedObject.RepliedComments)
		a.RepliedUsers.Combine(typedObject.RepliedUsers)
		a.ReplyComments.Combine(typedObject.ReplyComments)
		a.ReplyUsers.Combine(typedObject.ReplyUsers)
		a.First.Combine(typedObject.First)
		a.Last.Combine(typedObject.Last)
	}
}

func (a *CommentStatisticsAccumulator) CommentStatistics(ctx context.Context) *CommentStatistics {
	commentStatistics := &CommentStatistics{
		Count:            a.Counts.Total("count"),
		RepliedCount:     a.Counts.Total("replied_count"),
		RepliedRatio:     a.Counts.Ratio("replied_count", "count"),
		ReplyCount:       a.Counts.Total("reply_count"),
		ReplyRatio:       a.Counts.Ratio("reply_count", "count"),
		WordCountAverage: float64(a.Counts.Total("word_count")) / float64(a.Counts.Total("count")),
	}

	if a.First.Valid {
		commentStatistics.First = a.First.Time
	}

	if a.Last.Valid {
		commentStatistics.Last = a.Last.Time
	}

	//  add unbound values if this isn't a reference-only request.
	if !OmitReferencesFromContext(ctx) {
		commentStatistics.RepliedToComments = a.RepliedComments.Keys()
		commentStatistics.RepliedToUsers = a.RepliedUsers.Keys()
		commentStatistics.ReplyComments = a.ReplyComments.Keys()
		commentStatistics.ReplyUsers = a.ReplyUsers.Keys()
	}

	return commentStatistics
}

type CommentTypes struct {
	All    *CommentStatistics            `json:"all,omitempty" bson:"all,omitempty"`
	Types  map[string]*CommentStatistics `json:"types,omitempty" bson:",inline"`
	Ratios map[string]float64            `json:"ratios,omitempty" bson:"ratios,omitempty"`
}

type CommentTypesAccumulator struct {
	All   *CommentStatisticsAccumulator
	Types map[string]*CommentStatisticsAccumulator
}

func NewCommentTypesAccumulator() *CommentTypesAccumulator {
	return &CommentTypesAccumulator{
		All:   NewCommentStatisticsAccumulator(),
		Types: make(map[string]*CommentStatisticsAccumulator),
	}
}

func (a *CommentTypesAccumulator) Accumulate(ctx context.Context, object interface{}) {
	a.All.Accumulate(ctx, object)

	if comment, ok := object.(*model.Comment); ok {
		if comment.Status != "" {
			if _, ok := a.Types[comment.Status]; !ok {
				a.Types[comment.Status] = NewCommentStatisticsAccumulator()
			}
			a.Types[comment.Status].Accumulate(ctx, comment)
		}
	}
}

func (a *CommentTypesAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("CommentTypesAccumulator error: unexpected combine type")
	case *CommentTypesAccumulator:
		a.All.Combine(typedObject.All)
		for key, value := range typedObject.Types {
			if _, ok := a.Types[key]; !ok {
				a.Types[key] = NewCommentStatisticsAccumulator()
			}
			a.Types[key].Combine(value)
		}
	}
}

func (a *CommentTypesAccumulator) CommentTypes(ctx context.Context) *CommentTypes {
	all := a.All.CommentStatistics(ctx)

	types := make(map[string]*CommentStatistics)
	ratios := make(map[string]float64)
	for key, value := range a.Types {
		types[key] = value.CommentStatistics(ctx)
		if all.Count > 0 {
			ratios[key] = float64(types[key].Count) / float64(all.Count)
		}
	}

	return &CommentTypes{
		All:    all,
		Types:  types,
		Ratios: ratios,
	}
}

type CommentDimensions struct {
	All   *CommentTypes                       `json:"all,omitempty" bson:"all,omitempty"`
	Types map[string]map[string]*CommentTypes `json:"types,omitempty" bson:",inline"`
}

type CommentDimensionsAccumulator struct {
	All   *CommentTypesAccumulator
	Types map[string]map[string]*CommentTypesAccumulator
}

func NewCommentDimensionsAccumulator() *CommentDimensionsAccumulator {
	return &CommentDimensionsAccumulator{
		All:   NewCommentTypesAccumulator(),
		Types: make(map[string]map[string]*CommentTypesAccumulator),
	}
}

func (a *CommentDimensionsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	a.All.Accumulate(ctx, object)

	if comment, ok := object.(*model.Comment); ok {
		if assetID := comment.AssetID; assetID != "" {

			b, ok := ctx.Value("backend").(backend.Backend)
			if !ok {
				log.Println("CommentDimensionsAccumulator accumulate error: backend not found")
				return
			}

			assetObject, err := b.FindID("assets", assetID)
			if err != nil {
				log.Println("CommentDimensionsAccumulator accumulate error:", err)
				return
			}

			asset, ok := assetObject.(*model.Asset)
			if !ok {
				log.Println("CommentDimensionsAccumulator accumulate error: find returned wrong type")
				return
			}

			// Handle the asset by ID.
			if _, ok := a.Types["assets"]; !ok {
				a.Types["assets"] = make(map[string]*CommentTypesAccumulator)
			}

			assetIDString := assetID.Hex()
			if _, ok := a.Types["assets"][assetIDString]; !ok {
				a.Types["assets"][assetIDString] = NewCommentTypesAccumulator()
			}

			a.Types["assets"][assetIDString].Accumulate(ctx, object)

			// Handle authors.
			if _, ok := a.Types["author"]; !ok {
				a.Types["author"] = make(map[string]*CommentTypesAccumulator)
			}

			for _, author := range asset.Authors {
				if author.ID != "" {
					if _, ok := a.Types["assets"][author.ID]; !ok {
						a.Types["author"][author.ID] = NewCommentTypesAccumulator()
					}

					a.Types["author"][author.ID].Accumulate(ctx, object)
				}
			}

			// Handle the section.
			if _, ok := a.Types["section"]; !ok {
				a.Types["section"] = make(map[string]*CommentTypesAccumulator)
			}

			if asset.Section != "" {
				if _, ok := a.Types["section"][asset.Section]; !ok {
					a.Types["section"][asset.Section] = NewCommentTypesAccumulator()
				}

				a.Types["section"][asset.Section].Accumulate(ctx, object)
			}
		}
	}
}

func (a *CommentDimensionsAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("CommentDimensionsAccumulator error: unexpected combine type")
	case *CommentDimensionsAccumulator:
		a.All.Combine(typedObject.All)
		for dimension, commentTypes := range typedObject.Types {
			if _, ok := a.Types[dimension]; !ok {
				a.Types[dimension] = make(map[string]*CommentTypesAccumulator)
			}
			for key, value := range commentTypes {
				if _, ok := a.Types[dimension][key]; !ok {
					a.Types[dimension][key] = NewCommentTypesAccumulator()
				}
				a.Types[dimension][key].Combine(value)
			}
		}
	}
}

func (a *CommentDimensionsAccumulator) CommentDimensions(ctx context.Context) *CommentDimensions {
	types := make(map[string]map[string]*CommentTypes)
	for dimension, commentTypes := range a.Types {

		if _, ok := types[dimension]; !ok {
			types[dimension] = make(map[string]*CommentTypes)
		}

		for key, value := range commentTypes {
			types[dimension][key] = value.CommentTypes(ctx)
		}
	}

	return &CommentDimensions{
		All:   a.All.CommentTypes(ctx),
		Types: types,
	}
}
