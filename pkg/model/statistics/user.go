package statistics

import (
	"github.com/coralproject/pillar/pkg/model"
)

type UserActions struct {
	Performed *ActionDimensions "json:performed bson:performed"
	Received  *ActionDimensions "json:received bson:received"
}

type UserComments CommentDimensions

type UserStatistics struct {
	Actions  *UserActions  `json:"actions" bson:"actions"`
	Comments *UserComments `json:"comments" bson:"comments"`
}

type User struct {
	*model.User `bson:",inline"`
	Statistics  *UserStatistics `json:"stats" bson:"stats"`
}
