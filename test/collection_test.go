package test

import (
	"fmt"
	"github.com/wangleilei2010/gogo/collection"
	"testing"
)

func convert(p []int) {

}

func TestSliceInit(t *testing.T) {
	var s collection.Slice[int]
	s.Push(1)
	s.Foreach(func(e int) { fmt.Println(e) })
	convert(s)
	s.Append(2, 3, 4, 5)
	s.Foreach(func(e int) { fmt.Println(e) })
	fmt.Println("len =", s.Len())
}

func TestGenericInit(t *testing.T) {
	var m collection.GenericMap[int, int]
	m.Put(1, 2)
	m.Put(3, 4)
	for k, v := range m {
		fmt.Println(k, v)
	}
	m.Keys().Foreach(func(k int) { fmt.Print(k) })
	m.Values().Foreach(func(v int) { fmt.Print(v) })
}
