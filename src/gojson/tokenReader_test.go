package gojson

import (
	"strings"
	"testing"
)

func TestTokenReader1(t *testing.T) {
	var json = `{"language": "golang", "version":1.9}`
	reader := newCharReader(strings.NewReader(json))
	r := &tokenReader{reader}
	if !r.isWhiteSpace('\t') || !r.isWhiteSpace('\n') || !r.isWhiteSpace(' ') || !r.isWhiteSpace('\r') {
		t.Error("isWhiteSpace() failed")
	}
	if toke, err := r.readNextToken(); err != nil {
		t.Error("readNextToken() failed")
	} else {
		if toke != START_OBJECT {
			t.Error("readNextToken() failed")
		}
	}
	if toke, err := r.readNextToken(); err != nil {
		t.Error("readNextToken() failed")
	} else {
		if toke != STRING {
			t.Error("readNextToken() failed")
		}
	}
	if str, err := r.readString(); err != nil {
		t.Error("readString() failed")
	} else {
		if str != "language" {
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("readToken() failed")
	} else {
		if token != COLON_SEPERATOR {
			t.Error("readToken() failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("readToken() failed")
	} else {
		if token != STRING {
			t.Error("readToken() failed")
		}
	}
	if str, err := r.readString(); err != nil {
		t.Error("readString() failed")
	} else {
		if str != "golang" {
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("readToken() failed")
	} else {
		if token != COMA_SEPERATOR {
			t.Error("readToken() failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("readToken() failed")
	} else {
		if token != STRING {
			t.Error("readToken() failed")
		}
	}
	if str, err := r.readString(); err != nil {
		t.Error("readString() failed")
	} else {
		if str != "version" {
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("readToken() failed")
	} else {
		if token != COLON_SEPERATOR {
			t.Error("readToken() failed--")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("readToken() failed")
	} else {
		if token != NUMBER {
			t.Error("readToken() failed--")
		}
	}
	if number, err := r.readNumber(); err != nil {
		t.Error("read number failed")
	} else {
		if number != float64(1.9) {
			t.Error("read number failed")
		}
	}
}

func TestTokeReader2(t *testing.T) {
	var json = `{"isread": false, "company": [1, 2.5], "leader": null}`
	reader := newCharReader(strings.NewReader(json))
	r := &tokenReader{reader}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != START_OBJECT {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != STRING {
			t.Error("read token failed")
		}
	}
	if str, err := r.readString(); err != nil {
		t.Error("readString() failed")
	} else {
		if str != "isread" {
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != COLON_SEPERATOR {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != BOOLEAN {
			t.Error("read token failed")
		}
	}
	if b, err := r.readBoolean(); err != nil {
		t.Error("read boolean failed")
	} else {
		if b != false {
			t.Error("read boolean failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != COMA_SEPERATOR {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != STRING {
			t.Error("read token failed")
		}
	}
	if str, err := r.readString(); err != nil {
		t.Error("readString() failed")
	} else {
		if str != "company" {
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != COLON_SEPERATOR {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != START_ARRAY {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != NUMBER {
			t.Error("read token failed")
		}
	}
	if num, err := r.readNumber(); err != nil {
		t.Error("read number failed")
	} else {
		if num != float64(1) {
			t.Error("read number failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != COMA_SEPERATOR {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != NUMBER {
			t.Error("read token failed")
		}
	}
	if num, err := r.readNumber(); err != nil {
		t.Error("read number failed")
	} else {
		if num != float64(2.5) {
			t.Error("read number failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != END_ARRAY {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != COMA_SEPERATOR {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != STRING {
			t.Error("read token failed")
		}
	}
	if str, err := r.readString(); err != nil {
		t.Error("readString() failed")
	} else {
		if str != "leader" {
			t.Error("readString() failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != COLON_SEPERATOR {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != NULL {
			t.Error("read token failed")
		}
	}
	if err := r.readNull(); err != nil {
		t.Error("read null failed")
	}
}

func TestTokenReader3(t *testing.T) {
	var json = `{"\\ba": 1.3e2}`
	reader := newCharReader(strings.NewReader(json))
	r := &tokenReader{reader}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != START_OBJECT {
			t.Error("Read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read Token failed")
	} else {
		if token != STRING {
			t.Error("read token failed")
		}
	}
	if str, err := r.readString(); err != nil {
		t.Error("read string failed")
	} else {
		if str != "\\ba" {
			t.Error("read string failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read token failed")
	} else {
		if token != COLON_SEPERATOR {
			t.Error("read token failed")
		}
	}
	if token, err := r.readNextToken(); err != nil {
		t.Error("read token failed")
	} else {
		if token != NUMBER {
			t.Error("read token failed")
		}
	}
	if num, err := r.readNumber(); err != nil {
		t.Error("read number failed")
	} else {
		if num != float64(130) {
			t.Errorf("%f", num)
			t.Error("number doesn't equal")
		}
	}
}


func createTokenReaderFromString(str string) *tokenReader {
	reader := newCharReader(strings.NewReader(str))
	return &tokenReader{reader}
}

func TestTokeReaderNumber(t *testing.T) {
	str := "123 123.2 -123 -12.3 12e2 12e-2 -1.2e2"
	r := createTokenReaderFromString(str)
	if val, err := r.readNumber(); err != nil {
		t.Error(err.Error())
	}else{
		if val != 123 {
			t.Error("123 doesn't equal")
		}
	}
	r.readNextToken()
	if val, err := r.readNumber(); err != nil {
		t.Error(err.Error())
	}else{
		if val != 123.2 {
			t.Error("123.2 doesn't equal")
		}
	}
	r.readNextToken()
	if val, err := r.readNumber(); err != nil {
		t.Error(err.Error())
	}else{
		if val != -123 {
			t.Error("-123 doesn't equal")
		}
	}
	r.readNextToken()
	if val, err := r.readNumber(); err != nil{
		t.Error(err.Error())
	}else{
		if val != -12.3{
			t.Error("-12.3 doesn't equal")
		}
	}
	r.readNextToken()
	if val, err := r.readNumber(); err != nil {
		t.Error(err.Error())
	}else{
		if val != 1200 {
			t.Error("12e2 doesn't equal")
		}
	}
	r.readNextToken()
	if val, err := r.readNumber(); err!=nil{
		t.Error(err.Error())
	}else{
		if val != 0.12 {
			t.Error("12e-2 doesn't equal")
		}
	}
	r.readNextToken()
	if val, err := r.readNumber(); err != nil {
		t.Error(err.Error())
	}else{
		if val != -120 {
			t.Error("1.2e2 doesn't equal")
		}
	}
}

func TestTokenReaderBoolean(t *testing.T){
	str := "true false flase"
	r := createTokenReaderFromString(str)
	if b, err := r.readBoolean(); err != nil{
		t.Error(err.Error())
	}else{
		if b!=true{
			t.Error("true failed")
		}
	}
	r.readNextToken()
	if b, err := r.readBoolean(); err != nil {
		t.Error(err.Error())
	}else{
		if b!=false{
			t.Error("false failed")
		}
	}
	r.readNextToken()
	if _, err := r.readBoolean(); err == nil {
		t.Error(err.Error())
	}
}

func TestTokenReaderNull(t *testing.T) {
	str := "null"
	r := createTokenReaderFromString(str)
	if err := r.readNull(); err != nil{
		t.Error(err.Error())
	}
	str = "nULL"
	r = createTokenReaderFromString(str)
	if err := r.readNull(); err == nil {
		t.Error(err.Error())
	}
}

func TestTokenReaderString(t *testing.T){
	str := `"golang" "\\\"\/\\b\\f\\n\\r\\t\\u09f1"`
	r := createTokenReaderFromString(str)
	if str, err := r.readString(); err !=nil {
		t.Error(err.Error())
	}else{
		if str != "golang" {
			t.Error(`read "golang" failed`)
		}
	}
	r.readNextToken()
	if val, err := r.readString(); err != nil {
		t.Error(err.Error())
	}else{
		if val != `\"/\b\f\n\r\t\u09f1`{
			t.Error("read string failed")
		}
	}
}