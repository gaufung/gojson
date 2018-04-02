package gojson

import (
	"errors"
	"strconv"
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
	reader    *tokenReader
	stack     *stack
	lastToken Token
	status    int
}

func newJsonStreamFromTokenReader(r *tokenReader) *JsonStream {
	return &JsonStream{reader: r}
}

func (j *JsonStream) hasStatus(expectStatus int) bool {
	return (j.status & expectStatus) > 0
}

func (j *JsonStream) newMap() map[string]interface{} {
	return make(map[string]interface{})
}

func (j *JsonStream) newArray() []interface{} {
	return make([]interface{}, 0)
}

func parse(r *tokenReader) (interface{}, error) {
	j := newJsonStreamFromTokenReader(r)
	j.stack = newStack()
	j.status = STATUS_READ_BEGIN_OBJECT | STATUS_READ_BEGIN_ARRAY
	for {
		currentToken, err := j.reader.readNextToken()
		if err != nil {
			return nil, err
		}
		switch currentToken {
		case BOOLEAN:
			b, err := j.reader.readBoolean()
			if err != nil {
				return nil, err
			}
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_OBJECT_KEY); ok {
					key := sv.valueAsKey()
					if tsv, o := j.stack.peek(TYPE_OBJECT); o {
						tsv.valueAsObject()[key] = b
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_ARRAY); ok {
					temp := sv.valueAsArray()
					temp = append(temp, b)
					j.stack.push(newStackValueFromSlice(temp))
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			return nil, errors.New("Read boolean failed at " + strconv.Itoa(r.position()))
		case NUMBER:
			number, err := j.reader.readNumber()
			if err != nil {
				return nil, err
			}
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_OBJECT_KEY); ok {
					key := sv.valueAsKey()
					if tsv, o := j.stack.peek(TYPE_OBJECT); o {
						tsv.valueAsObject()[key] = number
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_ARRAY); ok {
					temp := sv.valueAsArray()
					temp = append(temp, number)
					j.stack.push(newStackValueFromSlice(temp))
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			return nil, errors.New("Read number failed at " + strconv.Itoa(r.position()))
		case NULL:
			j.reader.readNull()
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_OBJECT_KEY); ok {
					key := sv.valueAsKey()
					if tsv, o := j.stack.peek(TYPE_OBJECT); o {
						tsv.valueAsObject()[key] = nil
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_ARRAY); ok {
					temp := sv.valueAsArray()
					temp = append(temp, nil)
					j.stack.push(newStackValueFromSlice(temp))
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			return nil, errors.New("Read null failed at " + strconv.Itoa(r.position()))
		case STRING:
			str, _ := j.reader.readString()
			if j.hasStatus(STATUS_READ_OBJECT_KEY) {
				j.stack.push(newStackValueFromKey(str))
				j.status = STATUS_READ_COLON
				continue
			}
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_OBJECT_KEY); ok {
					key := sv.valueAsKey()
					if tsv, o := j.stack.peek(TYPE_OBJECT); o {
						tsv.valueAsObject()[key] = str
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_ARRAY); ok {
					temp := sv.valueAsArray()
					temp = append(temp, str)
					j.stack.push(newStackValueFromSlice(temp))
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			return nil, errors.New("Read string failed at " + strconv.Itoa(r.position()))
		case COLON_SEPERATOR:
			if j.status == STATUS_READ_COLON {
				j.status = STATUS_READ_OBJECT_VALUE | STATUS_READ_BEGIN_ARRAY
				continue
			}
			return nil, errors.New("Read colon seperator failed at " + strconv.Itoa(r.position()))
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
			return nil, errors.New("Read coma seperator failed at " + strconv.Itoa(r.position()))
		case START_OBJECT:
			if j.hasStatus(STATUS_READ_BEGIN_OBJECT) {
				j.stack.push(newStackValueFromObject(j.newMap()))
				j.status = STATUS_READ_OBJECT_KEY | STATUS_READ_END_OBJECT
				continue
			}
			if j.hasStatus(STATUS_READ_OBJECT_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_OBJECT_KEY); ok {
					key := sv.valueAsKey()
					j.reader.backToken()
					if val, err := parse(j.reader); err == nil {
						if tsv, o := j.stack.peek(TYPE_OBJECT); o{
							tsv.valueAsObject()[key] =val
							j.status =  STATUS_READ_COMMA | STATUS_READ_END_OBJECT
							continue
						}
					}
				}
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_ARRAY); ok {
					temp := sv.valueAsArray()
					j.reader.backToken()
					if val, err := parse(j.reader); err == nil {
						temp = append(temp, val)
						j.stack.push(newStackValueFromSlice(temp))
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
						continue
					}
				}
			}
			return nil, errors.New("Read { failed at " + strconv.Itoa(r.position()))
		case START_ARRAY:
			if j.hasStatus(STATUS_READ_BEGIN_ARRAY) {
				j.stack.push(newStackValueFromSlice(j.newArray()))
				j.status = STATUS_READ_ARRAY_VALUE
				continue
			}
			if j.hasStatus(STATUS_READ_ARRAY_VALUE) {
				if sv, ok := j.stack.popKind(TYPE_ARRAY); ok {
					temp := sv.valueAsArray()
					j.reader.backToken()
					if val, err := parse(j.reader); err != nil {
						temp = append(temp, val)
						j.stack.push(newStackValueFromSlice(temp))
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
						continue
					} else {
						return nil, err
					}
				}
			}
			return nil, errors.New("Read [ failed at " + strconv.Itoa(r.position()))
		case END_OBJECT:
			if j.hasStatus(STATUS_READ_END_OBJECT) {
				object, _ := j.stack.popKind(TYPE_OBJECT)
				if j.stack.isEmpty() {
					j.stack.push(object)
					if j.reader.isEmpty() {
						j.status = STATUS_READ_END_DOCUMENT
						continue
					} else {
						return object.value, nil
					}

				}
				kind := j.stack.getTopValueType()
				if kind == TYPE_OBJECT_KEY {
					if sv, ok := j.stack.popKind(TYPE_OBJECT_KEY); ok {
						key := sv.valueAsKey()
						if tsv, o := j.stack.peek(TYPE_OBJECT); o {
							tsv.valueAsObject()[key] = object.value
							j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
						}
					}
					continue
				}
				if kind == TYPE_ARRAY {
					if sv, ok := j.stack.popKind(TYPE_ARRAY); ok {
						temp := sv.valueAsArray()
						temp = append(temp, object.value)
						j.stack.push(newStackValueFromSlice(temp))
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
					}
					continue
				}
			}
			return nil, errors.New("Read } failed at " + strconv.Itoa(r.position()))
		case END_ARRAY:
			if j.hasStatus(STATUS_READ_END_ARRAY) {
				array, _ := j.stack.popKind(TYPE_ARRAY)
				if j.stack.isEmpty() {
					j.stack.push(array)
					if j.reader.isEmpty() {
						j.status = STATUS_READ_END_DOCUMENT
						continue
					} else {
						return array.value, nil
					}
				}
				kind := j.stack.getTopValueType()
				if kind == TYPE_OBJECT_KEY {
					if sv, ok := j.stack.popKind(TYPE_OBJECT_KEY); ok {
						key := sv.valueAsKey()
						if tsv, o := j.stack.peek(TYPE_OBJECT); o {
							tsv.valueAsObject()[key] = array.value
							j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
						}
					}
					continue
				}
				if kind == TYPE_ARRAY {
					if sv, ok := j.stack.popKind(TYPE_ARRAY); ok {
						temp := sv.valueAsArray()
						temp = append(temp, array.value)
						j.stack.push(newStackValueFromSlice(temp))
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
					}
					continue
				}
			}
			return nil, errors.New("Read ] failed at " + strconv.Itoa(r.position()))
		case END_DOCUMENT:
			if j.hasStatus(STATUS_READ_END_DOCUMENT) {
				v, _ := j.stack.pop()
				if j.stack.isEmpty() {
					return v.value, nil
				}
			}
			return nil, errors.New("Read document failed at " + strconv.Itoa(r.position()))
		}
	}
}
