package gojson

import "testing"

func TestStackValue(t *testing.T) {
	if TYPE_OBJECT != 0 {
		t.Error("TYPE_OBJECT failed")
	}
	if TYPE_OBJECT_KEY != 1 {
		t.Error("TYPE_OBJECT_KEY failed")
	}
	if TYPE_ARRAY != 2 {
		t.Error("TYPE_ARRAY failed")
	}
	if TYPE_SINGLE != 3 {
		t.Error("TYPE_SINGLE failed")
	}

	s1 := NewJsonObjectFromObject(1)
	if s1.Value.(int) != 1 || s1.Kind != TYPE_OBJECT {
		t.Error("NewJsonObjectFromObject() failed")
	}
	s2 := NewJsonObjectFromKey("key")
	if s2.ValueAsKey() != "key" || s2.Value.(string) != "key" || s2.Kind != TYPE_OBJECT_KEY {
		t.Error("NewJsonObjectFromKey() failed")
	}
	arr := []interface{}{1, 2, 3, 4}
	s3 := NewJsonObjectFromSlice(arr)
	if s3.Kind != TYPE_ARRAY {
		t.Error("NewJsonObjectFromSlice() failed")
	}
	s4 := NewJsonObjectFromSingle(123)
	if s4.Kind != TYPE_SINGLE || s4.Value.(int) != 123 {
		t.Error("NewJsonObjectFromSingle() failed")
	}

	m := make(map[string]interface{})
	m["index"] = 1
	m["name"] = "golang"
	s5 := NewJsonObjectFromMap(m)
	if s5.Kind != TYPE_OBJECT {
		t.Error("NewJsonObjectFromMap() failed")
	}
}
