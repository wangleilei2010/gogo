package types

import (
	"fmt"
	"reflect"
)

type Slice struct {
	innerSlice []interface{}
}

func (rs Slice) ForEach(cb func(p interface{})) {
	for _, r := range rs.innerSlice {
		cb(r)
	}
}

func (rs Slice) Len() int {
	return len(rs.innerSlice)
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

func (rs Slice) IndexOf(p interface{}) int {
	for idx, r := range rs.innerSlice {
		if r == p {
			return idx
		}
	}
	return -1
}

func (rs Slice) GetAt(index int) interface{} {
	for idx, r := range rs.innerSlice {
		if idx == index {
			return r
		}
	}
	return nil
}

func (rs Slice) Contains(p interface{}) bool {
	for _, r := range rs.innerSlice {
		if r == p {
			return true
		}
	}
	return false
}

func (rs *Slice) Append(v interface{}) {
	rs.innerSlice = append(rs.innerSlice, v)
}

func (rs *Slice) Extend(il interface{}) {
	switch reflect.TypeOf(il).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(il)
		for i := 0; i < s.Len(); i++ {
			rs.Append(s.Index(i).Interface())
		}
	default:
		fmt.Println("[WARNING]Unsupported Type:", reflect.TypeOf(il).Kind())
	}
}

func (rs *Slice) Remove(v interface{}) {
	idx := rs.IndexOf(v)
	if idx != -1 {
		rs.innerSlice = append(rs.innerSlice[:idx], rs.innerSlice[idx+1:]...)
	}
}

func (rs *Slice) RemoveAll(v interface{}) {
	for {
		idx := rs.IndexOf(v)
		if idx != -1 {
			rs.innerSlice = append(rs.innerSlice[:idx], rs.innerSlice[idx+1:]...)
		} else {
			break
		}
	}
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
