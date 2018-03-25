package gojson

import (
	"io"
	"strings"
)

const (
	STATUS_READ_END_DOCUMENT = 0x0002
	STATUS_READ_BEGIN_OBJECT = 0x0004
	STATUS_READ_END_OBJECT   = 0x0008
	STATUS_READ_OBJECT_KEY   = 0x0010
	STATUS_READ_OBJECT_VALUE = 0x0020
	STATUS_READ_COLON        = 0x0040
	STATUS_READ_COMMA        = 0x0080
	STATUS_READ_BEGIN_ARRAY  = 0x0100
	STATUS_READ_END_ARRAY    = 0x0200
	STATUS_READ_ARRAY_VALUE  = 0x0400
)

type JsonStream struct {
	reader    *TokenReader
	stack     *Stack
	lastToken Token
	status    int
}

func NewJsonStreamFromReader(r io.Reader) *JsonStream {
	return &JsonStream{
		reader: &TokenReader{NewCharReader(r)},
	}
}
func NewJsonStreamFromTokenReader(r *TokenReader) *JsonStream {
	return &JsonStream{reader: r}
}

func NewJsonStreamFromString(str string) *JsonStream {
	return NewJsonStreamFromReader(strings.NewReader(str))
}

func (j *JsonStream) hasStatus(expectStatus int) bool {
	return (j.status & expectStatus) != 0
}

func (j *JsonStream) newMap() map[string]interface{} {
	return make(map[string]interface{})
}

func (j *JsonStream) newArray() []interface{} {
	return make([]interface{}, 0)
}

