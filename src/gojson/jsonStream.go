package gojson

import (
	"io"
	"strings"
)

const (
	STATUS_READ_END_DOCUMENT = 0x0002
	STATUS_READ_BEGIN_OBJECT = 0x0004
	STATUS_READ_END_OBJECT = 0x0008
	STATUS_READ_OBJECT_KEY = 0x0010
	STATUS_READ_OBJECT_VALUE = 0x0020
	STATUS_READ_COLON = 0x0040
	STATUS_READ_COMMA= 0x0080
	STATUS_READ_BEGIN_ARRAY = 0x0100
	STATUS_READ_END_ARRAY = 0x0200
	STATUS_READ_ARRAY_VALUE = 0x0400
	STATUS_READ_SINGLE_VALUE = 0x0800
)

type JsonStream struct {
	reader *TokenReader
	stack *Stack
	isSingleValue bool
	lastToken Token
	status int
}

func NewJsonStreamFromReader(r io.Reader) *JsonStream {
	return &JsonStream{
		reader:&TokenReader{NewCharReader(r)},
	}
}

func NewJsonStreamFromString(str string) *JsonStream {
	return NewJsonStreamFromReader(strings.NewReader(str))
}

func (j *JsonStream) hasStatus(expectStatus int) bool{
	return j.status == expectStatus
}

func (j *JsonStream) newMap() map[string] interface{} {
	return make(map[string]interface{})
}

func (j *JsonStream) newArray() []interface{} {
	return make([]interface{}, 0)
}

//func (j *JsonStream) checkExpectedType(obj interface{}, t types.Map) bool{
//
//}


