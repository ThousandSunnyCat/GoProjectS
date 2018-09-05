package RabbitMQ

import (
	"fmt"
	"time"
	"encoding/json"
	"sync"
	"github.com/streadway/amqp"
	"github.com/satori/go.uuid"
)

type Client struct {
	wg sync.WaitGroup

	secret			string
	method			string

	connection		*amqp.Connection
	channel			*amqp.Channel
}

func (r *RabbitMQ) createClient(secret, method string) (client *Client, err error) {
	// config
	connection, e := amqp.Dial(r.uri)
	if e != nil {
		return nil, e
	}

	channel, e := connection.Channel()
	if e != nil {
		return nil, e
	}

	var res = &Client {
		connection: connection,
		channel: channel,
		secret: secret,
		method: method,
	}
	
	return res, nil
}

func (r *Client) Send(exchangeName, routerKey, mbody string) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %v", e)
		}
	}()

	// 生成message对象
	m := &Message {
		MessageId: uuid.Must(uuid.NewV4()).String(),
		MessageBody: mbody,
		CreateTime: time.Now().Unix(),
		IsDurable: false,
	}

	m.Signature = m.Sign(r.secret, r.method)

	// struct转json
	body, e := json.Marshal(m)
	if e != nil {
		return e
	}

	if e := r.channel.Publish(exchangeName, routerKey, false, false, amqp.Publishing {
		ContentType: "application/json",
		Body: []byte(body),
		DeliveryMode: amqp.Transient,
		Priority: 0,
	}); e != nil {
		return e
	}

	return nil
}

func (r *Client) Dispose() {
	defer func() {
		if err := recover(); err != nil {
			// log panic
		}
	}()

	r.channel.Close()
	r.connection.Close()
}