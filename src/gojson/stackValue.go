package gojson

const (
	TYPE_OBJECT = iota
	TYPE_OBJECT_KEY
	TYPE_ARRAY
	TYPE_SINGLE
)

type StackValue struct {
	Kind  int
	Value interface{}
}

func NewJsonObjectFromObject(obj interface{}) *StackValue {
	return &StackValue{Kind: TYPE_OBJECT, Value: obj}
}

func NewJsonObjectFromMap(m map[string]interface{}) *StackValue {
	return &StackValue{Kind: TYPE_OBJECT, Value: m}
}

func NewJsonObjectFromKey(key string) *StackValue {
	return &StackValue{Kind: TYPE_OBJECT_KEY, Value: key}
}

func NewJsonObjectFromSlice(arr []interface{}) *StackValue {
	return &StackValue{Kind: TYPE_ARRAY, Value: arr}
}

func NewJsonObjectFromSingle(obj interface{}) *StackValue {
	return &StackValue{Kind: TYPE_SINGLE, Value: obj}
}

func (s *StackValue) ValueAsKey() string {
	return s.Value.(string)
}

func (s *StackValue) ValueAsObject() map[string]interface{} {
	return s.Value.(map[string]interface{})
}

func (s *StackValue) ValueAsArray() []interface{} {
	return s.Value.([]interface{})
}
