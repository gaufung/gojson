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
	reader := newCharReader(stringReader)
	if reader == nil {
		t.Error("construct failed")
	}
	if !reader.hasMore() {
		t.Error("HasMore() methods failed")
	}
	substring := reader.next()
	if substring != '语' {
		t.Error(substring)
		t.Error("Next() method failed")
	}
	nextChar := reader.peek()
	if nextChar != '言' {
		t.Error("Peek failed")
	}
	nextChar = reader.next()
	if nextChar != '言' {
		t.Error("NextChar() failed")
	}
	nextChar = reader.next()
	if nextChar != 'A' {
		t.Error("NextChar() failed")
	}
}
