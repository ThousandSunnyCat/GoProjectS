/*
\---------------------------------/
		需要重新考虑解决方案		
\---------------------------------/
*/

package Cache

import (
	"fmt"
)

func (c *cache) Increment(key string, n int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[key]
	if !found || v.Expired() {
		return fmt.Errorf("Item %s not found", key)
	}
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) + int(n)
	case int8:
		v.Object = v.Object.(int8) + int8(n)
	case int16:
		v.Object = v.Object.(int16) + int16(n)
	case int32:
		v.Object = v.Object.(int32) + int32(n)
	case int64:
		v.Object = v.Object.(int64) + n
	case uint:
		v.Object = v.Object.(uint) + uint(n)
	case uintptr:
		v.Object = v.Object.(uintptr) + uintptr(n)
	case uint8:
		v.Object = v.Object.(uint8) + uint8(n)
	case uint16:
		v.Object = v.Object.(uint16) + uint16(n)
	case uint32:
		v.Object = v.Object.(uint32) + uint32(n)
	case uint64:
		v.Object = v.Object.(uint64) + uint64(n)
	case float32:
		v.Object = v.Object.(float32) + float32(n)
	case float64:
		v.Object = v.Object.(float64) + float64(n)
	default:
		return fmt.Errorf("The value for %s is not an integer", key)
	}
	c.items[key] = v
	return nil
}

func (c *cache) IncrementFloat(key string, n float64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[key]
	if !found || v.Expired() {
		return fmt.Errorf("Item %s not found", key)
	}
	switch v.Object.(type) {
	case float32:
		v.Object = v.Object.(float32) + float32(n)
	case float64:
		v.Object = v.Object.(float64) + n
	default:
		return fmt.Errorf("The value for %s does not have type float32 or float64", key)
	}
	c.items[key] = v
	return nil
}

func (c *cache) IncrementInt(k string, n int) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementInt8(k string, n int8) (int8, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int8)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int8", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementInt16(k string, n int16) (int16, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int16)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int16", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementInt32(k string, n int32) (int32, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int32)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int32", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementInt64(k string, n int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(int64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int64", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementUint(k string, n uint) (uint, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementUintptr(k string, n uintptr) (uintptr, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uintptr)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uintptr", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementUint8(k string, n uint8) (uint8, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint8)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint8", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementUint16(k string, n uint16) (uint16, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint16)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint16", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementUint32(k string, n uint32) (uint32, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint32)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint32", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementUint64(k string, n uint64) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(uint64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an uint64", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementFloat32(k string, n float32) (float32, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(float32)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an float32", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}

func (c *cache) IncrementFloat64(k string, n float64) (float64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.items[k]
	if !found || v.Expired() {
		return 0, fmt.Errorf("Item %s not found", k)
	}
	rv, ok := v.Object.(float64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an float64", k)
	}
	nv := rv + n
	v.Object = nv
	c.items[k] = v
	return nv, nil
}