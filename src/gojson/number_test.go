package gojson

import "testing"

func TestNumber(t *testing.T) {
	n := NewNumber(int64(10))
	if value, ok := n.int64Value(); ok != nil {
		t.Error("int64Value() failed")
		if value != int64(10) {
			t.Error("int64Value() failed")
		}
	}
	n = NewNumber(float64(0.0123))
	if value, ok := n.float64Value(); ok != nil {
		t.Error("float64Value() failed")
		if value != float64(0.0123) {
			t.Error("float64Value() failed")
		}
	}
}
