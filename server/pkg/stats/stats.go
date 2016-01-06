// The stats package is responsible for updating relevant cached statistics when events are fired

package stats

import (
	"errors"

	"github.com/coralproject/pillar/server/model"

	"github.com/ardanlabs/kit/log"
)

/*


Event:

Add Comment
-- count user.stats.comments
-- append user.lists.comments

-- count asset.stats.comments
-- append asset.lists.comments

*/

type Message struct {
	Name    string
	Payload map[string]interface{}
}

var (
	Event chan Message
)

func initChan() {
	Event = make(chan Message, 1)
}

func Init() {
	initChan()

	go listen()

}

func listen() {

	for {
		select {

		case event, ok := <-Event:

			if ok != true {
				log.Error("stats", "chan", errors.New("stats channel closed"), "Error")
				return
			}

			// this event name requires a comment in the payload
			if event.Name == "entity.comment.create" {

				// validate that event.Payload["comment"] is of type model.Comment

				onCreateComment(event.Payload["comment"].(model.Comment))
			}

		}

	}

}
