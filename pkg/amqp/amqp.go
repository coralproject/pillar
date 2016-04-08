package amqp

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

const (
	defaultAMQP string = "amqp://guest:guest@localhost:5672/"
)

var (
	amqpConnection *amqp.Connection
)

func init() {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		log.Printf("$AMQP_URL not found, trying to connect locally [%s]", defaultAMQP)
		url = defaultAMQP
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("Error connecting to AMQP: %s", err)
	}

	//save the primary connection
	amqpConnection = conn
}

//MQ denotes a wrapper structure around amqp.Exchange and amqp.Channel
type MQ struct {
	Exchange string
	Channel  *amqp.Channel
}

func NewMQ(exchange string) *MQ {
	mq := MQ{exchange, nil}
	if amqpConnection == nil {
		return &mq
	}

	ch, _ := amqpConnection.Channel()
	if ch != nil {
		//declare exchange
		ch.ExchangeDeclare(
			exchange, // name
			"fanout", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		mq.Channel = ch
	}

	return &mq
}

func (m *MQ) Close() {
	m.Channel.Close()
}

func (m *MQ) IsValid() bool {
	return m.Channel != nil
}

func (m *MQ) Publish(payload []byte) error {
	return m.Channel.Publish(
		m.Exchange, // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		})
}

func (m *MQ) Receive() (<-chan amqp.Delivery, error) {

	q, err := m.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = m.Channel.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		m.Exchange, // exchange
		false,
		nil)
	if err != nil {
		return nil, err
	}

	return m.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}
