package gojson

import (
	"io"
	"strings"
	"bytes"
	"reflect"
	"fmt"
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


func Encode(object interface{}) ([]byte, error){
	var buf bytes.Buffer
	if err:=encode(&buf, reflect.ValueOf(object)); err!=nil{
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v reflect.Value) error{
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("null")
	case reflect.Int,reflect.Int8, reflect.Int16,reflect.Int32,reflect.Int64:
		fmt.Fprintf(buf,"%d", v.Int())
	case reflect.Uint, reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.String:
		fmt.Fprintf(buf,"%q", v.String())
	case reflect.Ptr:
		return encode(buf,v.Elem())
	case reflect.Array, reflect.Slice:
		buf.WriteByte('[')
		for i:=0; i< v.Len();i++{
			if i>0{
				buf.WriteByte(',')
			}
			if err:=encode(buf,v.Index(i)); err!=nil{
				return err
			}
		}
		buf.WriteByte(']')
	case reflect.Struct:
		buf.WriteByte('{')
		for i:=0; i<v.NumField(); i++{
			if i>0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "\"%s\":", v.Type().Field(i).Name)
			if err:=encode(buf,v.Field(i)); err!=nil{
				return err
			}
		}
		buf.WriteByte('}')
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf,"%g",v.Float())
	case reflect.Bool:
		fmt.Fprintf(buf, "%t", v.Bool())
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
