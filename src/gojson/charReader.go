package gojson

//Copyright 2018 by gau fung. All rights reserved.
//Use of source code is governed by MIT License

import (
	"io"
	"io/ioutil"
	"unicode/utf8"
)

//CharReader struct
//Including buffer, reader, pos and size
type charReader struct {
	buffer []byte
	reader io.Reader
	pos    int
	size   int
}

//create CharReader
func newCharReader(r io.Reader) *charReader {
	reader := &charReader{reader: r}
	reader.fillBuffer()
	return reader
}

// fill buffer
func (r *charReader) fillBuffer() {
	if buffer, err := ioutil.ReadAll(r.reader); err == nil {
		r.buffer = buffer
		r.pos = 0
		r.size = len(r.buffer)
	} else {
		r.buffer = make([]byte, 0)
	}
}

//determine whether has more byte to read, if yes return true, or else false
func (r *charReader) hasMore() bool {
	return r.pos < r.size
}

// next char, using utf-8 encoding
func (r *charReader) next() rune {
	bytes := make([]byte, 0)
	for {
		bytes = append(bytes, r.nextByte())
		if utf8.Valid(bytes) {
			r, _ := utf8.DecodeRune(bytes)
			return r
		}
	}
}

//backward previous char, using utf-8 encoding
func (r *charReader) backward() {
	index := 1
	for {
		if utf8.Valid(r.buffer[r.pos-index : r.pos]) {
			r.pos = r.pos - index
			return
		} else {
			index++
		}
	}
}

//next byte
func (r *charReader) nextByte() byte {
	ch := r.buffer[r.pos]
	r.pos++
	return ch
}

//peek next char, using utf-8 encoding
func (r *charReader) peek() rune {
	idx := r.pos
	bytes := make([]byte, 0)
	for {
		bytes = append(bytes, r.buffer[idx])
		if utf8.Valid(bytes) {
			r, _ := utf8.DecodeRune(bytes)
			return r
		} else {
			idx++
		}
	}
}
