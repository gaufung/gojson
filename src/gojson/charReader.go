package gojson
//Copyright 2018 by gau fung. All rights reserved.
//Use of source code is governed by MIT License


import (
	"io"
	"unicode/utf8"
	"io/ioutil"
)

//CharReader struct
//Including buffer, reader, pos and size
type CharReader struct {
	buffer []byte
	reader io.Reader
	pos    int
	size   int
}

//create CharReader
func newCharReader(r io.Reader) *CharReader {
	reader := &CharReader{reader:r}
	reader.fillBuffer()
	return reader
}

// fill buffer
func (r *CharReader) fillBuffer() {
	if buffer, err:=ioutil.ReadAll(r.reader); err==nil{
		r.buffer=buffer
		r.pos = 0
		r.size = len(r.buffer)
	}else{
		r.buffer = make([]byte,0)
	}
}

//determine whether has more byte to read, if yes return true, or else false
func (r *CharReader) hasMore() bool {
	return r.pos < r.size
}

// next char, using utf-8 encoding
func (r *CharReader) next() rune {
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
func (r *CharReader) backward() {
	index := 1
	for {
		if utf8.Valid(r.buffer[r.pos-index:r.pos]){
			r.pos = r.pos -index
			return
		}else{
			index++
		}
	}
}

//next byte
func (r *CharReader) nextByte() byte {
	ch := r.buffer[r.pos]
	r.pos++
	return ch
}

//peek next char, using utf-8 encoding
func (r *CharReader) peek() rune {
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

