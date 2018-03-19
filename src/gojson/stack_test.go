package gojson

import "testing"

func TestStack(t *testing.T) {
	stack := NewStack()
	if stack.Size != 100 || stack.pos != 0 {
		t.Error("size failed")
	}
	if !stack.IsEmpty() {
		t.Error("IsEmpty() failed")
	}
	item1 := NewJsonObjectFromKey("key")
	stack.Push(item1)
	if stack.IsEmpty() || stack.pos != 1 {
		t.Error("IsEmpty() failed")
	}
	if item2, err := stack.Pop(); err != nil {
		t.Error("Pop() failed")
	} else {
		if item2.value != "key" {
			t.Error("Pop() failed")
		}
	}

	if _, err := stack.Pop(); err == nil {
		t.Error("Pop() failed")
	}

	item3 := NewJsonObjectFromObject(42)
	stack.Push(item3)
	if _, err := stack.PopKind(TYPE_OBJECT_KEY); err == nil {
		t.Error("PopKind() failed")
	}
	item4 := NewJsonObjectFromSlice([]interface{}{1, 2, 3})
	stack.Push(item4)
	if item, err := stack.Peek(TYPE_ARRAY); err != nil {
		t.Error("Peek() failed")
		arr := item.ValueAsArray()
		if arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
			t.Error("Peek() failed")
		}
	}
	m := make(map[string]interface{})
	m["language"] = "Golang"
	m["version"] = 1.9
	item5 := NewJsonObjectFromMap(m)
	stack.Push(item5)
	if stack.GetTopValueType() != TYPE_OBJECT {
		t.Error("GetTopValueType() failed")
	}
}
