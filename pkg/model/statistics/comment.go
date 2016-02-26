package statistics

import (
	"log"

	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/model"
)

type CommentStatistics struct {
	Count int `json:"count" bson:"count"`

	// Replied concerns the comments of this group.
	RepliedCount      int      `json:"replied_count" bson:"replied_count"`
	RepliedToComments []string `json:"replied_comments" bson:"replied_comments"`
	RepliedToUsers    []string `json:"replied_users" bson:"replied_users"`
	RepliedRatio      float64  `json:"replied_ratio" bson:"replied_ratio"`

	// Reply concerns replies to the comments of this group.
	ReplyCount    int      `json:"reply_count" bson:"reply_count"`
	ReplyComments []string `json:"reply_comments" bson:"reply_comments"`
	ReplyUsers    []string `json:"reply_users" bson:"reply_users"`
	ReplyRatio    float64  `json:"reply_ratio" bson:"reply_ratio"`
}

type CommentStatisticsAccumulator struct {
	Counts, RepliedComments, RepliedUsers, ReplyComments, ReplyUsers aggregate.Int
}

func NewCommentStatisticsAccumulator() *CommentStatisticsAccumulator {
	return &CommentStatisticsAccumulator{
		Counts:          aggregate.NewInt(),
		RepliedComments: aggregate.NewInt(),
		RepliedUsers:    aggregate.NewInt(),
		ReplyComments:   aggregate.NewInt(),
		ReplyUsers:      aggregate.NewInt(),
	}
}

func (a *CommentStatisticsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	if comment, ok := object.(*model.Comment); ok {
		a.Counts.Add("count", 1)

		// Handle replied.
		if comment.ParentID.String() != "" {
			a.Counts.Add("replied_count", 1)
			a.RepliedComments.Add(comment.ParentID.String(), 1)
		}

		// Handle reply.
		for _, reply := range comment.Children {
			a.Counts.Add("reply_count", 1)
			a.ReplyComments.Add(reply.String(), 1)
		}
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
	}
}

func (a *CommentStatisticsAccumulator) CommentStatistics() *CommentStatistics {
	return &CommentStatistics{
		Count:             a.Counts.Total("count"),
		RepliedCount:      a.Counts.Total("replied_count"),
		RepliedToComments: a.RepliedComments.Keys(),
		RepliedToUsers:    a.RepliedUsers.Keys(),
		RepliedRatio:      a.Counts.Ratio("replied_count", "count"),
		ReplyCount:        a.Counts.Total("reply_count"),
		ReplyComments:     a.ReplyComments.Keys(),
		ReplyUsers:        a.ReplyUsers.Keys(),
		ReplyRatio:        a.Counts.Ratio("reply_count", "count"),
	}
}

type CommentTypes struct {
	All   *CommentStatistics
	Types map[string]*CommentStatistics `json:"types" bson:",inline"`
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

func (a *CommentTypesAccumulator) CommentTypes() *CommentTypes {
	types := make(map[string]*CommentStatistics)
	for key, value := range a.Types {
		types[key] = value.CommentStatistics()
	}

	return &CommentTypes{
		All:   a.All.CommentStatistics(),
		Types: types,
	}
}

type CommentDimensions struct {
	All   *CommentTypes
	Types map[string]map[string]*CommentTypes `json:"types" bson:",inline"`
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

		if assetID := comment.AssetID.String(); assetID != "" {

			if _, ok := a.Types["assets"]; !ok {
				a.Types["assets"] = make(map[string]*CommentTypesAccumulator)
			}

			if _, ok := a.Types["assets"][assetID]; !ok {
				a.Types["assets"][assetID] = NewCommentTypesAccumulator()
			}

			a.Types["assets"][assetID].Accumulate(ctx, object)
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

func (a *CommentDimensionsAccumulator) CommentDimensions() *CommentDimensions {
	types := make(map[string]map[string]*CommentTypes)
	for dimension, commentTypes := range a.Types {

		if _, ok := types[dimension]; !ok {
			types[dimension] = make(map[string]*CommentTypes)
		}

		for key, value := range commentTypes {
			types[dimension][key] = value.CommentTypes()
		}
	}

	return &CommentDimensions{
		All:   a.All.CommentTypes(),
		Types: types,
	}
}
