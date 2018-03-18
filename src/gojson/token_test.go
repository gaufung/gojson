package gojson

import "testing"

func TestToken(t *testing.T){
	i := END_DOCUMENT
	if i != 0 {
		t.Error("const error")
	}
	i = START_OBJECT
	if i != 1 {
		t.Error("const error")
	}
	i = END_OBJECT
	if i != 2 {
		t.Error("const error")
	}
}