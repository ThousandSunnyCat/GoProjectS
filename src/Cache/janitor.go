package Cache

import (
	"time"
)

type janitor struct {
	Interval	time.Duration
	stop		chan bool
}

// 开启缓存的定时器
func (j *janitor) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <- ticker.C:
			c.DeleteExpired()
		case <- j.stop:
			ticker.Stop()
			return
		}
	}
}

// 停止定时器
func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

// 定时器工厂
func runJanitor(c *cache, ci time.Duration) {
	j := &janitor {
		Interval:	ci,
		stop:		make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}