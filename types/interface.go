package types

type IJsonValue interface {
	Get(key string) IJsonValue
	DeepGet(key string) IJsonValue
	ContainsKey(key string) bool

	GetAt(index int) IJsonValue

	Contains(v interface{}) bool

	Value() interface{}

	IsMap() bool
	IsSliceOrArray() bool
}
