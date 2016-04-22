package amqp

import (
	"github.com/streadway/amqp"
)

//MQ denotes a wrapper structure around amqp.Exchange and amqp.Channel
type MQ struct {
	Exchange string
	Channel  *amqp.Channel
}

var (
	amqpConnection *amqp.Connection
)

func connect(url string) *amqp.Connection {

	if amqpConnection != nil {
		return amqpConnection
	}

	c, err := amqp.Dial(url)
	if err != nil {
		return nil
	}

	amqpConnection = c
	return amqpConnection
}


func NewMQ(url string, exchange string) *MQ {

	//create an MQ anyway
	mq := MQ{exchange, nil}

	conn := connect(url)
	if conn == nil {
		return &mq
	}

	ch, _ := conn.Channel()
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
	if m.Channel == nil {
		return
	}

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
