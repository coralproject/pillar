package statistics

import (
	"log"

	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/model"
)

type ActionStatistics struct {
	Count    int            `json:"count" bson:"count"`
	Users    map[string]int `json:"users" bson:"users"`
	Comments []string       `json:"comments" bson:"comments"`
	Assets   []string       `json:"assets" bson:"assets"`
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
		a.Users.Add(action.UserID.String(), 1)
		a.Comments.Add(action.TargetID.String(), 1)
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

func (a *ActionStatisticsAccumulator) ActionStatistics() *ActionStatistics {
	return &ActionStatistics{
		Count:    a.Counts.Total("count"),
		Users:    a.Users,
		Comments: a.Comments.Keys(),
		Assets:   a.Assets.Keys(),
	}
}

type ActionDimensions struct {
	All   *ActionStatistics
	Types map[string]*ActionStatistics `json:"types" bson:",inline"`
}
