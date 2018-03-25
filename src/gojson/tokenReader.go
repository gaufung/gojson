package gojson

import (
	"math"
	"strings"
	"io"
	"errors"
	"strconv"
)

const (
	READ_NUMBER_INT_PART = iota
	READ_NUMBER_FRA_PART
	READ_NUMBER_EXP_PART
	READ_NUMBER_END
)

type TokenReader struct {
	reader *CharReader
}

func NewTokenReaderFromString(s string) *TokenReader {
	return newTokenReader(strings.NewReader(s))
}

func newTokenReader(r io.Reader) *TokenReader{
	return &TokenReader{newCharReader(r)} 
}

func (t *TokenReader) position() int {
	return t.reader.pos
}

func (t *TokenReader) errors(info string) error{
	return errors.New(info+" at " + strconv.Itoa(t.position()))
}

// white space to ignore
func (t *TokenReader) isWhiteSpace(ch rune) bool {
	return ch == '\n' || ch == '\t' || ch == ' ' || ch == '\r'
}

// read next token
func (t *TokenReader) readNextToken() (Token,error) {
	ch := '?'
	for {
		if !t.reader.hasMore() {
			return END_DOCUMENT, nil
		}
		ch = t.reader.peek() //peek it
		if !t.isWhiteSpace(ch) {
			break
		}
		t.reader.next()
	}
	switch ch {
	case '{':
		t.reader.next() //skip
		return START_OBJECT, nil
	case '}':
		t.reader.next() //skip
		return END_OBJECT, nil
	case '[':
		t.reader.next() //skip
		return START_ARRAY, nil
	case ']':
		t.reader.next() //skip
		return END_ARRAY, nil
	case ':':
		t.reader.next() //skip
		return COLON_SEPERATOR, nil
	case ',':
		t.reader.next() //skip
		return COMA_SEPERATOR, nil
	case '"':
		return STRING, nil
	case 'n':
		return NULL, nil
	case 't':
		return BOOLEAN, nil
	case 'f':
		return BOOLEAN, nil
	case '-':
		return NUMBER, nil
	}
	if ch >= '0' && ch <= '9' {
		return NUMBER, nil
	}
	//panic(NewJsonParserError("Parse error when try to guess next token.", t.reader.pos))
	return -1, t.errors("Unexpected token ")
}

// read string
func (t *TokenReader) readString() (string,error) {
	result := make([]rune, 0)
	ch := t.reader.next()
	if ch != '"' {
		return "", t.errors("Unexpected string")
	}
	for {
		ch = t.reader.next()
		if ch == '\\' {
			ech := t.reader.next()
			switch ech {
			case '"':
				result = append(result, ech)
			case '\\':
				result = append(result, ech)
			case '/':
				result = append(result, ech)
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			case 'u':
				u := 0
				for i := 0; i < 4; i++ {
					uch := t.reader.next()
					if uch >= '0' && uch <= '9' {
						u = (u << 4) + (int(uch) - int('0'))
					} else if uch >= 'a' && uch <= 'f' {
						u = (u << 4) + (int(uch) - int('0')) + 10
					} else if uch >= 'A' && uch <= 'F' {
						u = (u << 4) + (int(uch) - int('A')) + 10
					} else {
						//panic(NewJsonParserError("Unexpected char", t.reader.pos))
						return "", t.errors("Read string: unexpected char ")
					}
				}
				result = append(result, rune(u))
			default:
				return "", t.errors("Read string: unexpected char")
			}
		} else if ch == '"' {
			break
		} else if ch == '\r' || ch == '\n' {
			return "", t.errors("Read string: unexpected char")
		} else {
			result = append(result, ch)
		}
	}
	return string(result), nil
}

// read boolean
func (t *TokenReader) readBoolean() (bool,error) {
	ch := t.reader.next()
	expect := ""
	if ch == 't' {
		expect = "rue"
	} else if ch == 'f' {
		expect = "alse"
	} else {
		//panic(NewJsonParserError("Unexpected char", t.reader.pos))
		return false, t.errors("Read boolean: unexpected char ")
	}
	for _, c := range []rune(expect) {
		theChar := t.reader.next()
		if theChar != c {
			return false, t.errors("Read boolean: unexpected char ")
		}
	}
	return ch == 't', nil
}

func (t *TokenReader) readNull() error {
	expect := "null"
	for _, c := range []rune(expect) {
		theChar := t.reader.next()
		if theChar != c {
			return t.errors("Read null: unexpected char")
		}
	}
	return nil
}

