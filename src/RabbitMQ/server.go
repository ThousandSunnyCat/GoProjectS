package RabbitMQ

import (
	"encoding/json"
	"time"
	"fmt"
	"sync"
	"github.com/streadway/amqp"
)

type Server struct {
	wg					sync.WaitGroup

	secret				string
	method				string

	connection			*amqp.Connection
	channels			map[string]*amqp.Channel
}

type Receiver interface {
	QueueName()			string
	RouterKeys()		[]string
	Retry()				int
	OnError(error)
	OnReceive(*Message)	bool
}

// 绑定连接
func (r *Server) Bind(exchangeName, exchangeType string, receivers []Receiver) (err error) {
	defer func() {
		r.Dispose()
		if e := recover(); e != nil {
			err = fmt.Errorf("队列绑定到虚拟机异常: %s", e)
		}
	}()

	channel, e := r.getChannel(exchangeName)
	if e != nil {
		return e
	}

	// 定义交换机(需要等待返回error)
	if e := channel.ExchangeDeclare(exchangeName, exchangeType, false, false, false, false, nil); e != nil {
		return e
	}
	
	// 绑定监听队列
	for _, v := range receivers {
		r.wg.Add(1)
		go r.listen(channel, v, exchangeName)
	}
	
	r.wg.Wait()

	return
}

// 清理连接
func (r *Server) Dispose() {
	for _, v := range r.channels {
		v.Close()
	}
	r.connection.Close()
}

// 创建消费者
func (r *RabbitMQ) createReceiver(secret, method string) (client *Server, err error) {
	// config
	connection, e := amqp.Dial(r.uri)
	if e != nil {
		return nil, e
	}

	var res = &Server {
		secret: secret,
		method: method,
		connection: connection,
		channels: make(map[string]*amqp.Channel),
	}
	
	return res, nil
}

func (r *Server) getChannel(key string) (channel *amqp.Channel, err error) {
	
	channel, ok := r.channels[key]
	if !ok {
		ch, e := r.connection.Channel()
		if e != nil {
			return nil, e
		}

		r.channels[key] = ch
		channel = ch
	}

	return channel, nil
}

func (r *Server) listen(channel *amqp.Channel, v Receiver, exchangeName string) {
	// 组锁
	defer r.wg.Done()

	queueName := v.QueueName()

	// 定义队列(可以不等待?)
	if _, e := channel.QueueDeclare(queueName, true, false, false, false, nil); e != nil {
		v.OnError(fmt.Errorf("队列创建失败: %s", e))
		return
	}

	// 消息投递设置
	if e := channel.Qos(1, 0, false); e != nil {
		v.OnError(fmt.Errorf("投递设置失败: %s", e))
		return
	}

	// 绑定队列交换机
	for _, key := range v.RouterKeys() {
		if e := channel.QueueBind(queueName, key, exchangeName, false, nil); e != nil {
			v.OnError(fmt.Errorf("队列交换机绑定失败: %s", e))
			return
		}
	}

	// 获取消费通道(可启动AutoACK)
	msgs, e := channel.Consume(queueName, "", false, false, false, false, nil)
	if e != nil {
		v.OnError(fmt.Errorf("消费通道获取失败: %s", e))
		return
	}
	// 处理方法
	for msg := range msgs {
		m, e := r.getMessage(msg.Body)
		if e != nil {
			v.OnError(fmt.Errorf("消息解码失败: %s", e))
		} else {

			// 验证签名
			if r.method != "" && r.secret != "" && m.Signature != m.Sign(r.secret, r.method) {
				v.OnError(fmt.Errorf("验签失败: %s", e))
				goto MESSAGEACK
			}

			// 当接收者消息处理失败的时候，
			// 比如网络问题导致的数据库连接失败，redis连接失败等等这种
			// 通过重试可以成功的操作，那么这个时候是需要重试的
			// 直到数据处理成功后再返回，然后才会回复rabbitmq ack
			// 重试次数由消费者自身决定

			times := v.Retry()
			for !v.OnReceive(m) {
				//log.Warnf("receiver 数据处理失败，将要重试")
				time.Sleep(500 * time.Millisecond)

				if times <= 0 {
					goto MESSAGEACK
				}

				times--
			}
		}
		// goto 标签 跳出读消息循环
		MESSAGEACK:

		// 确认收到本条消息, multiple必须为false
		msg.Ack(false)
	}
}

func (r *Server) getMessage(body []byte) (*Message, error) {
	// 转json
	m := &Message{}
	if e := json.Unmarshal(body, m); e != nil {
		return nil, e
	}

	return m, nil
}