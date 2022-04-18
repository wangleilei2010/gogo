package collection

type Slice[T any] []T

func (s *Slice[T]) Push(t T) {
	ns := make([]T, 0)
	for _, e := range *s {
		ns = append(ns, e)
	}
	ns = append(ns, t)
	*s = ns
}

func (s Slice[T]) Len() int {
	l := 0
	for i, _ := range s {
		l = i
	}
	return l + 1
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

type GenericMap[K int | string, V any] map[K]V
