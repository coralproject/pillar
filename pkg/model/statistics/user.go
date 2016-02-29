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
	switch typedObject := object.(type) {
	default:
		log.Println("UserStatisticsAccumulator error: unexpected combine type")
	case *UserStatisticsAccumulator:
		a.Comments.Combine(typedObject.Comments)
	}
}

func (a *UserStatisticsAccumulator) UserStatistics(ctx context.Context) *UserStatistics {
	return &UserStatistics{
		Comments: a.Comments.CommentDimensions(ctx),
	}
}

type User struct {
	model.User `bson:",inline"`
	Statistics *UserStatistics `json:"statistics,omitempty" bson:"statistics,omitempty"`
	Reference  *UserStatistics `json:"reference,omitempty" bson:"reference,omitempty"`
}

type UserAccumulator struct {
	DimensionAccumulator map[string]aggregate.Int
}

func NewUserAccumulator() *UserAccumulator {
	return &UserAccumulator{
		DimensionAccumulator: make(map[string]aggregate.Int),
	}
}

func (a *UserAccumulator) Accumulate(ctx context.Context, object interface{}) {

	user, ok := object.(*model.User)
	if !ok {
		return
	}

	b, ok := ctx.Value("backend").(backend.Backend)
	if !ok {
		return
	}

	iter, err := b.Find("comments", map[string]interface{}{"user_id": user.ID})
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
			log.Println("Comment error:", err)
			return
		}
	}()

	accumulator := aggregate.Pipeline(ctx, comments, func() aggregate.Accumulator {
		return NewUserStatisticsAccumulator()
	})

	userStatisticsAccumulator, ok := accumulator.(*UserStatisticsAccumulator)
	if !ok {
		return
	}

	userStatistics := userStatisticsAccumulator.UserStatistics(ctx)
	if count := user.Stats["comments"]; count != nil && userStatistics.Comments.All.All.Count != count {
		log.Printf("Comment count didn't match, got %d, expected %d for %s", userStatistics.Comments.All.All.Count, count, user.ID.Hex())
	}

	if userStatistics.Comments.All.All.Count > 0 {
		if err := b.UpsertID("user_statistics", user.ID, &User{
			User:       *user,
			Statistics: userStatistics,
		}); err != nil {
			log.Println("User statistics error:", err)
		}

		userReference := userStatisticsAccumulator.UserStatistics(NewReferenceOnlyContext(ctx))
		if err := b.UpsertID("user_reference", user.ID, &User{
			User:      *user,
			Reference: userReference,
		}); err != nil {
			log.Println("User statistics error:", err)
		}
	}

	for dimension, commentTypes := range userStatistics.Comments.Types {
		if _, ok := a.DimensionAccumulator[dimension]; !ok {
			a.DimensionAccumulator[dimension] = aggregate.NewInt()
		}

		for key, _ := range commentTypes {
			a.DimensionAccumulator[dimension].Add(key, 1)
		}
	}
}

func (a *UserAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("UserAccumulator error: unexpected combine type")
	case *UserAccumulator:
		for key, value := range typedObject.DimensionAccumulator {
			if _, ok := a.DimensionAccumulator[key]; !ok {
				a.DimensionAccumulator[key] = aggregate.NewInt()
			}
			a.DimensionAccumulator[key].Combine(value)
		}
	}
}

func (a *UserAccumulator) Dimensions() []*model.Dimension {

	dimensions := make([]*model.Dimension, 0, len(a.DimensionAccumulator))

	for key, value := range a.DimensionAccumulator {
		dimensions = append(dimensions, &model.Dimension{
			Name:         key,
			Constituents: value.Keys(),
		})
	}

	return dimensions
}