func (j *JsonStream) parse() interface{} {
	j.stack = NewStack()
	j.status = STATUS_READ_SINGLE_VALUE | STATUS_READ_BEGIN_OBJECT | STATUS_READ_BEGIN_ARRAY
	for {
		currentToken := j.reader.readNextToken()
		switch currentToken {
		case BOOLEAN:
			b := j.reader.readBoolean()
			if j.hasStatus(STATUS_READ_SINGLE_VALUE) {
				j.stack.Push(NewJsonObjectFromSingle(b))
				j.status = STATUS_READ_END_DOCUMENT;
				continue
			}
			if j.status == STATUS_READ_OBJECT_VALUE {
				if sv, err := j.stack.PopKind(TYPE_OBJECT_KEY); err!=nil{
					key:= sv.ValueAsKey()
					if tsv, e:=j.stack.Peek(TYPE_OBJECT);e!=nil{
						tsv.ValueAsObject()[key]=b
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.status == STATUS_READ_ARRAY_VALUE {
				if sv, err := j.stack.Peek(TYPE_ARRAY); err!=nil{
					temp := sv.ValueAsArray()
					temp = append(temp, b)
					//sv.ValueAsArray()=append(sv.ValueAsArray(),b)
					j.status=STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			panic(NewJsonParserError("should not reach here.", 0))
		case NUMBER:
			number := j.reader.readNumber()
			if j.hasStatus(STATUS_READ_SINGLE_VALUE) {
				j.stack.Push(NewJsonObjectFromObject(number))
				j.status = STATUS_READ_END_DOCUMENT
				continue
			}
			if j.status == STATUS_READ_OBJECT_VALUE {
				if sv, err := j.stack.PopKind(TYPE_OBJECT_KEY); err!=nil{
					key := sv.ValueAsKey()
					if tsv, e:= j.stack.Peek(TYPE_OBJECT); e!=nil{
						tsv.ValueAsObject()[key]=number
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.status == STATUS_READ_ARRAY_VALUE {
				if sv, err := j.stack.Peek(TYPE_ARRAY); err!=nil{
					temp := sv.ValueAsArray()
					temp = append(temp, number)
					//sv.ValueAsArray()=append(sv.ValueAsArray(), number)
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			panic(NewJsonParserError("should not reach here.", 0))
		case NULL:
			j.reader.readNull()
			if j.hasStatus(STATUS_READ_SINGLE_VALUE) {
				j.stack.Push(NewJsonObjectFromObject(nil))
				j.status = STATUS_READ_END_DOCUMENT
				continue
			}
			if j.status == STATUS_READ_OBJECT_VALUE {
				if sv, err:=j.stack.PopKind(TYPE_OBJECT_KEY); err!=nil{
					key := sv.ValueAsKey()
					if tsv, e:=j.stack.Peek(TYPE_OBJECT); e!=nil{
						tsv.ValueAsObject()[key]=nil
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.status == STATUS_READ_ARRAY_VALUE {
				if sv, err := j.stack.Peek(TYPE_ARRAY); err!=nil {
					temp := sv.ValueAsArray()
					temp = append(temp, nil)
					//sv.ValueAsArray() = append(sv.ValueAsArray(), nil)
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			panic(NewJsonParserError("should not reach here.", 0))
		case STRING:
			str := j.reader.readString()
			if j.hasStatus(STATUS_READ_SINGLE_VALUE) {
				j.stack.Push(NewJsonObjectFromObject(str))
				j.status = STATUS_READ_END_DOCUMENT
				continue
			}
			if j.status == STATUS_READ_OBJECT_KEY {
				j.stack.Push(NewJsonObjectFromKey(str))
				j.status = STATUS_READ_COLON
				continue
			}
			if j.status == STATUS_READ_OBJECT_VALUE {
				if sv, err := j.stack.PopKind(TYPE_OBJECT_KEY); err!=nil{
					key := sv.ValueAsKey()
					if tsv, e := j.stack.Peek(TYPE_OBJECT); e!=nil{
						tsv.ValueAsObject()[key]= str
						j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
					}
				}
				continue
			}
			if j.status == STATUS_READ_ARRAY_VALUE {
				if sv, err:= j.stack.Peek(TYPE_ARRAY); err!=nil {
					temp := sv.ValueAsArray()
					temp = append(temp, str)
					//sv.ValueAsArray() = append(sv.ValueAsArray(), str)
					j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
				}
				continue
			}
			panic(NewJsonParserError("should not reach here.", 0))
		case COLON_SEPERATOR:
			if j.status == STATUS_READ_COLON {
				j.status = STATUS_READ_OBJECT_VALUE;
				continue
			}
			panic(NewJsonParserError("should not reach here.", 0))
		case COMA_SEPERATOR:
			if j.hasStatus(STATUS_READ_COMMA) {
				if j.hasStatus(STATUS_READ_END_OBJECT){
					j.status = STATUS_READ_OBJECT_KEY
					continue
				}
				if j.hasStatus(STATUS_READ_END_ARRAY) {
					j.status =STATUS_READ_ARRAY_VALUE
					continue
				}
			}
			panic(NewJsonParserError("should not reach here.", 0))
		case END_ARRAY:
			if j.hasStatus(STATUS_READ_END_ARRAY) {
				array,_ := j.stack.PopKind(TYPE_ARRAY)
				if j.stack.IsEmpty() {
					j.stack.Push(array)
					j.status = STATUS_READ_END_DOCUMENT
					continue
				}
				kind := j.stack.GetTopValueType()
				if kind == TYPE_OBJECT_KEY {
					if sv, err := j.stack.PopKind(TYPE_OBJECT_KEY); err!=nil {
						key := sv.ValueAsKey()
						if tsv, e:=j.stack.Peek(TYPE_OBJECT); e!=nil{
							tsv.ValueAsObject()[key]=array.value
							j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
						}
					}
					continue
				}
				if kind == TYPE_ARRAY {
					if sv, err:=j.stack.PopKind(TYPE_ARRAY); err!=nil{
						temp := sv.ValueAsArray()
						temp = append(temp, array.value)
						//sv.ValueAsArray() = append(sv.ValueAsArray(), array.value)
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
					}
					continue
				}
			}
			panic(NewJsonParserError("should not reach here.", 0))
		//case END_OBJECT:
		//	if j.hasStatus(STATUS_READ_END_OBJECT) {
		//		object,_ := j.stack.PopKind(TYPE_OBJECT)
		//		if j.stack.IsEmpty() {
		//			j.stack.Push(object)
		//			j.status = STATUS_READ_END_DOCUMENT
		//			continue
		//		}
		//		kind := j.stack.GetTopValueType()
		//		if kind == TYPE_OBJECT_KEY {
		//			if sv, err:= j.stack.PopKind(TYPE_OBJECT_KEY); err!=nil{
		//				key:= sv.ValueAsKey()
		//				if tsv, e:=j.stack.Peek(TYPE_OBJECT); e!=nil{
		//					tsv.ValueAsObject()[key]=object.value
		//					j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
		//				}
		//			}
		//			continue
		//		}
		//		if kind == TYPE_ARRAY {
		//			if sv, err:= j.stack.PopKind(TYPE_ARRAY); err!=nil{
		//				sv.ValueAsArray()=append(sv.ValueAsArray(), object.value)
		//				j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
		//			}
		//			continue
		//		}
		//	}
		//	panic(NewJsonParserError("should not reach here.", 0))
		case END_OBJECT:
			if j.hasStatus(STATUS_READ_END_OBJECT) {
				object, _ := j.stack.PopKind(TYPE_OBJECT)
				if j.stack.IsEmpty(){
					j.stack.Push(object)
					j.status = STATUS_READ_END_DOCUMENT
					continue
				}
				kind := j.stack.GetTopValueType()
				if kind == TYPE_OBJECT_KEY {
					if sv, err := j.stack.PopKind(TYPE_OBJECT_KEY); err!=nil {
						key := sv.ValueAsKey()
						if tsv, e:= j.stack.Peek(TYPE_OBJECT); e!=nil{
							tsv.ValueAsObject()[key]=object.value
							j.status = STATUS_READ_COMMA | STATUS_READ_END_OBJECT
						}
					}
					continue
				}
				if kind == TYPE_ARRAY {
					if sv, err:= j.stack.PopKind(TYPE_ARRAY); err!=nil{
						temp := sv.ValueAsArray()
						temp = append(temp, object.value)
						//sv.ValueAsArray()=append(sv.ValueAsArray(), object.value)
						j.status = STATUS_READ_COMMA | STATUS_READ_END_ARRAY
					}
					continue
				}
			}
			panic(NewJsonParserError("should not reach here.", 0))
		case END_DOCUMENT:
			if j.hasStatus(STATUS_READ_END_DOCUMENT) {
				v,_ := j.stack.Pop()
				if j.stack.IsEmpty(){
					return v.value
				}
			}
			panic(NewJsonParserError("Unexpected EOF.", 0))
		case START_ARRAY:
			if j.hasStatus(STATUS_READ_BEGIN_ARRAY){
				j.stack.Push(NewJsonObjectFromSlice(j.newArray()))
				j.status = STATUS_READ_ARRAY_VALUE
				continue
			}
			panic(NewJsonParserError("Unexpected [.", 0))
		case START_OBJECT:
			if j.hasStatus(STATUS_READ_BEGIN_OBJECT) {
				j.stack.Push(NewJsonObjectFromSingle(j.newMap()))
				j.status= STATUS_READ_OBJECT_KEY
				continue
			}
			panic(NewJsonParserError("Unexpected {.", 0))

		}
	}
}