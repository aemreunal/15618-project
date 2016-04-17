package benchmark

import (
	"sync"
)

type LockMap struct {
	data map[interface{}] interface{}
	lock *sync.Mutex
}

func NewLockMap() *LockMap {
	return &LockMap{make(map[interface{}]interface{}), new (sync.Mutex)}
}

func (lockmap *LockMap) Get(k interface{}) (interface{}, bool) {
	lockmap.lock.Lock()
	defer lockmap.lock.Unlock()
	v, ok := lockmap.data[k]
	return v,ok
}

func (lockmap *LockMap) Put(k,v interface{}) interface{} {
	lockmap.lock.Lock()
	defer lockmap.lock.Unlock()
	/* Save old value */
	old, _ := lockmap.data[k]
	lockmap.data[k] = v
	return old
}

func (lockmap *LockMap) Remove(k interface{}) (interface{}, bool){
	lockmap.lock.Lock()
	defer lockmap.lock.Unlock()
	/* Save old value */
	old, ok := lockmap.data[k];
	if ok {
		delete(lockmap.data, k)
	}
	return old, ok
}

func (lockmap *LockMap) clear() {
	lockmap.lock.Lock()
	defer lockmap.lock.Unlock()
	lockmap.data = make(map[interface{}] interface{})
}
