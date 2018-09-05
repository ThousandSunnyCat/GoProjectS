package businiss

import (
	"fmt"
	"time"
	"RabbitMQ"
)

type receiver struct {
	queueName	string
	routerKeys	[]string
	times		int
	errorFunc	func(error)
	success		func(string)bool
}
func (r *receiver) QueueName() string {
	return r.queueName
}
func (r *receiver) RouterKeys() []string {
	return r.routerKeys
}
func (r *receiver) OnError(err error) {
	// log
	r.errorFunc(err)
}
func (r *receiver) OnReceive(m *RabbitMQ.Message) bool {

	if m.CreateTime < time.Now().Unix() {
		// 公共操作
	}
	
	return r.success(m.MessageBody)
}
func (r *receiver) Retry() int {
	return r.times
}

func SetReceiver() {

	defer func(){
		if x := recover(); x != nil {
			fmt.Printf("SetReceiver panic err: %v", x)
		}
	}()

	option := &RabbitMQ.Option {
		HostName: "xxx.xxx.xxx.xxx",
		Port: 5673,
		UserName: "xxxxxx",
		Password: "xxxxxx",
		VirtualHost: "xxxxxx",
	}
	
	if e := SetRabbitMQServer(option, "Test111", "MD5"); e != nil {
		// log，异常处理
	}

	server := GetRabbitServer()
	receivers := []RabbitMQ.Receiver{}
	receivers = append(receivers, &receiver{
		queueName: "SendMessage",
		routerKeys: []string{"SendMessageService"},
		times: 3,
		errorFunc: func(e error) {
			// log
			fmt.Printf("SendMessage error: %v", e)
		},
		success: func(s string) bool {
			// 发送消息
			fmt.Printf("SendMessage: %v", s)
			//panic("test")
			return true
		},
	})

	go func() {
		if err := server.Bind("MessageService", "direct", receivers); err != nil {
			// log 处理异常
		}
	}()

	test()
}

func test(){
	option := &RabbitMQ.Option {
		HostName: "112.74.56.154",
		Port: 5673,
		UserName: "yingxing",
		Password: "waimaisaas2017",
		VirtualHost: "YingXingTest",
	}
	
	if e := SetRabbitMQClient(option, "Test111", "MD5"); e != nil {
		fmt.Printf("client: %v", e)
	}

	r := GetRabbitClient()
	if e := r.Send("MessageService", "SendMessageService", "golang testestest"); e != nil {
		fmt.Printf("senderr: %v", e)
	}
}