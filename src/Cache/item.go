package Cache

import (
	"time"
)

type Item struct {
	Object		interface{}
	Expiration	int64
}

// 判断item是否过期
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}