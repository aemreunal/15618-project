package benchmark

import (
	"sync"
)

type RWLockMap struct {
	data map[interface{}] interface{}
	lock *sync.RWMutex
}

func NewRWLockMap() *RWLockMap {
	return &RWLockMap{make(map[interface{}]interface{}), new (sync.RWMutex)}
}

func (rwlockmap *RWLockMap) Get(k interface{}) (interface{}, bool) {
	rwlockmap.lock.RLock()
	defer rwlockmap.lock.RUnlock()
	v, ok := rwlockmap.data[k]
	return v,ok
}

func (rwlockmap *RWLockMap) Put(k,v interface{}) interface{} {
	rwlockmap.lock.Lock()
	defer rwlockmap.lock.Unlock()
	/* Save old value */
	old, _ := rwlockmap.data[k]
	rwlockmap.data[k] = v
	return old
}

func (rwlockmap *RWLockMap) Remove(k interface{}) (interface{}, bool){
	rwlockmap.lock.Lock()
	defer rwlockmap.lock.Unlock()
	/* Save old value */
	old, ok := rwlockmap.data[k];
	if ok {
		delete(rwlockmap.data, k)
	}
	return old, ok
}

func (rwlockmap *RWLockMap) Clear() {
	rwlockmap.lock.Lock()
	defer rwlockmap.lock.Unlock()
	rwlockmap.data = make(map[interface{}] interface{})
}
