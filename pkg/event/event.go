package event

//Various event related constants
const (
	EventCommentNew     string = "event_comment_new"
	EventCommentUpdated string = "event_comment_updated"
	EventActionNew      string = "event_action_new"
)

// Event denotes an event in Pillar
type Event struct {
	Name   string   `json:"name"
	Payload bson.ObjectId `json:"target_id"
}
