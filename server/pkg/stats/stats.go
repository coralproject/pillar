// The stats package is responsible for updating relevant cached statistics when events are fired

package stats

import (
	"errors"

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
	Payload map[string]string
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

			if event.Name == "entity.comment.create" {
				onCreateComment(event.Payload)
			}

		}

	}

}
