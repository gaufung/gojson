package gojson

import (
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
)

const (
	READ_NUMBER_INT_PART = iota
	READ_NUMBER_FRA_PART
	READ_NUMBER_EXP_PART
)

const (
	READ_NUMBER_PLUS  = 0x0001
	READ_NUMBER_MINUS = 0x0002
	READ_NUMBER_DIGIT = 0X0004
	READ_NUMBER_DOT   = 0x0008
	READ_NUMBER_E     = 0x0010
	READ_NUMBER_END = 0x0100
)

type tokenReader struct {
	reader *charReader
}

func newTokenReaderFromString(s string) *tokenReader {
	return newTokenReader(strings.NewReader(s))
}

func newTokenReader(r io.Reader) *tokenReader {
	return &tokenReader{newCharReader(r)}
}

func (t *tokenReader) position() int {
	return t.reader.pos
}

func (t *tokenReader) errors(info string) error {
	return errors.New(info + " at " + strconv.Itoa(t.position()))
}

// white space to ignore
func (t *tokenReader) isWhiteSpace(ch rune) bool {
	return ch == '\n' || ch == '\t' || ch == ' ' || ch == '\r'
}

// read next token
func (t *tokenReader) readNextToken() (Token, error) {
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
func (t *tokenReader) readString() (string, error) {
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
func (t *tokenReader) readBoolean() (bool, error) {
	ch := t.reader.next()
	expect := ""
	if ch == 't' {
		expect = "rue"
	} else if ch == 'f' {
		expect = "alse"
	} else {
		return false, t.errors("Read boolean: unexpected char")
	}
	for _, c := range []rune(expect) {
		theChar := t.reader.next()
		if theChar != c {
			return false, t.errors("Read boolean: unexpected char")
		}
	}
	return ch == 't', nil
}

func (t *tokenReader) readNull() error {
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
func (t *tokenReader) readNumber() (float64, error) {
	intPart, fraPart, expPart := make([]rune, 0), make([]rune, 0), make([]rune, 0)
	minusSign, expMinusSign := false, false
	phase := READ_NUMBER_INT_PART
	status := READ_NUMBER_MINUS | READ_NUMBER_PLUS | READ_NUMBER_DIGIT

	for {
		if !t.reader.hasMore() {
			break
		}
		ch := t.reader.peek()
		token := numberToken(ch)
		if token == READ_NUMBER_END {
			break
		}
		switch token {
		case READ_NUMBER_MINUS:
			if phase == READ_NUMBER_INT_PART {
				if status&READ_NUMBER_MINUS > 0 {
					t.reader.next()
					minusSign = true
					status = READ_NUMBER_DIGIT | READ_NUMBER_DOT
					continue

				}
			} else if phase == READ_NUMBER_EXP_PART {
				if status&READ_NUMBER_MINUS > 0 {
					t.reader.next()
					expMinusSign = true
					status = READ_NUMBER_DIGIT
					continue
				}

			}
			return 0.0, t.errors("Read float64: unexpect char")
		case READ_NUMBER_PLUS:
			if phase == READ_NUMBER_INT_PART {
				if status&READ_NUMBER_PLUS > 0 {
					t.reader.next()
					status = READ_NUMBER_DIGIT | READ_NUMBER_DOT
					continue
				}
			} else if phase == READ_NUMBER_EXP_PART {
				if status&READ_NUMBER_PLUS > 0 {
					t.reader.next()
					status = READ_NUMBER_DIGIT
					continue
				}
			}
			return 0.0, t.errors("Read float64: unexpect char")
		case READ_NUMBER_DOT:
			if phase == READ_NUMBER_INT_PART {
				if status&READ_NUMBER_DOT > 0 {
					t.reader.next()
					phase = READ_NUMBER_FRA_PART
					status = READ_NUMBER_DIGIT | READ_NUMBER_E
					continue
				}
			}
			return 0.0, t.errors("Read float64: unexpect char")
		case READ_NUMBER_DIGIT:
			if phase == READ_NUMBER_INT_PART {
				if status&READ_NUMBER_DIGIT > 0 {
					char := t.reader.next()
					intPart = append(intPart, char)
					status = READ_NUMBER_DIGIT | READ_NUMBER_DOT | READ_NUMBER_E
					continue
				}
			}
			if phase == READ_NUMBER_FRA_PART {
				if status&READ_NUMBER_DIGIT > 0 {
					char := t.reader.next()
					fraPart = append(fraPart, char)
					status = READ_NUMBER_DIGIT | READ_NUMBER_E
					continue
				}
			}
			if phase == READ_NUMBER_EXP_PART {
				if status&READ_NUMBER_DIGIT > 0 {
					char := t.reader.next()
					expPart = append(expPart, char)
					status = READ_NUMBER_DIGIT
					continue
				}
			}
			return 0.0, t.errors("Read float64: unexpect char")
		case READ_NUMBER_E:
			if phase == READ_NUMBER_INT_PART || phase == READ_NUMBER_FRA_PART {
				if status&READ_NUMBER_E > 0 {
					t.reader.next()
					phase = READ_NUMBER_EXP_PART
					status = READ_NUMBER_PLUS | READ_NUMBER_MINUS | READ_NUMBER_DIGIT
					continue
				}
			}
			return 0.0, t.errors("Read float64: unexpect char")
		}
	}
	fraction := fracPartConvert(fraPart)
	firstPart := 0.0
	if fraction == 0.0 {
		firstPart = intPartConvert(intPart)
	} else {
		firstPart = intPartConvert(intPart) + fraction
	}
	secondPart := intPartConvert(expPart)
	if minusSign {
		firstPart = -1.0 * firstPart
	}
	if expMinusSign {
		secondPart = -1.0 * secondPart
	}
	if secondPart == 0.0 {
		return firstPart, nil
	} else {
		return firstPart * math.Pow(10.0, secondPart), nil
	}
}

func intPartConvert(digits []rune) float64 {
	result := 0.0
	for _, digit := range digits {
		result = result*10.0 + float64(digit-'0')
	}
	return result
}
func fracPartConvert(digits []rune) float64 {
	result := 0.0
	factor := 0.1
	for _, digit := range digits {
		result = result + factor*float64(digit-'0')
		factor = factor * 0.1
	}
	return result
}

func numberToken(ch rune) int {
	if ch >= '0' && ch <= '9' {
		return READ_NUMBER_DIGIT
	} else if ch == '.' {
		return READ_NUMBER_DOT
	} else if ch == 'e' || ch == 'E' {
		return READ_NUMBER_E
	} else if ch == '+' {
		return READ_NUMBER_PLUS
	} else if ch == '-' {
		return READ_NUMBER_MINUS
	} else {
		return READ_NUMBER_END
	}
}

//back token
func (t *tokenReader) backToken() {
	t.reader.backward()
}

//determine whether is empty or not
func (t *tokenReader) isEmpty() bool {
	if !t.reader.hasMore() {
		return true
	} else {
		return false
	}
}
