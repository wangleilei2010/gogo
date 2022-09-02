package generics

type MapKey interface {
	comparable
}

type GenericMap[K MapKey, V any] map[K]V

func NewMap[K MapKey, V any]() GenericMap[K, V] {
	return make(map[K]V)
}

func (m GenericMap[K, V]) ContainsKey(k K) bool {
	if m == nil {
		return false
	}
	if _, ok := m[k]; ok {
		return true
	}
	return false
}

func (m GenericMap[K, V]) DelKey(k K) {
	delete(m, k)
}

func (m GenericMap[K, V]) Put(k K, v V) {
	m[k] = v
}

func (m GenericMap[K, V]) Keys() *Slice[K] {
	var keys = NewSlice[K]()
	for k, _ := range m {
		keys.Append(k)
	}
	return keys
}

func (m GenericMap[K, V]) Values() *Slice[V] {
	var values = NewSlice[V]()
	for _, v := range m {
		values.Append(v)
	}
	return values
}
