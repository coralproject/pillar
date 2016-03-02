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
	Performed *ActionTypes `json:"performed" bson:"performed"`
	Received  *ActionTypes `json:"received" bson:"received"`
}

type UserActionsAccumulator struct {
	Performed *ActionTypesAccumulator
	Received  *ReceivedActionTypesAccumulator
}

func NewUserActionsAccumulator() *UserActionsAccumulator {
	return &UserActionsAccumulator{
		Performed: NewActionTypesAccumulator(),
		Received:  NewReceivedActionTypesAccumulator(),
	}
}

func (a *UserActionsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	if action, ok := object.(*model.Action); ok {
		userID := ctx.Value("user_id")
		if action.UserID == userID {
			a.Performed.Accumulate(ctx, object)
		} else {
			a.Received.Accumulate(ctx, object)
		}
	}
}

func (a *UserActionsAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("UserActionsAccumulator error: unexpected combine type")
	case *UserActionsAccumulator:
		a.Performed.Combine(typedObject.Performed)
		a.Received.Combine(typedObject.Received)
	}
}

func (a *UserActionsAccumulator) UserActions(ctx context.Context) *UserActions {
	return &UserActions{
		Performed: a.Performed.ActionStatistics(ctx),
		Received:  a.Received.ActionStatistics(ctx),
	}
}

type UserStatistics struct {
	Actions  *UserActions       `json:"actions" bson:"actions"`
	Comments *CommentDimensions `json:"comments" bson:"comments"`
}

type UserStatisticsAccumulator struct {
	Comments     *CommentDimensionsAccumulator
	Actions      *UserActionsAccumulator
	CommentIDMap map[interface{}]struct{}
}

func NewUserStatisticsAccumulator() *UserStatisticsAccumulator {
	return &UserStatisticsAccumulator{
		Comments:     NewCommentDimensionsAccumulator(),
		Actions:      NewUserActionsAccumulator(),
		CommentIDMap: make(map[interface{}]struct{}),
	}
}

func (a *UserStatisticsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	a.Comments.Accumulate(ctx, object)
	a.Actions.Accumulate(ctx, object)
	if comment, ok := object.(*model.Comment); ok {
		a.CommentIDMap[comment.ID] = struct{}{}
	}
}

func (a *UserStatisticsAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("UserStatisticsAccumulator error: unexpected combine type")
	case *UserStatisticsAccumulator:
		a.Comments.Combine(typedObject.Comments)
		a.Actions.Combine(typedObject.Actions)
		for key, _ := range typedObject.CommentIDMap {
			a.CommentIDMap[key] = struct{}{}
		}
	}
}

func (a *UserStatisticsAccumulator) UserStatistics(ctx context.Context) *UserStatistics {
	return &UserStatistics{
		Comments: a.Comments.CommentDimensions(ctx),
		Actions:  a.Actions.UserActions(ctx),
	}
}

func (a *UserStatisticsAccumulator) CommentIDs(ctx context.Context) []interface{} {
	commentIDs := make([]interface{}, 0, len(a.CommentIDMap))
	for key, _ := range a.CommentIDMap {
		commentIDs = append(commentIDs, key)
	}
	return commentIDs
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

	// Add the user ID to the context.
	ctx = context.WithValue(ctx, "user_id", user.ID)

	b, ok := ctx.Value("backend").(backend.Backend)
	if !ok {
		return
	}

	commentsIterator, err := b.Find("comments", map[string]interface{}{
		"user_id": user.ID,
	})
	if err != nil {
		return
	}

	commentsAccumulator :=
		aggregate.Pipeline(ctx, iterator.EachChannel(commentsIterator), func() aggregate.Accumulator {
			return NewUserStatisticsAccumulator()
		})

	userStatisticsAccumulator, ok := commentsAccumulator.(*UserStatisticsAccumulator)
	if !ok {
		return
	}

	actionsRevceivedIterator, err := b.Find("actions", map[string]interface{}{
		"target":    "comments",
		"target_id": userStatisticsAccumulator.CommentIDs(ctx),
	})
	if err != nil {
		return
	}

	actionsReceivedAccumulator :=
		aggregate.Pipeline(ctx, iterator.EachChannel(actionsRevceivedIterator), func() aggregate.Accumulator {
			return NewUserStatisticsAccumulator()
		})

	userStatisticsAccumulator.Combine(actionsReceivedAccumulator)

	actionsPerformedIterator, err := b.Find("actions", map[string]interface{}{
		"user_id": user.ID,
	})
	if err != nil {
		return
	}

	actionsPerformedAccumulator :=
		aggregate.Pipeline(ctx, iterator.EachChannel(actionsPerformedIterator), func() aggregate.Accumulator {
			return NewUserStatisticsAccumulator()
		})

	userStatisticsAccumulator.Combine(actionsPerformedAccumulator)

	// 		actionsReceivedIterator, err := b.Find("actions", map[string]interface{}{"":"","target_id": user.ID})
	// if err != nil {
	// 	return
	// }

	userStatisticsReference := userStatisticsAccumulator.UserStatistics(ctx)
	if count := user.Stats["comments"]; count != nil && userStatisticsReference.Comments.All.All.Count != count {
		log.Printf("Comment count didn't match, got %d, expected %d for %s", userStatisticsReference.Comments.All.All.Count, count, user.ID.Hex())
	}

	if userStatisticsReference.Comments.All.All.Count > 0 {
		if err := b.UpsertID("user_reference", user.ID, &User{
			User:      *user,
			Reference: userStatisticsReference,
		}); err != nil {
			log.Println("User statistics error:", err)
		}

		userStatistics := userStatisticsAccumulator.UserStatistics(NewOmitReferencesContext(ctx))
		if err := b.UpsertID("user_statistics", user.ID, &User{
			User:       *user,
			Statistics: userStatistics,
		}); err != nil {
			log.Println("User statistics error:", err)
		}
	}

	for dimension, commentTypes := range userStatisticsReference.Comments.Types {
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
