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

	stringReader := strings.NewReader("A fox jumps over the lazy dog")

	reader := NewCharReader(stringReader)

	if reader == nil {
		t.Error("construct failed")
	}

	if !reader.HasMore() {
		t.Error("HasMore() methods failed")
	}

	substring := reader.Next(5)
	if substring != "A fox" {
		t.Error("Next() method failed")
	}

	nextChar := reader.Peek()
	if nextChar != byte(' ') {
		t.Error("Peek failed")
	}
	substring = reader.Next(6)
	if substring != " jumps" {
		t.Error("Next() method failed")
	}
	substring = reader.Next(18)
	if substring != " over the lazy dog" {
		t.Error("Next() method failed")
	}
	if reader.HasMore() {
		t.Error("HasMore() method failed")
	}
	reader.Next(10)
	t.Errorf("Next() has failed")
}
