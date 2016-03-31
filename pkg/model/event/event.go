package event

import (
	"github.com/coralproject/pillar/pkg/model"
)

//PayloadComment denotes a payload to be used when a comment is created/updated.
type PayloadComment struct {
	Comment model.Comment `json:"comment" bson:"comment"`
	Asset   model.Asset   `json:"asset" bson:"asset"`
	User    model.User    `json:"user" bson:"user"`
}

//PayloadAction denotes an message to be used when an action is created/updated.
type PayloadAction struct {
	Action  model.Action  `json:"action" bson:"action"`
	Actor   model.User    `json:"actor" bson:"actor"`
	User    model.User    `json:"user" bson:"user"`
	Comment model.Comment `json:"comment" bson:"comment"`
}
