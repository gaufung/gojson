package gojson

import (
	"strings"
	"testing"
)

func TestTokenReader1(t *testing.T) {
	var json = `{"language": "golang", "version":1.9}`
	reader := NewCharReader(strings.NewReader(json))
	r := &TokenReader{reader}
	if !r.isWhiteSpace('\t') || !r.isWhiteSpace('\n') || !r.isWhiteSpace(' ') || !r.isWhiteSpace('\r') {
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
	if r.readNumber() != float64(1.9) {
		t.Error("readNumber() failed")
	}
}

func TestTokeReader2(t *testing.T) {
	var json = `{"isread": false, "company": [1, 2.5, 3], "leader": null}`
	reader := NewCharReader(strings.NewReader(json))
	r := &TokenReader{reader}
	if r.readNextToken() != START_OBJECT {
		t.Error("read START_OBJECT failed")
	}
	if r.readNextToken() != STRING {
		t.Error("read STRING failed")
	}
	if r.readString() != "isread" {
		t.Error(`read "isread" failed`)
	}
	if r.readNextToken() != COLON_SEPERATOR {
		t.Error(`read COLON_SEPERATOR failed`)
	}
	if r.readNextToken() != BOOLEAN {
		t.Error(`read BOOLEAN failed`)
	}
	if r.readBoolean() != false {
		t.Error(`read false failed`)
	}
	if r.readNextToken() != COMA_SEPERATOR {
		t.Error(`read COMA_SEPERATOR failed`)
	}
	if r.readNextToken() != STRING {
		t.Error(`read STRING failed`)
	}
	if r.readString() != "company" {
		t.Error(`read "company" failed`)
	}
	if r.readNextToken() != COLON_SEPERATOR {
		t.Error(`read COLON_SEPERATOR failed`)
	}
	if r.readNextToken() != START_ARRAY {
		t.Error(`read START_ARRAY failed`)
	}
	if r.readNextToken() != NUMBER {
		t.Error(`read NUMBER failed`)
	}
	if r.readNumber() != 1 {
		t.Error("read 1 failed")
	}
	if r.readNextToken() != COMA_SEPERATOR {
		t.Error("read COMA_SEPERATOR failed")
	}
	if r.readNextToken() != NUMBER {
		t.Error(`read NUMBER failed`)
	}
	if r.readNumber() != 2.5 {
		t.Error("read 2.5 failed")
	}
	if r.readNextToken() != COMA_SEPERATOR {
		t.Error("read COMA_SEPERATOR failed")
	}
	if r.readNextToken() != NUMBER {
		t.Error("read NUMBER failed")
	}
	if r.readNumber() != 3 {
		t.Error("read 2 failed")
	}
	if r.readNextToken() != END_ARRAY {
		t.Error("read END_ARRAY failed")
	}
	if r.readNextToken() != COMA_SEPERATOR {
		t.Error(`read COMA_SEPERATOR failed`)
	}
	if r.readNextToken() != STRING {
		t.Error(`read STRING failed`)
	}
	if r.readString() != `leader` {
		t.Error(`read "leader" failed`)
	}
	if r.readNextToken() != COLON_SEPERATOR {
		t.Error(`read COLON_SEPERATOR failed`)
	}
	if r.readNextToken() != NULL {
		t.Error(`read NULL failed`)
	}
	r.readNull()
	if r.readNextToken() != END_OBJECT {
		t.Error("read END_OBJECT failed")
	}
}
