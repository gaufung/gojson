package gojson

const (
	TYPE_OBJECT = iota
	TYPE_OBJECT_KEY
	TYPE_ARRAY
)

type stackValue struct {
	kind  int
	value interface{}
}

func newStackValueFromObject(obj interface{}) *stackValue {
	return &stackValue{kind: TYPE_OBJECT, value: obj}
}

func newStackValueFromMap(m map[string]interface{}) *stackValue {
	return &stackValue{kind: TYPE_OBJECT, value: m}
}

func newStackValueFromKey(key string) *stackValue {
	return &stackValue{kind: TYPE_OBJECT_KEY, value: key}
}

func newStackValueFromSlice(arr []interface{}) *stackValue {
	return &stackValue{kind: TYPE_ARRAY, value: arr}
}

func (s *stackValue) valueAsKey() string {
	return s.value.(string)
}

func (s *stackValue) valueAsObject() map[string]interface{} {
	return s.value.(map[string]interface{})
}

func (s *stackValue) valueAsArray() []interface{} {
	return s.value.([]interface{})
}
