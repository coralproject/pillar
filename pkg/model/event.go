package model

//Various Events
const (
	EventUserImport       string = "user_import"
	EventAssetImport      string = "asset_import"
	EventActionImport     string = "action_import"
	EventCommentImport    string = "comment_import"
	EventNoteImport       string = "note_import"

	EventUserAddUpdate    string = "user_add_update"
	EventAssetAddUpdate   string = "asset_add_update"
	EventActionAddUpdate  string = "action_add_update"
	EventCommentAddUpdate string = "comment_add_update"
	EventNoteAddUpdate    string = "note_add_update"
	EventSearchAddUpdate  string = "search_add_update"
	EventAuthorAddUpdate  string = "author_add_update"
	EventSectionAddUpdate string = "section_add_update"

	EventTagAdded         string = "tag_added"
	EventTagRemoved       string = "tag_removed"
)

//Event denotes an event in Pillar
type Event struct {
	Name    string      `json:"name" bson:"name"`
	Payload interface{} `json:"payload" bson:"payload"`
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
	Tag  string `json:"tag" bson:"tag"`
	User User   `json:"user" bson:"user"`
}
