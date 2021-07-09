package types

import "reflect"

type Slice struct {
	innerSlice []interface{}
}

func (rs Slice) ForEach(cb func(p interface{})) {
	for _, r := range rs.innerSlice {
		cb(r)
	}
}

func (rs Slice) MapToStrList(cb func(p interface{}) string) []string {
	res := make([]string, 0)
	for _, r := range rs.innerSlice {
		sr := cb(r)
		res = append(res, sr)
	}
	return res
}

func (rs Slice) Map(cb func(p interface{}) interface{}) *Slice {
	res := NewSlice()
	for _, r := range rs.innerSlice {
		sr := cb(r)
		res.Append(sr)
	}
	return res
}

func (rs Slice) Filter(cb func(p interface{}) bool) *Slice {
	res := NewSlice()
	for _, r := range rs.innerSlice {
		if ok := cb(r); ok {
			res.Append(r)
		}
	}
	return res
}

func (rs Slice) FindIndex(cb func(p interface{}) bool) int {
	for idx, r := range rs.innerSlice {
		if ok := cb(r); ok {
			return idx
		}
	}
	return -1
}

func (rs *Slice) Append(v interface{}) {
	rs.innerSlice = append(rs.innerSlice, v)
}

func ConvSlice(il interface{}) (ret *Slice) {
	ret = NewSlice()
	switch reflect.TypeOf(il).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(il)
		for i := 0; i < s.Len(); i++ {
			ret.Append(s.Index(i).Interface())
		}
	default:
		ret = nil
	}
	return
}

func NewSlice() *Slice {
	i := make([]interface{}, 0)
	return &Slice{innerSlice: i}
}
