package statistics

import (
	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/model"
)

type UserActions struct {
	Performed *ActionDimensions `json:"performed" bson:"performed"`
	Received  *ActionDimensions `json:"received" bson:"received"`
}

type UserComments CommentDimensions

type UserStatistics struct {
	Actions  *UserActions  `json:"actions" bson:"actions"`
	Comments *UserComments `json:"comments" bson:"comments"`
}

type UserStatisticsAccumulator struct {
	Comments *CommentDimensionsAccumulator
}

func NewUserStatisticsAccumulator() *UserStatisticsAccumulator {
	return &UserStatisticsAccumulator{
		Comments: NewCommentDimensionsAccumulator(),
	}
}

func (a *UserStatisticsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	a.Comments.Accumulate(ctx, object)
}

func (a *UserStatisticsAccumulator) Combine(object interface{}) {
	a.Comments.Combine(object)
}

type User struct {
	*model.User `bson:",inline"`
	Statistics  *UserStatistics `json:"stats" bson:"stats"`
}
