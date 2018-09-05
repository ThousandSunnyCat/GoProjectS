package Cache

import (
	"os"
	"encoding/gob"
	"io"
	"fmt"
	"sync"
	"time"
)

const (
	NoExpiration time.Duration = -1

	DefaultExpiration time.Duration = 0
)

type Cache struct {
	*cache
}

type cache struct {
	defaultExpiration	time.Duration
	items				map[string]Item
	mu					sync.RWMutex
	onEvicted			func(string, interface{})
	janitor				*janitor
}

type keyAndValue struct {
	key		string
	value	interface{}
}

// 设置（覆盖）
func (c *cache) Set(key string, value interface{}, t time.Duration) {
	e := c.expiration(t)

	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value, e)
}

// 设置默认时间（覆盖
func (c *cache) SetDefault(key string, value interface{}) {
	c.Set(key, value, DefaultExpiration)
}

// 添加（不覆盖）
func (c *cache) Add(key string, value interface{}, t time.Duration) error {
	c.mu.Lock()
	defer c.mu.Lock()

	_, found := c.get(key)
	if found {
		return fmt.Errorf("Item %s already exists", key)
	}

	c.set(key, value, c.expiration(t))
	return nil
}

// 替换
func (c *cache) Replace(key string, value interface{}, t time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.get(key)
	if !found {
		return fmt.Errorf("Item %s doesn't exist", key)
	}
	c.set(key, value, c.expiration(t))
	return nil
}

// 获取
func (c *cache) Get(key string) (interface{}, bool) {

	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.get(key)
}

// 暂无好的整合方案
// 获取（带时间戳）
func (c *cache) GetWithExpiration(key string) (interface{}, time.Time, bool){
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, time.Time{}, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, time.Time{}, false
		}
		// unix时间，之后返回看看是否需要修改时区
		return item.Object, time.Unix(0, item.Expiration), true
	}

	return item.Object, time.Time{}, true
}

// 删除指定缓存
func (c *cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, evicted := c.delete(key)
	// 起协程，onEvicted不宜执行长时间阻塞的方法，避免协程暴涨问题
	if evicted {
		go func() {
			c.onEvicted(key, v)
		}()
	}
}

// 删除过期缓存
func (c *cache) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		// "Inlining" of expired
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}

	// 起协程，onEvicted不宜执行长时间阻塞的方法，避免协程暴涨问题
	go func() {
		for _, v := range evictedItems {
			c.onEvicted(v.key, v.value)
		}
	}()
}

// 设置缓存被清除时执行的方法（单个缓存单个方法）
func (c *cache) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.onEvicted = f
}

// 存储到文件
func (c *cache) SaveFile(fname string) (err error) {
	fp, e := os.Create(fname)
	if e != nil {
		return e
	}
	
	defer func() {
		e = fp.Close()
		if err == nil {
			err = e
		}
	}()
	
	return c.save(fp)
}

// 从文件恢复
func (c *cache) LoadFile(fname string) (err error) {
	fp, e := os.Open(fname)
	if e != nil {
		return e
	}

	defer func() {
		e := fp.Close()
		if err == nil {
			err = e
		}
	}()

	return c.load(fp)
}

// 获取整个缓存数据
func (c *cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m := make(map[string]Item, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			continue
		}
		m[k] = v
	}
	return m
}

// 查看缓存条数
func (c *cache) ItemCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// 清空缓存
func (c *cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[string]Item{}
}

/*
\----------------------------------------------/
					分	割	线					
\----------------------------------------------/
*/

func (c *cache) expiration(t time.Duration) int64 {
	var e int64
	// 默认值控制
	if t == DefaultExpiration {
		t = c.defaultExpiration
	}
	if t > 0 {
		e = time.Now().Add(t).UnixNano()
	}
	
	return e
}

func (c *cache) set(key string, value interface{}, e int64) {
	c.items[key] = Item {
		Object:		value,
		Expiration:	e,
	}
}

func (c *cache) get(key string) (interface{}, bool) {
	item, found := c.items[key]
	if !found || (item.Expiration > 0 && time.Now().UnixNano() > item.Expiration) {
		return nil, false
	}
	return item.Object, true
}

// 原delete
/*
func (c *cache) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}
*/
func (c *cache) delete(key string) (interface{}, bool) {
	if v, found := c.items[key]; found {
		delete(c.items, key)
		if c.onEvicted != nil {
			return v.Object, true
		}
	}
	return nil, false
}


func (c *cache) save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)

	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, v := range c.items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&c.items)
	return
}


func (c *cache) load(r io.Reader) (err error) {
	dec := gob.NewDecoder(r)

	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()

	items := map[string]Item{}
	err = dec.Decode(&items)
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range items {
		ov, found := c.items[k]
		if !found || ov.Expired() {
			c.items[k] = v
		}
	}

	return
}