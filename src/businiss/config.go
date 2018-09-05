package businiss

import (
	"Redis"
	"RabbitMQ"
)

type config struct {
	RabbitClient	*RabbitMQ.Client
	RabbitServer	*RabbitMQ.Server
	Redis			*Redis.Client
}

var c = &config{}

func SetRabbitMQClient(mqOption *RabbitMQ.Option, secret, method string) error {

	if secret != "" && method == "" {
		method = "SHA256"
	}

	r, e := RabbitMQ.SetConfig(mqOption).BuildRabbitMQClient(secret, method)
	if e != nil {
		return e
	}
	c.RabbitClient = r
	return nil
}

func SetRabbitMQServer(mqOption *RabbitMQ.Option, secret, method string) error {

	if secret != "" && method == "" {
		method = "SHA256"
	}

	r, e := RabbitMQ.SetConfig(mqOption).BuildRabbitMQReceiver(secret, method)
	if e != nil {
		return e
	}
	c.RabbitServer = r
	return nil
}

func SetRedisConfig(redisOption *Redis.Option) error {
	r, e := Redis.SetOption(redisOption)
	if e != nil {
		return e
	}
	c.Redis = r
	return nil
}

func GetRabbitClient() *RabbitMQ.Client {
	return c.RabbitClient
}

func GetRabbitServer() *RabbitMQ.Server {
	return c.RabbitServer
}

func GetRedis() *Redis.Client {
	return c.Redis
}