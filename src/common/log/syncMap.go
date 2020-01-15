//==============================================================================
// syncMap Map锁
// -- 自定义带锁的Map
//==============================================================================
package log

import (
	"sync"
)

type SyncMap struct {
	lock *sync.RWMutex
	bm   map[interface{}]interface{}
}

func NewSyncmap() *SyncMap {
	return &SyncMap{
		lock: new(sync.RWMutex),
		bm:   make(map[interface{}]interface{}),
	}
}

func (m *SyncMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

func (m *SyncMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

func (m *SyncMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

func (m *SyncMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

func (m *SyncMap) getMap() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.bm
}
