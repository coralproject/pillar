package statistics

import (
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
	switch typedObject := object.(type) {
	default:
		// May want to log here to indicate an unhandleable object.
	case *model.Comment:
		a.Counts.Add("count", 1)

		// Handle replied.
		if typedObject.ParentID.String() != "" {
			a.Counts.Add("replied_count", 1)
			a.RepliedComments.Add(typedObject.ParentID.String(), 1)
		}

		// Handle reply.
		for _, reply := range typedObject.Children {
			a.Counts.Add("reply_count", 1)
			a.ReplyComments.Add(reply.String(), 1)
		}
	}
}

func (a *CommentStatisticsAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		// May want to log here to indicate an unhandleable object.
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
		ReplyCount:        a.Counts.Total("reply_count"),
		ReplyComments:     a.ReplyComments.Keys(),
		ReplyUsers:        a.ReplyUsers.Keys(),
	}
}

type CommentDimensions struct {
	All   *CommentStatistics
	Types map[string]*CommentStatistics `json:"types" bson:",inline"`
}

type CommentDimensionsAccumulator struct {
	All   *CommentStatisticsAccumulator
	Types map[string]*CommentStatisticsAccumulator
}

func NewCommentDimensionsAccumulator() *CommentDimensionsAccumulator {
	return &CommentDimensionsAccumulator{
		All:   NewCommentStatisticsAccumulator(),
		Types: make(map[string]*CommentStatisticsAccumulator),
	}
}

func (a *CommentDimensionsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	a.All.Accumulate(ctx, object)

	// Handle types
}

func (a *CommentDimensionsAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		// May want to log here to indicate an unhandleable object.
	case *CommentDimensionsAccumulator:
		a.All.Combine(object)
		for key, value := range typedObject.Types {
			a.Types[key].Combine(value)
		}
	}
}

func (a *CommentDimensionsAccumulator) CommentDimensions() *CommentDimensions {
	types := make(map[string]*CommentStatistics)
	for key, value := range a.Types {
		types[key] = value.CommentStatistics()
	}

	return &CommentDimensions{
		All:   a.All.CommentStatistics(),
		Types: types,
	}
}
