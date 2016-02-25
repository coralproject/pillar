package statistics

import (
	"github.com/coralproject/pillar/pkg/model"
)

type UserActions struct{}

type UserComments struct{}

type UserStatistics struct {
	Actions  *UserActions  `json:"actions" bson:"actions"`
	Comments *UserComments `json:"comments" bson:"comments"`
}

type User struct {
	*model.User `bson:",inline"`
	Statistics  *UserStatistics `json:"stats" bson:"stats"`
}
