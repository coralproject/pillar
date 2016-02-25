package statistics

import (
	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
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

func NewActionStatisticsAccumulator(ctx context.Context) *ActionStatisticsAccumulator {
	return &ActionStatisticsAccumulator{
		Counts:   aggregate.NewInt(),
		Comments: aggregate.NewInt(),
		Assets:   aggregate.NewInt(),
		Users:    aggregate.NewInt(),
	}
}

type ActionDimensions struct {
	All   *ActionStatistics
	Types map[string]*ActionStatistics `json:"types" bson:",inline"`
}
