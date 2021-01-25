package collection

import "reflect"

type GoSlice []interface{}

func (rs GoSlice) ForEach(cb func(p interface{})) {
	for _, r := range rs {
		cb(r)
	}
}

func (rs GoSlice) MapToStrList(cb func(p interface{}) string) []string {
	res := make([]string, 0)
	for _, r := range rs {
		sr := cb(r)
		res = append(res, sr)
	}
	return res
}

func (rs GoSlice) Map(cb func(p interface{}) interface{}) GoSlice {
	res := make(GoSlice, 0)
	for _, r := range rs {
		sr := cb(r)
		res = append(res, sr)
	}
	return res
}

func (rs GoSlice) Filter(cb func(p interface{}) bool) GoSlice {
	res := make(GoSlice, 0)
	for _, r := range rs {
		if ok := cb(r); ok {
			res = append(res, r)
		}
	}
	return res
}

func (rs GoSlice) FindIndex(cb func(p interface{}) bool) int {
	for idx, r := range rs {
		if ok := cb(r); ok {
			return idx
		}
	}
	return -1
}

func NewGoSlice(il interface{}) (ret GoSlice) {
	ret = make(GoSlice, 0)
	switch reflect.TypeOf(il).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(il)
		for i := 0; i < s.Len(); i++ {
			ret = append(ret, s.Index(i).Interface())
		}
	default:
		ret = nil
	}
	return
}
