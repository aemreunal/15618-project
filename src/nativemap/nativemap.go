package nativemap

import ()

type NativeMap struct {
	data map[interface{}]interface{}
}

func NewNativeMap() *NativeMap {
	return &NativeMap{make(map[interface{}]interface{})}
}

func (nativemap *NativeMap) Get(k interface{}) (interface{}, bool) {
	v, ok := nativemap.data[k]
	return v, ok
}

func (nativemap *NativeMap) Put(k, v interface{}) interface{} {
	old := nativemap.data[k]
	nativemap.data[k] = v
	return old
}

func (nativemap *NativeMap) Remove(k interface{}) (interface{}, bool) {
	/* Save old value */
	old, ok := nativemap.data[k]
	if ok {
		delete(nativemap.data, k)
	}
	return old, ok
}

func (nativemap *NativeMap) clear() {
	nativemap.data = make(map[interface{}]interface{})
}