func Parser(r *TokenReader) (interface{}, bool) {
	j := NewJsonStreamFromTokenReader(r)
	j.stack = NewStack()
	j.status = STATUS_READ_BEGIN_OBJECT | STATUS_READ_BEGIN_ARRAY
	for {
		currentToken := j.reader.readNextToken()
		switch currentToken {
		case BOOLEAN:
			b := j.reader.readBoolean()
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.PopKind(TYPE_OBJECT_KEY); ok {
					key := sv.ValueAsKey()
					if tsv, o := j.stack.Peek(TYPE_OBJECT); o {
						tsv.ValueAsObject()[key] = b
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.Peek(TYPE_ARRAY); ok {
					temp := sv.ValueAsArray()
					temp = append(temp, b)
					j.stack.Pop()
					j.stack.Push(NewJsonObjectFromSlice(temp))
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			return nil, false
		case NUMBER:
			number := j.reader.readNumber()
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.PopKind(TYPE_OBJECT_KEY); ok {
					key := sv.ValueAsKey()
					if tsv, o := j.stack.Peek(TYPE_OBJECT); o {
						tsv.ValueAsObject()[key] = number
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.Peek(TYPE_ARRAY); ok {
					temp := sv.ValueAsArray()
					temp = append(temp, number)
					j.stack.Pop()
					j.stack.Push(NewJsonObjectFromSlice(temp))
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			return nil, false
		case NULL:
			j.reader.readNull()
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.PopKind(TYPE_OBJECT_KEY); ok {
					key := sv.ValueAsKey()
					if tsv, o := j.stack.Peek(TYPE_OBJECT); o {
						tsv.ValueAsObject()[key] = nil
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.Peek(TYPE_ARRAY); ok {
					temp := sv.ValueAsArray()
					temp = append(temp, nil)
					j.stack.Pop()
					j.stack.Push(NewJsonObjectFromSlice(temp))
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			return nil, false
		case STRING:
			str := j.reader.readString()
			if j.hasStatus(STATUS_READ_OBJECT_KEY) {
				j.stack.Push(NewJsonObjectFromKey(str))
				j.status = STATUS_READ_COLON
				continue
			}
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.PopKind(TYPE_OBJECT_KEY); ok {
					key := sv.ValueAsKey()
					if tsv, o := j.stack.Peek(TYPE_OBJECT); o {
						tsv.ValueAsObject()[key] = str
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.PopKind(TYPE_ARRAY); ok {
					temp := sv.ValueAsArray()
					temp = append(temp, str)
					j.stack.Push(NewJsonObjectFromSlice(temp))
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			return nil, false
		case COLON_SEPERATOR:
			if j.status == STATUS_READ_COLON {
				j.status = STATUS_READ_OBJECT_VALUE | STATUS_READ_BEGIN_OBJECT | STATUS_READ_BEGIN_ARRAY
				continue
			}
			return nil, false
		case COMA_SEPERATOR:
			if j.hasStatus(STATUS_READ_COMMA) {
				if j.hasStatus(STATUS_READ_END_OBJECT) {
					j.status = STATUS_READ_OBJECT_KEY
					continue
				}
				if j.hasStatus(STATUS_READ_END_ARRAY) {
					j.status = STATUS_READ_ARRAY_VALUE | STATUS_READ_BEGIN_ARRAY | STATUS_READ_BEGIN_OBJECT
					continue
				}
			}
			return nil, false
		case START_OBJECT:
			if j.hasStatus(STATUS_READ_BEGIN_OBJECT) {
				j.stack.Push(NewJsonObjectFromObject(j.newMap()))
				j.status = STATUS_READ_OBJECT_KEY | STATUS_READ_END_OBJECT
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.PopKind(TYPE_ARRAY); ok {
					temp := sv.ValueAsArray()
					j.reader.BackToken()
					if val, ok := Parser(j.reader); ok {
						temp = append(temp, val)
						j.stack.Push(NewJsonObjectFromSlice(temp))
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
						continue
					}
				}
			}
			return nil, false
		case START_ARRAY:
			if j.hasStatus(STATUS_READ_BEGIN_ARRAY) {
				j.stack.Push(NewJsonObjectFromSlice(j.newArray()))
				j.status = STATUS_READ_ARRAY_VALUE
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.PopKind(TYPE_ARRAY); ok {
					temp := sv.ValueAsArray()
					j.reader.BackToken()
					if val, ok := Parser(j.reader); ok {
						temp = append(temp, val)
						j.stack.Push(NewJsonObjectFromSlice(temp))
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
						continue
					}
				}
			}
			return nil, false
		case END_OBJECT:
			if j.hasStatus(STATUS_READ_END_OBJECT) {
				object, _ := j.stack.PopKind(TYPE_OBJECT)
				if j.stack.IsEmpty() {
					j.stack.Push(object)
					//j.status = STATUS_READ_END_DOCUMENT
					if j.reader.IsEmpty() {
						j.status = STATUS_READ_END_DOCUMENT
						continue
					} else {
						return object.Value, true
					}

				}
				kind := j.stack.GetTopValueType()
				if kind == TYPE_OBJECT_KEY {
					if sv, ok := j.stack.PopKind(TYPE_OBJECT_KEY); ok {
						key := sv.ValueAsKey()
						if tsv, o := j.stack.Peek(TYPE_OBJECT); o {
							tsv.ValueAsObject()[key] = object.Value
							j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
						}
					}
					continue
				}
				if kind == TYPE_ARRAY {
					if sv, ok := j.stack.PopKind(TYPE_ARRAY); ok {
						temp := sv.ValueAsArray()
						temp = append(temp, object.Value)
						j.stack.Push(NewJsonObjectFromSlice(temp))
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
					}
					continue
				}
			}
			return nil, false
		case END_ARRAY:
			if j.hasStatus(STATUS_READ_END_ARRAY) {
				array, _ := j.stack.PopKind(TYPE_ARRAY)
				if j.stack.IsEmpty() {
					j.stack.Push(array)
					if j.reader.IsEmpty() {
						j.status = STATUS_READ_END_DOCUMENT
						continue
					} else {
						return array.Value, true
					}
				}
				kind := j.stack.GetTopValueType()
				if kind == TYPE_OBJECT_KEY {
					if sv, ok := j.stack.PopKind(TYPE_OBJECT_KEY); ok {
						key := sv.ValueAsKey()
						if tsv, o := j.stack.Peek(TYPE_OBJECT); o {
							tsv.ValueAsObject()[key] = array.Value
							j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
						}
					}
					continue
				}
				if kind == TYPE_ARRAY {
					if sv, ok := j.stack.PopKind(TYPE_ARRAY); ok {
						temp := sv.ValueAsArray()
						temp = append(temp, array.Value)
						j.stack.Push(NewJsonObjectFromSlice(temp))
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
					}
					continue
				}
			}
			return nil, false
		case END_DOCUMENT:
			if j.hasStatus(STATUS_READ_END_DOCUMENT) {
				v, _ := j.stack.Pop()
				if j.stack.IsEmpty() {
					return v.Value, true
				}
			}
			return nil, false
		}
	}
}