// read number
func (t *TokenReader) readNumber() (float64, error) {
	intPart, fraPart, expPart := make([]rune, 0), make([]rune, 0), make([]rune, 0)
	hasFraPart, hasExpPart := false, false
	ch := t.reader.peek()
	minusSign := ch == '-'
	expMinusSign := false
	if minusSign {
		t.reader.next()
	}
	status := READ_NUMBER_INT_PART
	for {
		if t.reader.hasMore() {
			ch = t.reader.peek()
		} else {
			status = READ_NUMBER_END
		}
		switch status {
		case READ_NUMBER_INT_PART:
			if ch >= '0' && ch <= '9' {
				intPart = append(intPart, t.reader.next())
			} else if ch == '.' {
				if len(intPart) == 0 {
					return 0.0, t.errors("Read float64: unexpect char ")
				}
				t.reader.next()
				hasFraPart = true
				status = READ_NUMBER_FRA_PART
			} else if ch == 'e' || ch == 'E' {
				t.reader.next()
				hasExpPart = true
				signChar := t.reader.peek()
				if signChar == '-' || signChar == '+' {
					expMinusSign = signChar == '-'
					t.reader.next()
				}
				status = READ_NUMBER_EXP_PART
			} else {
				if len(intPart) == 0 {
					return 0.0, t.errors("Read float64: unexpect char ")
				}
				status = READ_NUMBER_END
			}
			continue
		case READ_NUMBER_FRA_PART:
			if ch >= '0' && ch <= '9' {
				fraPart = append(fraPart, t.reader.next())
			} else if ch == 'e' || ch == 'E' {
				t.reader.next()
				hasExpPart = true
				signChar := t.reader.peek()
				if signChar == '-' || signChar == '+' {
					expMinusSign = signChar == '-'
					t.reader.next()
				}
				status = READ_NUMBER_EXP_PART
			} else {
				if len(fraPart) == 0 {
					//panic(NewJsonParserError("Unexpected char", t.reader.pos))
					return 0.0, t.errors("Read float64, unexpected char ")
				}
				status = READ_NUMBER_END
			}
			continue
		case READ_NUMBER_EXP_PART:
			if ch >= '0' && ch <= '9' {
				expPart = append(expPart, t.reader.next())
			} else {
				if len(expPart) == 0 {
					return 0.0, t.errors("Read float64, unexpected char ")
				}
				status = READ_NUMBER_END
			}
			continue
		case READ_NUMBER_END:
			if len(intPart) == 0 {
				return 0.0, t.errors("Read float64, unexpected char ")
			}
			lint := int64(0)
			if minusSign {
				//lint = -1 * string2long(intPart, readed)
				if val, err:=string2long(intPart, t.position()); err==nil{
					lint = -1 * val
				}else{
					return 0.0, err
				}
			} else {
				//lint = string2long(intPart, readed)
				if val, err := string2long(intPart, t.position()); err==nil{
					lint = val
				}else{
					return 0.0, err
				}
			}
			if hasExpPart && len(expPart) == 0 {
				return float64(lint), nil
			}
			if hasFraPart && len(fraPart) == 0 {
				return 0.0, t.errors("Read float64, unexpected char ")
			}
			dFraPart := float64(0.0)
			if hasFraPart {
				if minusSign {
					//dFraPart = -float64(1.0) * string2Fraction(fraPart, t.reader.pos)
					if val, err := string2Fraction(fraPart, t.position()); err == nil{
						dFraPart = -float64(1.0) * val
					}else{
						return 0.0, err
					}
				} else {
					//dFraPart = string2Fraction(fraPart, t.reader.pos)
					if val, err := string2Fraction(fraPart, t.position()); err==nil{
						dFraPart = val
					}else{
						return 0.0, err
					}
				}
			}
			number := float64(0.0)
			if hasExpPart {
				index := int64(0.0)
				if expMinusSign {
					//index = -1.0 * string2long(expPart, readed)
					if val, err := string2long(expPart, t.position()); err==nil{
						index = -1.0 * val
					}else{
						return 0.0, err
					}
				} else {
					//index = string2long(expPart, readed)
					if val, err := string2long(expPart, t.position()); err!=nil{
						index = val
					}else{
						return 0.0, err
					}
				}
				number = (float64(lint) + dFraPart) * math.Pow(10.0, float64(index))
			} else {
				number = float64(lint) + dFraPart
			}
			if number > maxSafeFloat64 {
				return 0.0, t.errors("Read float64, unexpected char ")
			}
			return number, nil
		}
	}
}

//back token
func (t *TokenReader) BackToken() {
	t.reader.backward()
}

//determine whether is empty or not
func (t *TokenReader) IsEmpty() bool {
	if !t.reader.hasMore() {
		return true
	} else {
		return false
	}
}

var maxSafeInt64 = int64(9007199254740991)

func string2long(str []rune, readed int) (int64, error) {
	if len(str) > 16 {
		return 0.0, errors.New("Read long: too many digits at "+strconv.Itoa(readed))
	}
	result := int64(0)
	for _, digit := range str {
		result = result*10 + int64(digit-'0')
		if result > maxSafeInt64 {
			return 0.0, errors.New("Read long: too big at "+ strconv.Itoa(readed))
		}
	}
	return result, nil
}

var maxSafeFloat64 = 1.7976931348623157e+308

func string2Fraction(str []rune, readed int) (float64, error) {
	if len(str) > 16 {
		return 0.0, errors.New("Read fraction: too many digits at " +strconv.Itoa(readed))
	}
	result := float64(0.0)
	for idx, digit := range str {
		n := int(digit - '0')
		if n == 0 {
			result = result + 0
		} else {
			result = result + float64(n)/math.Pow(10, float64(idx+1))
		}
	}
	return result, nil
}
