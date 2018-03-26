package gojson

import (
	"strings"
	"testing"
)

func TestTokenReader1(t *testing.T) {
	var json = `{"language": "golang", "version":1.9}`
	reader := newCharReader(strings.NewReader(json))
	r := &TokenReader{reader}
	if !r.isWhiteSpace('\t') || !r.isWhiteSpace('\n') || !r.isWhiteSpace(' ') || !r.isWhiteSpace('\r') {
		t.Error("isWhiteSpace() failed")
	}
	if toke, err := r.readNextToken(); err!=nil{
		t.Error("readNextToken() failed")
	}else{
		if toke!=START_OBJECT{
			t.Error("readNextToken() failed")
		}
	}
	if toke, err := r.readNextToken(); err!=nil{
		t.Error("readNextToken() failed")
	}else{
		if toke != STRING{
			t.Error("readNextToken() failed")
		}
	}
	if str, err := r.readString(); err!=nil{
		t.Error("readString() failed")
	}else{
		if str!= "language"{
			t.Error("readString() failed")
		}
	}
	if token, err:= r.readNextToken(); err!=nil{
		t.Error("readToken() failed")
	}else{
		if token != COLON_SEPERATOR{
			t.Error("readToken() failed")
		}
	}
	if token, err:= r.readNextToken(); err!=nil{
		t.Error("readToken() failed")
	}else{
		if token != STRING{
			t.Error("readToken() failed")
		}
	}
	if str, err := r.readString(); err!=nil{
		t.Error("readString() failed")
	}else{
		if str!= "golang"{
			t.Error("readString() failed")
		}
	}
	if token, err:= r.readNextToken(); err!=nil{
		t.Error("readToken() failed")
	}else{
		if token != COMA_SEPERATOR{
			t.Error("readToken() failed")
		}
	}
	if token, err:= r.readNextToken(); err!=nil{
		t.Error("readToken() failed")
	}else{
		if token != STRING{
			t.Error("readToken() failed")
		}
	}
	if str, err := r.readString(); err!=nil{
		t.Error("readString() failed")
	}else{
		if str!= "version"{
			t.Error("readString() failed")
		}
	}
	if token, err:= r.readNextToken(); err!=nil{
		t.Error("readToken() failed")
	}else{
		if token != COLON_SEPERATOR{
			t.Error("readToken() failed--")
		}
	}
	if token, err:= r.readNextToken(); err!=nil{
		t.Error("readToken() failed")
	}else{
		if token != NUMBER{
			t.Error("readToken() failed--")
		}
	}
	if number, err := r.readNumber(); err!=nil{
		t.Error("read number failed")
	}else{
		if number != float64(1.9) {
			t.Error("read number failed")
		}
	}
}

func TestTokeReader2(t *testing.T) {
	var json = `{"isread": false, "company": [1, 2.5], "leader": null}`
	reader := newCharReader(strings.NewReader(json))
	r := &TokenReader{reader}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != START_OBJECT{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != STRING{
			t.Error("read token failed")
		}
	}
	if str, err := r.readString(); err!=nil{
		t.Error("readString() failed")
	}else{
		if str!= "isread"{
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != COLON_SEPERATOR{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != BOOLEAN{
			t.Error("read token failed")
		}
	}
	if b, err := r.readBoolean(); err!=nil{
		t.Error("read boolean failed")
	}else{
		if b!=false{
			t.Error("read boolean failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != COMA_SEPERATOR{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != STRING{
			t.Error("read token failed")
		}
	}
	if str, err := r.readString(); err!=nil{
		t.Error("readString() failed")
	}else{
		if str!= "company"{
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != COLON_SEPERATOR{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != START_ARRAY{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != NUMBER{
			t.Error("read token failed")
		}
	}
	if num, err := r.readNumber(); err!=nil{
		t.Error("read number failed")
	}else{
		if num != float64(1){
			t.Error("read number failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != COMA_SEPERATOR{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != NUMBER{
			t.Error("read token failed")
		}
	}
	if num, err := r.readNumber(); err!=nil{
		t.Error("read number failed")
	}else{
		if num != float64(2.5){
			t.Error("read number failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != END_ARRAY{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != COMA_SEPERATOR{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != STRING{
			t.Error("read token failed")
		}
	}
	if str, err := r.readString(); err!=nil{
		t.Error("readString() failed")
	}else{
		if str!= "leader"{
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != COLON_SEPERATOR{
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken();err!=nil{
		t.Error("read Token failed")
	}else{
		if token != NULL{
			t.Error("read token failed")
		}
	}
	if err := r.readNull(); err !=nil{
		t.Error("read null failed")
	}
}

func TestTokenReader3(t *testing.T){
	var json = `{"\\ba": 1.3e2}`
	reader := newCharReader(strings.NewReader(json))
	r := &TokenReader{reader}
	if token, err := r.readNextToken(); err!=nil{
		t.Error("read Token failed")
	}else{
		if token != START_OBJECT{
			t.Error("Read token failed")
		}
	}
	if token, err := r.readNextToken(); err!=nil{
		t.Error("read Token failed")
	}else{
		if token != STRING {
			t.Error("read token failed")
		}
	}
	if str, err := r.readString(); err!=nil {
		t.Error("read string failed")
	}else{
		if str != "\\ba"{
			t.Error("read string failed")
		}
	}
	if token, err:= r.readNextToken(); err!=nil{
		t.Error("read token failed")
	}else{
		if token != COLON_SEPERATOR {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err!=nil{
		t.Error("read token failed")
	}else{
		if token != NUMBER {
			t.Error("read token failed")
		}
	}
	if num, err := r.readNumber(); err!=nil {
		t.Error("read number failed")
	}else{
		if num != float64(130) {
			t.Errorf("%f", num)
			t.Error("number doesn't equal")
		}
	}
}