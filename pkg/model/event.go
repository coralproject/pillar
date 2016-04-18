package model

//Various Events
const (
	EventAssetAdded     string = "asset_added"
	EventAssetUpdated   string = "asset_updated"
	EventCommentAdded   string = "comment_added"
	EventCommentUpdated string = "comment_updated"
	EventTagAdded       string = "tag_added"
	EventTagRemoved     string = "tag_removed"
)

//PayloadAsset denotes an message to be used when an Asset is added/updated
type PayloadAsset struct {
	Event string `json:"event" bson:"event"`
	Asset Asset  `json:"asset" bson:"asset"`
}

//PayloadComment denotes a payload to be used when a comment is created/updated.
type PayloadComment struct {
	Comment Comment `json:"comment" bson:"comment"`
	Asset   Asset   `json:"asset" bson:"asset"`
	User    User    `json:"user" bson:"user"`
}

//PayloadAction denotes an message to be used when an action is created/updated.
type PayloadAction struct {
	Action  Action  `json:"action" bson:"action"`
	Actor   User    `json:"actor" bson:"actor"`
	User    User    `json:"user" bson:"user"`
	Comment Comment `json:"comment" bson:"comment"`
}

//PayloadTag denotes an message to be used when a tag is added/removed
type PayloadTag struct {
	Event string `json:"event" bson:"event"`
	Tag   string `json:"tag" bson:"tag"`
	User  User   `json:"user" bson:"user"`
}
