package gojson

import (
	"testing"
	"strings"
)

var json = `{"language": "golang", "version":1.9}`

func TestTokenReader(t *testing.T) {
	reader := NewCharReader(strings.NewReader(json))
	r := &TokenReader{reader}
	if !r.isWhiteSpace('\t') || !r.isWhiteSpace('\n') || !r.isWhiteSpace(' ') || !r.isWhiteSpace('\r'){
		t.Error("isWhiteSpace() failed")
	}
	if r.readNextToken() != START_OBJECT {
		t.Error("readNextToken() failed")
	}
	if r.readNextToken() != STRING {
		t.Error("readNextToken() failed")
	}
	if r.readString() != "language" {
		t.Error("readString() failed")
	}
	if r.readNextToken() != COLON_SEPERATOR {
		t.Error("readNextToken() failed")
	}
	if r.readNextToken() != STRING {
		t.Error("readNextToken() failed")
	}
	if r.readString() != "golang" {
		t.Error("readNextToken() failed")
	}
	if r.readNextToken() != COMA_SEPERATOR {
		t.Error("readNextToken() failed")
	}
	if r.readNextToken() != STRING {
		t.Error("readNextToken() failed")
	}
	if r.readString() != "version" {
		t.Error("readNextToken() failed")
	}
	if r.readNextToken() != COLON_SEPERATOR {
		t.Error("readNextToken() failed")
	}
	if r.readNextToken() != NUMBER {
		t.Error("readNextToken() failed")
	}
	if val, err := r.readNumber().float64Value(); err!= nil{
		t.Error("readNumber() failed")
	}else{
		if val != 1.9 {
			t.Error("readNumber() failed")
		}
	}
}
