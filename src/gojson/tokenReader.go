package gojson

import "math"

const (
	READ_NUMBER_INT_PART = iota
	READ_NUMBER_FRA_PART
	READ_NUMBER_EXP_PART
	READ_NUMBER_END
)

type TokenReader struct {
	reader *CharReader
}

// is white space
func (t *TokenReader) isWhiteSpace(ch rune) bool {
	return ch == '\n' || ch == '\t' || ch == ' ' || ch == '\r'
}

func (t *TokenReader) readNextToken() Token {
	ch := '?'
	for {
		if !t.reader.HasMore() {
			return END_DOCUMENT
		}
		ch = t.reader.Peek() //peek it
		if !t.isWhiteSpace(ch) {
			break
		}
		t.reader.Next()
	}
	switch ch {
	case '{':
		t.reader.Next() //skip
		return START_OBJECT
	case '}':
		t.reader.Next() //skip
		return END_OBJECT
	case '[':
		t.reader.Next() //skip
		return START_ARRAY
	case ']':
		t.reader.Next() //skip
		return END_ARRAY
	case ':':
		t.reader.Next() //skip
		return COLON_SEPERATOR
	case ',':
		t.reader.Next() //skip
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
	panic(NewJsonParserError("Parse error when try to guess next token.", t.reader.readed))
}

func (t *TokenReader) readString() string {
	result := make([]rune, 0)
	ch := t.reader.Next()
	if ch != '"' {
		panic(NewJsonParserError("Expected \" but actual is: ", t.reader.readed))
	}
	for {
		ch = t.reader.Next()
		if ch == '\\' {
			ech := t.reader.Next()
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
					uch := t.reader.Next()
					if uch >= '0' && uch <= '9' {
						u = (u << 4) + (int(uch) - int('0'))
					} else if uch >= 'a' && uch <= 'f' {
						u = (u << 4) + (int(uch) - int('0')) + 10
					} else if uch >= 'A' && uch <= 'F' {
						u = (u << 4) + (int(uch) - int('A')) + 10
					} else {
						panic(NewJsonParserError("Unexpected char", t.reader.readed))
					}
				}
				result = append(result, rune(u))
			default:
				panic(NewJsonParserError("Unexpected char", t.reader.readed))
			}
		} else if ch == '"' {
			break
		} else if ch == '\r' || ch == '\n' {
			panic(NewJsonParserError("Unexpected char", t.reader.readed))
		} else {
			result = append(result, ch)
		}
	}
	return string(result)
}

func (t *TokenReader) readBoolean() bool {
	ch := t.reader.Next()
	expect := ""
	if ch == 't' {
		expect = "rue"
	} else if ch == 'f' {
		expect = "alse"
	} else {
		panic(NewJsonParserError("Unexpected char", t.reader.readed))
	}
	for _, c := range []rune(expect) {
		theChar := t.reader.Next()
		if theChar != c {
			panic(NewJsonParserError("Unexpected char", t.reader.readed))
		}
	}
	return ch == 't'
}

func (t *TokenReader) readNull() {
	expect := "null"
	for _, c := range []rune(expect) {
		theChar := t.reader.Next()
		if theChar != c {
			panic(NewJsonParserError("Unexpected char", t.reader.readed))
		}
	}
}

func (t *TokenReader) readNumber() *Number {
	intPart, fraPart, expPart := make([]rune, 0), make([]rune, 0), make([]rune, 0)
	hasFraPart, hasExpPart := false, false
	ch := t.reader.Peek()
	minusSign := ch == '-'
	expMinusSign := false
	if minusSign {
		t.reader.Next()
	}
	status := READ_NUMBER_INT_PART
	for {
		if t.reader.HasMore() {
			ch = t.reader.Peek()
		} else {
			status = READ_NUMBER_END
		}
		switch status {
		case READ_NUMBER_INT_PART:
			if ch >= '0' && ch <= '9' {
				intPart = append(intPart, t.reader.Next())
			} else if ch == '.' {
				if len(intPart) == 0 {
					panic(NewJsonParserError("Unexpected char", t.reader.readed))
				}
				t.reader.Next()
				hasFraPart = true
				status = READ_NUMBER_FRA_PART
			} else if ch == 'e' || ch == 'E' {
				t.reader.Next()
				hasExpPart = true
				signChar := t.reader.Peek()
				if signChar == '-' || signChar == '+' {
					expMinusSign = signChar == '-'
					t.reader.Next()
				}
				status = READ_NUMBER_EXP_PART
			} else {
				if len(intPart) == 0 {
					panic(NewJsonParserError("Unexpected char", t.reader.readed))
				}
				status = READ_NUMBER_END
			}
			continue
		case READ_NUMBER_FRA_PART:
			if ch >= '0' && ch <= '9' {
				fraPart = append(fraPart, t.reader.Next())
			} else if ch == 'e' || ch == 'E' {
				t.reader.Next()
				hasExpPart = true
				signChar := t.reader.Peek()
				if signChar == '-' || signChar == '+' {
					expMinusSign = signChar == '-'
					t.reader.Next()
				}
				status = READ_NUMBER_EXP_PART
			} else {
				if len(fraPart) == 0 {
					panic(NewJsonParserError("Unexpected char", t.reader.readed))
				}
				status = READ_NUMBER_END
			}
			continue
		case READ_NUMBER_EXP_PART:
			if ch >= '0' && ch <= '9' {
				expPart = append(expPart, t.reader.Next())
			} else {
				if len(expPart) == 0 {
					panic(NewJsonParserError("Unexpected char", t.reader.readed))
				}
				status = READ_NUMBER_END
			}
			continue
		case READ_NUMBER_END:
			readed := t.reader.readed
			if len(intPart) == 0 {
				panic(NewJsonParserError("Unexpected char", t.reader.readed))
			}
			lint := int64(0)
			if minusSign {
				lint = -1 * string2long(intPart, readed)
			} else {
				lint = string2long(intPart, readed)
			}
			if hasExpPart && len(expPart) == 0 {
				return NewNumber(lint)
			}
			if hasFraPart && len(fraPart) == 0 {
				panic(NewJsonParserError("Unexpected char", t.reader.readed))
			}
			dFraPart := float64(0.0)
			if hasFraPart {
				if minusSign {
					dFraPart = -float64(1.0) * string2Fraction(fraPart, t.reader.readed)
				} else {
					dFraPart = string2Fraction(fraPart, t.reader.readed)
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
			if number > MAX_SAFE_DOUBLE {
				panic(NewJsonParserError("Exceeded maximum value", t.reader.readed))
			}
			return NewNumber(number)
		}
	}

}

var MAX_SAFE_INTEGER = int64(9007199254740991)

func string2long(str []rune, readed int) int64 {
	if len(str) > 16 {
		panic(NewJsonParserError("Number string is too long", readed))
	}
	result := int64(0)
	for _, digit := range str {
		result = result*10 + int64(digit-'0')
		if result > MAX_SAFE_INTEGER {
			panic(NewJsonParserError("Exceeded maximum value", readed))
		}
	}
	return result
}

var MAX_SAFE_DOUBLE = 1.7976931348623157e+308

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
