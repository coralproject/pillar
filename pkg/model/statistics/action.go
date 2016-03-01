package statistics

import (
	"log"

	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/model"
)

type ActionStatistics struct {
	Count    int            `json:"count" bson:"count,minsize"`
	Users    map[string]int `json:"users,omitempty" bson:"users,omitempty,minsize"`
	Comments []string       `json:"comments,omitempty" bson:"comments,omitempty"`
	Assets   []string       `json:"assets,omitempty" bson:"assets,omitempty"`
	Sections []string       `json:"sections,omitempty" bson:"sections,omitempty"`
}

type PerformedActionStatisticsAccumulator struct {
	Counts, Comments, Assets, Users aggregate.Int
}

func NewPerformedActionStatisticsAccumulator() *PerformedActionStatisticsAccumulator {
	return &PerformedActionStatisticsAccumulator{
		Counts:   aggregate.NewInt(),
		Comments: aggregate.NewInt(),
		Assets:   aggregate.NewInt(),
		Users:    aggregate.NewInt(),
	}
}

func (a *PerformedActionStatisticsAccumulator) Accumulate(ctx context.Context, object interface{}) {
	if action, ok := object.(*model.Action); ok {
		a.Counts.Add("count", 1)
		switch {
		case action.Target == "comment":
			a.Comments.Add(action.TargetID.Hex(), 1)
		case action.Target == "user":
			a.Users.Add(action.TargetID.Hex(), 1)
		case action.Target == "asset":
			a.Assets.Add(action.TargetID.Hex(), 1)
		}
	}
}

func (a *PerformedActionStatisticsAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("PerformedActionStatisticsAccumulator error: unexpected combine type")
	case *PerformedActionStatisticsAccumulator:
		a.Counts.Combine(typedObject.Counts)
		a.Comments.Combine(typedObject.Comments)
		a.Assets.Combine(typedObject.Assets)
		a.Users.Combine(typedObject.Users)
	}
}

func (a *PerformedActionStatisticsAccumulator) ActionStatistics(ctx context.Context) *ActionStatistics {
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

type ActionTypes struct {
	All    *ActionStatistics            `json:"all,omitempty" bson:"all,omitempty"`
	Types  map[string]*ActionStatistics `json:"types" bson:",inline"`
	Ratios map[string]float64           `json:"ratios,omitempty" bson:"ratios,omitempty"`
}

type ActionTypesAccumulator struct {
	All   *PerformedActionStatisticsAccumulator
	Types map[string]*PerformedActionStatisticsAccumulator
}

func NewActionTypesAccumulator() *ActionTypesAccumulator {
	return &ActionTypesAccumulator{
		All:   NewPerformedActionStatisticsAccumulator(),
		Types: make(map[string]*PerformedActionStatisticsAccumulator),
	}
}

func (a *ActionTypesAccumulator) Accumulate(ctx context.Context, object interface{}) {
	if action, ok := object.(*model.Action); ok {
		a.All.Accumulate(ctx, object)
		if action.Type != "" {
			if _, ok := a.Types[action.Type]; !ok {
				a.Types[action.Type] = NewPerformedActionStatisticsAccumulator()
			}
			a.Types[action.Type].Accumulate(ctx, object)
		}
	}
}

func (a *ActionTypesAccumulator) Combine(object interface{}) {
	switch typedObject := object.(type) {
	default:
		log.Println("ActionTypesAccumulator error: unexpected combine type")
	case *ActionTypesAccumulator:
		a.All.Combine(typedObject.All)
		for key, value := range typedObject.Types {
			if _, ok := a.Types[key]; !ok {
				a.Types[key] = NewPerformedActionStatisticsAccumulator()
			}
			a.Types[key].Combine(value)
		}
	}
}

func (a *ActionTypesAccumulator) ActionStatistics(ctx context.Context) *ActionTypes {
	all := a.All.ActionStatistics(ctx)
	types := make(map[string]*ActionStatistics)
	ratios := make(map[string]float64)
	for key, value := range a.Types {
		types[key] = value.ActionStatistics(ctx)
		if all.Count > 0 {
			ratios[key] = float64(types[key].Count) / float64(all.Count)
		}
	}

	return &ActionTypes{
		All:    all,
		Types:  types,
		Ratios: ratios,
	}
}
