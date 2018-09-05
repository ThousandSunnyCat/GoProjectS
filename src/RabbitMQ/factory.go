package RabbitMQ

import (
	"fmt"
)

// 连接配置
type Option struct {
	HostName	string
	Port		int
	UserName	string
	Password	string
	VirtualHost	string
}

// 如果有其他配置，往这里添加
type RabbitMQ struct {
	uri			string	// RabbitMQ 连接字符串
}

type Message struct {
	MessageId	string	// 消息ID
	IsDurable	bool	// 是否保存
	MessageBody	string	// 消息体
	CreateTime	int64		// 消息创建时间
	Signature	string	// 消息签名
}

// 消息队列设置
func SetConfig(o *Option) *RabbitMQ {
	return &RabbitMQ{
		uri: fmt.Sprintf("amqp://%s:%s@%s:%d/%s", o.UserName, o.Password, o.HostName, o.Port, o.VirtualHost),
	}
}

// 客户端(消息生产者)
func (r *RabbitMQ) BuildRabbitMQClient(secret, method string) (client *Client, err error) {
	defer func() {
		if x := recover(); x != nil {
			client = nil
			err = fmt.Errorf("Error occurred in create RabbitMQClient: %v", x)
		}
	}()

	return r.createClient(secret, method)
}

// 服务端(消息消费者)
func (r *RabbitMQ) BuildRabbitMQReceiver(secret, method string) (receiver *Server, err error) {
	defer func() {
		if x := recover(); x != nil {
			receiver = nil
			err = fmt.Errorf("Error occurred in create RabbitMQClient: %v", x)
		}
	}()

	return r.createReceiver(secret, method)
}