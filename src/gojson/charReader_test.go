package gojson

import (
	"fmt"
	"strings"
	"testing"
)

func TestCharReader(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("EOF")
		}
	}()

	stringReader := strings.NewReader("语言A fox jumps over the lazy dogΦx")
	reader := NewCharReader(stringReader)
	if reader == nil {
		t.Error("construct failed")
	}
	if !reader.HasMore() {
		t.Error("HasMore() methods failed")
	}
	substring := reader.Next()
	if substring != '语' {
		t.Error(substring)
		t.Error("Next() method failed")
	}
	nextChar := reader.Peek()
	if nextChar != '言' {
		t.Error("Peek failed")
	}
	nextChar = reader.Next()
	if nextChar != '言' {
		t.Error("NextChar() failed")
	}
	nextChar = reader.Next()
	if nextChar != 'A' {
		t.Error("NextChar() failed")
	}
	str := reader.NextString(30)
	if str != " fox jumps over the lazy dogΦx" {

		t.Error("Next() method failed")
	}
	if reader.HasMore() {
		t.Error("HasMore() method failed")
	}
}
