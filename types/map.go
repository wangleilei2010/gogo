package types

type Map struct {
	innerMap map[interface{}]interface{}
}

type MapEntry struct {
	Key   interface{}
	Value interface{}
}

func (m *Map) Get(key interface{}) interface{} {
	if v, ok := m.innerMap[key]; ok {
		return v
	} else {
		return nil
	}
}

func (m *Map) Set(key, value interface{}) {
	m.innerMap[key] = value
}

func (m *Map) ContainsKey(key interface{}) bool {
	if _, ok := m.innerMap[key]; ok {
		return true
	} else {
		return false
	}
}

func (m *Map) Keys() *Slice {
	ret := make([]interface{}, 0)
	for k, _ := range m.innerMap {
		ret = append(ret, k)
	}
	return ConvSlice(ret)
}

func (m *Map) Values() *Slice {
	ret := make([]interface{}, 0)
	for _, v := range m.innerMap {
		ret = append(ret, v)
	}
	return ConvSlice(ret)
}

func (m *Map) EntrySet() []MapEntry {
	ret := make([]MapEntry, 0)
	for k, v := range m.innerMap {
		entry := MapEntry{Key: k, Value: v}
		ret = append(ret, entry)
	}
	return ret
}

func (m *Map) ContainsValue(value interface{}) bool {
	values := m.Values()
	for _, e := range values.innerSlice {
		if e == value {
			return true
		}
	}
	return false
}

func NewMap() *Map {
	m := make(map[interface{}]interface{})
	iMap := &Map{innerMap: m}
	return iMap
}
