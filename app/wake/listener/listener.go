package listener

import (
	"encoding/json"
	"fmt"
	"github.com/WPMedia/coral/wake/worker"
	"github.com/coralproject/pillar/pkg/amqp"
	"github.com/coralproject/pillar/pkg/model"
	"log"
	"os"
)

func Listen() {
	mq := amqp.NewMQ(os.Getenv("AMQP_URL"), os.Getenv("AMQP_EXCHANGE"))
	defer mq.Close()

	msgs, err := mq.Receive()
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {

			var event model.Event
			json.Unmarshal(d.Body, &event)
			switch event.Name {
			case model.EventAssetImport, model.EventAssetAddUpdate:
				worker.UpdateAsset(event)
				break

			case model.EventTagAdded, model.EventTagRemoved:
				worker.UpdateUserTag(event)
				break

			default:
				break
			}
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
