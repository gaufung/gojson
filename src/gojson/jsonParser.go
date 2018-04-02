package gojson

import (
	"io"
	"strings"
	"reflect"
)

type JsonParser struct {
	reader *tokenReader
}


func NewJsonParser(reader io.Reader) *JsonParser {
	return &JsonParser{reader:newTokenReader(reader)}
}

func NewJsonParserFromString(json string) *JsonParser{
	return NewJsonParser(strings.NewReader(json))
}

func (p *JsonParser) Parse() (interface{}, error){
	return parse(p.reader)
}

func (p *JsonParser) Decode(target interface{}) error {
	if object, err := p.Parse(); err!=nil{
		return err
	}else{
		return decode(object,target)
	}
}

func decode(object, target interface{}) error {
	panic("Have not implemented")
}