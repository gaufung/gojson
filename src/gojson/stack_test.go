package gojson

import "testing"

func TestStack(t *testing.T) {
	stack := newStack()
	if stack.Size != 100 || stack.pos != 0 {
		t.Error("size failed")
	}
	if !stack.isEmpty() {
		t.Error("IsEmpty() failed")
	}
	item1 := newStackValueFromKey("key")
	stack.push(item1)
	if stack.isEmpty() || stack.pos != 1 {
		t.Error("IsEmpty() failed")
	}
	if item2, ok := stack.pop(); !ok {
		t.Error("Pop() failed")
	} else {
		if item2.value != "key" {
			t.Error("Pop() failed")
		}
	}

	if _, ok := stack.pop(); ok {
		t.Error("Pop() failed")
	}
	item3 := newStackValueFromObject(42)
	stack.push(item3)
	if _, ok := stack.popKind(TYPE_OBJECT); !ok {
		t.Error("PopKind() failed")
	}
	item4 := newStackValueFromSlice([]interface{}{1, 2, 3})
	stack.push(item4)
	if item, ok := stack.peek(TYPE_ARRAY); !ok {
		t.Error("Peek() failed")
		arr := item.valueAsArray()
		if arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
			t.Error("Peek() failed")
		}
	}
	m := make(map[string]interface{})
	m["language"] = "Golang"
	m["version"] = 1.9
	item5 := newStackValueFromMap(m)
	stack.push(item5)
	if stack.getTopValueType() != TYPE_OBJECT {
		t.Error("GetTopValueType() failed")
	}
}
