package async

import "sync"

const defaultConMapCap = 32

// ConMap 并发安全Map
// 仅用于 `写远多于读` 的场景，对于 `并发读` 操作较多的场景，建议用官方 `sync.Map`
type ConMap struct {
	locker *sync.RWMutex
	items  map[string]interface{}
}

// NewConMap 创建 ConMap
func NewConMap() *ConMap {
	return NewConMapWithCap(defaultConMapCap)
}

// NewConMapWithCap 创建指定 cap 的 ConMap
func NewConMapWithCap(cap int) *ConMap {
	if cap <= 0 {
		cap = defaultConMapCap
	}

	return &ConMap{
		locker: &sync.RWMutex{},
		items:  make(map[string]interface{}, cap),
	}
}

// Set 设置 key/value 对
func (m *ConMap) Set(key string, value interface{}) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.items[key] = value
}

// Get 根据 key 获取 value
func (m *ConMap) Get(key string) interface{} {
	m.locker.Lock()
	defer m.locker.Unlock()

	return m.items[key]
}

// SetValues 批量设置 key/value
func (m *ConMap) SetValues(data map[string]interface{}) {
	m.locker.Lock()
	defer m.locker.Unlock()

	for k, v := range data {
		m.items[k] = v
	}
}

// GetValues 根据 key 列表，获取指定key 的值
func (m *ConMap) GetValues(keys []string) map[string]interface{} {
	m.locker.RLock()
	defer m.locker.RUnlock()

	// 定义结果
	l := len(keys)
	if l > len(m.items) {
		l = len(m.items)
	}
	r := make(map[string]interface{}, l)

	for _, key := range keys {
		if value, exist := m.items[key]; exist {
			r[key] = value
		}
	}

	return r
}

// Data 获取所有值
func (m *ConMap) Data() map[string]interface{} {
	m.locker.RLock()
	defer m.locker.RUnlock()

	r := make(map[string]interface{}, len(m.items))
	for key, value := range m.items {
		r[key] = value
	}

	return r
}
