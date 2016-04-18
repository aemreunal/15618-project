package gotomic

/* Encapsulation of Hash */

type GotomicMap struct {
	hash *Hash
}

func NewGotomicMap() *GotomicMap {
	return &GotomicMap{NewHash()}
}

func (this *GotomicMap) GetHashableKey(k interface{}) Hashable {
	var key Hashable
	switch k.(type) {
	case int:
		key = IntKey(k.(int))
		break
	case int64:
		key = Int64Key(k.(int64))
		break
	case string:
		key = StringKey(k.(string))
		break
	}
	return key
}

func (this *GotomicMap) Get(k interface{}) (interface{}, bool) {
	return this.hash.Get(this.GetHashableKey(k))
}

func (this *GotomicMap) Put(k, v interface{}) interface{} {
	old, _ := this.hash.Put(this.GetHashableKey(k), v)
	return old
}

func (this *GotomicMap) Remove(k interface{}) (interface{}, bool) {
	return this.hash.Delete(this.GetHashableKey(k))
}
