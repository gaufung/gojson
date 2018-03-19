package gojson

const (
	TYPE_OBJECT = iota
	TYPE_OBJECT_KEY
	TYPE_ARRAY
	TYPE_SINGLE
)

type StackValue struct {
	Kind  int
	value interface{}
}

func NewJsonObjectFromObject(obj interface{}) *StackValue {
	return &StackValue{Kind: TYPE_OBJECT, value: obj}
}

func NewJsonObjectFromMap(m map[string]interface{}) *StackValue {
	return &StackValue{Kind: TYPE_OBJECT, value: m}
}

func NewJsonObjectFromKey(key string) *StackValue {
	return &StackValue{Kind: TYPE_OBJECT_KEY, value: key}
}

func NewJsonObjectFromSlice(arr []interface{}) *StackValue {
	return &StackValue{Kind: TYPE_ARRAY, value: arr}
}

func NewJsonObjectFromSingle(obj interface{}) *StackValue {
	return &StackValue{Kind: TYPE_SINGLE, value: obj}
}

func (s *StackValue) ValueAsKey() string {
	return s.value.(string)
}

func (s *StackValue) ValueAsObject() map[string]interface{} {
	return s.value.(map[string]interface{})
}

func (s *StackValue) ValueAsArray() []interface{} {
	return s.value.([]interface{})
}
