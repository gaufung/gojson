package gojson

import (
	"math"
	"strings"
	"io"
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

func newTokenReaderFromString(s string) *TokenReader {
	return newTokenReader(strings.NewReader(s))
}

func newTokenReader(r io.Reader) *TokenReader{
	return &TokenReader{newCharReader(r)} 
}

func (t *TokenReader) position() int {
	return t.reader.pos
}

// white space to ignore
func (t *TokenReader) isWhiteSpace(ch rune) bool {
	return ch == '\n' || ch == '\t' || ch == ' ' || ch == '\r'
}

// read next token
func (t *TokenReader) readNextToken() Token {
	ch := '?'
	for {
		if !t.reader.hasMore() {
			return END_DOCUMENT
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
		return START_OBJECT
	case '}':
		t.reader.next() //skip
		return END_OBJECT
	case '[':
		t.reader.next() //skip
		return START_ARRAY
	case ']':
		t.reader.next() //skip
		return END_ARRAY
	case ':':
		t.reader.next() //skip
		return COLON_SEPERATOR
	case ',':
		t.reader.next() //skip
		return COMA_SEPERATOR
	case '"':
		return STRING
	case 'n':
		return NULL
	case 't':
		return BOOLEAN
	case 'f':
		return BOOLEAN
	case '-':
		return NUMBER
	}
	if ch >= '0' && ch <= '9' {
		return NUMBER
	}
	panic(NewJsonParserError("Parse error when try to guess next token.", t.reader.pos))
}

// read string
func (t *TokenReader) readString() string {
	result := make([]rune, 0)
	ch := t.reader.next()
	if ch != '"' {
		panic(NewJsonParserError("Expected \" but actual is: ", t.reader.pos))
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
						panic(NewJsonParserError("Unexpected char", t.reader.pos))
					}
				}
				result = append(result, rune(u))
			default:
				panic(NewJsonParserError("Unexpected char", t.reader.pos))
			}
		} else if ch == '"' {
			break
		} else if ch == '\r' || ch == '\n' {
			panic(NewJsonParserError("Unexpected char", t.reader.pos))
		} else {
			result = append(result, ch)
		}
	}
	return string(result)
}

// read boolean
func (t *TokenReader) readBoolean() bool {
	ch := t.reader.next()
	expect := ""
	if ch == 't' {
		expect = "rue"
	} else if ch == 'f' {
		expect = "alse"
	} else {
		panic(NewJsonParserError("Unexpected char", t.reader.pos))
	}
	for _, c := range []rune(expect) {
		theChar := t.reader.next()
		if theChar != c {
			panic(NewJsonParserError("Unexpected char", t.reader.pos))
		}
	}
	return ch == 't'
}

func (t *TokenReader) readNull() {
	expect := "null"
	for _, c := range []rune(expect) {
		theChar := t.reader.next()
		if theChar != c {
			panic(NewJsonParserError("Unexpected char", t.reader.pos))
		}
	}
}

// read number
func (t *TokenReader) readNumber() float64 {
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
					panic(NewJsonParserError("Unexpected char", t.reader.pos))
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
					panic(NewJsonParserError("Unexpected char", t.reader.pos))
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
					panic(NewJsonParserError("Unexpected char", t.reader.pos))
				}
				status = READ_NUMBER_END
			}
			continue
		case READ_NUMBER_EXP_PART:
			if ch >= '0' && ch <= '9' {
				expPart = append(expPart, t.reader.next())
			} else {
				if len(expPart) == 0 {
					panic(NewJsonParserError("Unexpected char", t.reader.pos))
				}
				status = READ_NUMBER_END
			}
			continue
		case READ_NUMBER_END:
			readed := t.reader.pos
			if len(intPart) == 0 {
				panic(NewJsonParserError("Unexpected char", t.reader.pos))
			}
			lint := int64(0)
			if minusSign {
				lint = -1 * string2long(intPart, readed)
			} else {
				lint = string2long(intPart, readed)
			}
			if hasExpPart && len(expPart) == 0 {
				return float64(lint)
			}
			if hasFraPart && len(fraPart) == 0 {
				panic(NewJsonParserError("Unexpected char", t.reader.pos))
			}
			dFraPart := float64(0.0)
			if hasFraPart {
				if minusSign {
					dFraPart = -float64(1.0) * string2Fraction(fraPart, t.reader.pos)
				} else {
					dFraPart = string2Fraction(fraPart, t.reader.pos)
				}
			}
			number := float64(0.0)
			if hasExpPart {
				index := int64(0.0)
				if expMinusSign {
					index = -1.0 * string2long(expPart, readed)
				} else {
					index = string2long(expPart, readed)
				}
				number = (float64(lint) + dFraPart) * math.Pow(10.0, float64(index))
			} else {
				number = float64(lint) + dFraPart
			}
			if number > maxSafeFloat64 {
				panic(NewJsonParserError("Exceeded maximum value", t.reader.pos))
			}
			return number
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

func string2long(str []rune, readed int) int64 {
	if len(str) > 16 {
		panic(NewJsonParserError("Number string is too long", readed))
	}
	result := int64(0)
	for _, digit := range str {
		result = result*10 + int64(digit-'0')
		if result > maxSafeInt64 {
			panic(NewJsonParserError("Exceeded maximum value", readed))
		}
	}
	return result
}

var maxSafeFloat64 = 1.7976931348623157e+308

func string2Fraction(str []rune, readed int) float64 {
	if len(str) > 16 {
		panic(NewJsonParserError("Number string is too long", readed))
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
	return result
}
