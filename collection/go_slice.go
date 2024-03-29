package collection

type Slice[T any] []T

type MapKey interface {
	comparable
}

func (s *Slice[T]) Push(t T) {
	ns := make([]T, 0)
	for _, e := range *s {
		ns = append(ns, e)
	}
	ns = append(ns, t)
	*s = ns
}

func (s *Slice[T]) Append(t ...T) {
	ns := make([]T, 0)
	for _, e := range *s {
		ns = append(ns, e)
	}
	ns = append(ns, t...)
	*s = ns
}

func (s Slice[T]) Len() int {
	return len(s)
}

func (s Slice[T]) Foreach(f func(t T)) {
	for _, e := range s {
		f(e)
	}
}

func (s Slice[T]) MapToStrList(f func(t T) string) []string {
	l := make([]string, 0)
	for _, e := range s {
		m := f(e)
		l = append(l, m)
	}
	return l
}

func (s Slice[T]) Find(f func(t T) bool) Slice[T] {
	another := make(Slice[T], 0)
	for _, e := range s {
		if f(e) {
			another.Push(e)
		}
	}
	return another
}

func (s Slice[T]) FindOne(f func(t T) bool) T {
	for _, e := range s {
		if f(e) {
			return e
		}
	}
	return *new(T)
}

func (s Slice[T]) Contains(f func(t T) bool) bool {
	for _, e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

func Map[E, V any](s Slice[E], f func(e E) V) Slice[V] {
	d := make([]V, 0)
	for _, ele := range s {
		v := f(ele)
		d = append(d, v)
	}
	return d
}

func New[T any](a []T) *Slice[T] {
	s := &Slice[T]{}
	for _, i := range a {
		s.Push(i)
	}
	return s
}

//type ComparableSlice[T comparable] Slice[T]
//
//func (s *ComparableSlice[T]) Dedup() ComparableSlice[T] {
//	var n ComparableSlice[T]
//	var m GenericMap[T, interface{}]
//	for _, e := range *s {
//		if !m.ContainsKey(e) {
//			m.Put(e, nil)
//			n.Push(e)
//		}
//	}
//	return n
//}

type GenericMap[K MapKey, V any] map[K]V

func (m *GenericMap[K, V]) ContainsKey(k K) bool {
	if *m == nil {
		*m = make(map[K]V)
	}
	if _, ok := (*m)[k]; ok {
		return true
	}
	return false
}

func (m *GenericMap[K, V]) DelKey(k K) {
	if *m == nil {
		*m = make(map[K]V)
	}
	delete(*m, k)
}

func (m *GenericMap[K, V]) Put(k K, v V) {
	if *m == nil {
		*m = make(map[K]V)
	}
	(*m)[k] = v
}

func (m GenericMap[K, V]) Keys() Slice[K] {
	var keys Slice[K]
	for k, _ := range m {
		keys.Push(k)
	}
	return keys
}

func (m GenericMap[K, V]) Values() Slice[V] {
	var values Slice[V]
	for _, v := range m {
		values.Push(v)
	}
	return values
}
