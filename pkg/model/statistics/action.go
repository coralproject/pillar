package statistics

import (
	"log"

	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/model"
)

type ActionStatistics struct {
	Count    int            `json:"count" bson:"count"`
	Users    map[string]int `json:"users,omitempty" bson:"users,omitempty"`
	Comments []string       `json:"comments,omitempty" bson:"comments,omitempty"`
	Assets   []string       `json:"assets,omitempty" bson:"assets,omitempty"`
	Sections []string       `json:"sections,omitempty" bson:"sections,omitempty"`
}

type ActionStatisticsAccumulator struct {
	Counts, Comments, Assets, Users aggregate.Int
}

func NewActionStatisticsAccumulator() *ActionStatisticsAccumulator {
	return &ActionStatisticsAccumulator{
		Counts:   aggregate.NewInt(),
		Comments: aggregate.NewInt(),
		Assets:   aggregate.NewInt(),
		Users:    aggregate.NewInt(),
	}
}

func (a *ActionStatisticsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	if action, ok := object.(*model.Action); ok {
		a.Counts.Add("count", 1)
		a.Users.Add(action.UserID.Hex(), 1)
		if action.Target == "comment" {
			a.Comments.Add(action.TargetID.Hex(), 1)
		}
	}
}

func (a *ActionStatisticsAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("ActionStatisticsAccumulator error: unexpected combine type")
	case *ActionStatisticsAccumulator:
		a.Counts.Combine(typedObject.Counts)
		a.Comments.Combine(typedObject.Comments)
		a.Assets.Combine(typedObject.Assets)
		a.Users.Combine(typedObject.Users)
	}
}

func (a *ActionStatisticsAccumulator) ActionStatistics(ctx context.Context) *ActionStatistics {
	actionStatistics := &ActionStatistics{
		Count: a.Counts.Total("count"),
	}

	if !OmitReferencesFromContext(ctx) {
		actionStatistics.Users = a.Users
		actionStatistics.Comments = a.Comments.Keys()
		actionStatistics.Assets = a.Assets.Keys()
	}

	return actionStatistics
}

type ActionDimensions struct {
	All   *ActionStatistics            `json:"all,omitempty" bson:"all,omitempty"`
	Types map[string]*ActionStatistics `json:"types" bson:",inline"`
}
