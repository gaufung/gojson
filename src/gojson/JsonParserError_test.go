package gojson

import "testing"

func TestJsonParseError(t *testing.T){
	err := NewJsonParserError("json parser", 1)
	if err.Error() != "json parser"{
		t.Error("Error() failed")
	}
	if err.GetErrorIndex() != 1{
		t.Error("GetErrorIndex() failed")
	}
}