package generics

import (
	"errors"
	"fmt"
)

type Slice[T any] []T

func NewSlice[T any]() *Slice[T] {
	s := make(Slice[T], 0)
	return &s
}

func (s *Slice[T]) Append(t ...T) {
	*s = append(*s, t...)
}

func (s *Slice[T]) Len() int {
	return len(*s)
}

func (s *Slice[T]) Foreach(f func(t T)) {
	for _, e := range *s {
		f(e)
	}
}

func (s *Slice[T]) FindIndex(compare func(t T) bool) int {
	for i, e := range *s {
		if compare(e) {
			// 由于T的范围是any, 无法简单地判断t==e, 故需传入比较函数
			return i
		}
	}
	return -1
}

func (s *Slice[T]) Remove(compare func(t T) bool) {
	for {
		idx := s.FindIndex(compare)
		if idx == -1 {
			break
		} else {
			*s = append((*s)[:idx], (*s)[idx+1:]...)
		}
	}
}

func (s *Slice[T]) Get(index int) (notOutOfBounds bool, t T) {
	for i, e := range *s {
		if i == index {
			notOutOfBounds = true
			t = e
			return
		}
	}
	// 越界
	notOutOfBounds = false
	return
}

func (s *Slice[T]) Slice(start, end int) *Slice[T] {
	ns := NewSlice[T]()
	for i, e := range *s {
		if i >= start && i < end {
			ns.Append(e)
		}
	}
	return ns
}

func (s *Slice[T]) Filter(match func(t T) bool) *Slice[T] {
	ns := NewSlice[T]()
	for _, e := range *s {
		if match(e) {
			ns.Append(e)
		}
	}
	return ns
}

func (s *Slice[T]) Update(pos int, t T) error {
	if pos < 0 || pos > s.Len()-1 {
		return errors.New(fmt.Sprintf("OUT_OF_BOUNDS:%d", pos))
	} else {
		(*s)[pos] = t
		return nil
	}
}

func SliceMapping[E, V any](s *Slice[E], mapFunc func(e E) V) *Slice[V] {
	ns := NewSlice[V]()
	for _, e := range *s {
		v := mapFunc(e)
		ns.Append(v)
	}
	return ns
}
