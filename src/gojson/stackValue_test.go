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

	s1 := newStackValueFromObject(1)
	if s1.value.(int) != 1 || s1.kind != TYPE_OBJECT {
		t.Error("NewJsonObjectFromObject() failed")
	}
	s2 := newStackValueFromKey("key")
	if s2.valueAsKey() != "key" || s2.value.(string) != "key" || s2.kind != TYPE_OBJECT_KEY {
		t.Error("NewJsonObjectFromKey() failed")
	}
	arr := []interface{}{1, 2, 3, 4}
	s3 := newStackValueFromSlice(arr)
	if s3.kind != TYPE_ARRAY {
		t.Error("NewJsonObjectFromSlice() failed")
	}
	m := make(map[string]interface{})
	m["index"] = 1
	m["name"] = "golang"
	s5 := newStackValueFromMap(m)
	if s5.kind != TYPE_OBJECT {
		t.Error("NewJsonObjectFromMap() failed")
	}
}
