package Redis

import (
	"time"
	"github.com/gomodule/redigo/redis"
)

type Option struct {
	MaxIdle		int
	IdleTimeout	time.Duration
	Password	string
	Database	int
	Host		string
}

type Client struct {
	pool		*redis.Pool
}

func SetOption(o *Option) (cli *Client, err error) {

	pool := &redis.Pool {
		MaxIdle: o.MaxIdle,	// 最大空闲连接数
		MaxActive: o.MaxIdle,	// 最大活跃连接处
		IdleTimeout: o.IdleTimeout,	// 超时时间
		Dial: func() (redis.Conn, error) {
			pw := redis.DialPassword(o.Password)
			db := redis.DialDatabase(o.Database)

			conn, e := redis.Dial("tcp", o.Host, pw, db)
			if e != nil {
				return nil, e
			}

			return conn, nil
		},
	}

	c := &Client {
		pool: pool,
	}

	return c, nil
}

func (c *Client) Conn() redis.Conn {
	return c.pool.Get()
}

func (c *Client) Dispose() {
	c.pool.Close()
}

//

func (c *Client) Count() (int, error) {
	return redis.Int(c.Conn().Do("DBSIZE"))
}

func (c *Client) Keys() ([]string, error) {
	return redis.Strings(c.Conn().Do("KEYS", "*"))
}

// 基础操作

func (c *Client) ContainsKey(key string) (bool, error) {
	return  redis.Bool(c.Conn().Do("EXISTS", key))
}

func (c *Client) SetString(key, value string, second int) (bool, error) {
	var res interface{}
	var err error

	if (second <= 0) {
		res, err = c.Conn().Do("SET", key, value)
	} else {
		res, err = c.Conn().Do("SETEX", key, second, value)
	}

	s, e := redis.String(res, err);
	if e != nil {
		return false, e
	}
	return s == "OK", nil
}

func (c *Client) GetString(key string) (string, error) {
	return redis.String(c.Conn().Do("GET", key))
}

func (c *Client) Remove(keys ...string) (int, error) {
	return redis.Int(c.Conn().Do("DEL", keys))
}

func (c *Client) SetExpire(key string, second int) (bool, error) {
	r, e := redis.String(c.Conn().Do("EXPIRE", key, second))
	if e != nil {
		return false, e
	}
	return r == "OK", nil
}

func (c *Client) GetExpire(key string) (int, error) {
	r, e := redis.Int(c.Conn().Do("TTL", key))
	if e != nil {
		return 0, e
	}
	return r, nil
}

// 集合操作

func (c *Client) SetAll(values map[string]string) (bool, error) {
	r, e := redis.String(c.Conn().Do("MSET", values))
	if e != nil {
		return false, e
	}
	return r == "OK", nil
}

func (c *Client) GetAll(keys []string) (map[string]string, error) {
	return redis.StringMap(c.Conn().Do("MGET", keys))
}

// 高级操作

func (c *Client) Add(key, value string) (bool, error) {
	// SETNX 不具备设置过期时间，考虑使用EXPIRE
	return redis.Bool(c.Conn().Do("SETNX", key, value))
}

func (c *Client) Replace(key, value string) (string, error) {
	return redis.String(c.Conn().Do("GETSET", key, value))
}

func (c *Client) Increment(key string, value int) (int, error) {
	if value == 1 {
		return redis.Int(c.Conn().Do("INCR", key))
	}
	return redis.Int(c.Conn().Do("INCRBY", key, value))
}

func (c *Client) IncrementFloat(key string, value float64) (float64, error) {
	return redis.Float64(c.Conn().Do("INCRBYFLOAT", key, value))
}

func (c *Client) Decrement(key string, value int) (int, error) {
	if value == 1 {
		return redis.Int(c.Conn().Do("DECR", key, value))
	}
	return redis.Int(c.Conn().Do("DECRBY", key, value))
}

func (c *Client) DecrementFloat(key string, value float64) (float64, error) {
	return c.IncrementFloat(key, -value)
}

// HASH

func (c *Client) GetHash(key, hashKey string) (string, error) {
	return redis.String(c.Conn().Do("HGET", key, hashKey))
}

func (c *Client) SetHash(key, hashKey, value string) (bool, error) {
	return redis.Bool(c.Conn().Do("HSET", key, hashKey, value))
}

func (c *Client) GetInfo() (string, error) {
	return redis.String(c.Conn().Do("INFO"))
}