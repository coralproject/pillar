package statistics

import (
	"golang.org/x/net/context"
	"log"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/backend"
	"github.com/coralproject/pillar/pkg/backend/iterator"
	"github.com/coralproject/pillar/pkg/model"
)

type UserActions struct {
	Performed *ActionDimensions `json:"performed" bson:"performed"`
	Received  *ActionDimensions `json:"received" bson:"received"`
}

type UserStatistics struct {
	Actions  *UserActions       `json:"actions" bson:"actions"`
	Comments *CommentDimensions `json:"comments" bson:"comments"`
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

func (a *UserStatisticsAccumulator) UserStatistics() *UserStatistics {
	return &UserStatistics{
		Comments: a.Comments.CommentDimensions(),
	}
}

type User struct {
	*model.User `bson:",inline"`
	Statistics  *UserStatistics `json:"stats" bson:"stats"`
}

type UserAccumulator struct {
}

func NewUserAccumulator() *UserAccumulator {
	return &UserAccumulator{}
}

func (a *UserAccumulator) Accumulate(ctx context.Context, object interface{}) {

	user_id, ok := object.(string)
	if !ok {
		return
	}

	b, ok := ctx.Value("backend").(backend.Backend)
	if !ok {
		return
	}

	iter, err := b.Find("comments", map[string]interface{}{"user_id": user_id})
	if err != nil {
		return
	}

	comments := make(chan interface{})

	go func() {
		defer close(comments)
		if err := iterator.Each(iter, func(doc interface{}) error {

			// Assert that the document is the type we expect.
			comment, ok := doc.(*model.Comment)
			if !ok {
				return backend.BackendTypeError
			}

			comments <- comment
			return nil
		}); err != nil {
			return
		}
	}()

	accumulator := aggregate.Pipeline(ctx, comments,
		func() aggregate.Accumulator { return NewUserStatisticsAccumulator() },
	)

	UserStatisticsAccumulator, ok := accumulator.(*UserStatisticsAccumulator)
	if !ok {
		return
	}

	log.Println(UserStatisticsAccumulator.UserStatistics())
}

func (a *UserAccumulator) Combine(object interface{}) {
	// Noop.
}
