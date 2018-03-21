package gojson

import (
	"io"
	"unicode/utf8"
)

type CharReader struct {
	buffer []byte
	reader io.Reader
	readed int
	pos    int
	size   int
}

func NewCharReader(r io.Reader) *CharReader {
	buffer_size := 1024
	return &CharReader{
		buffer: make([]byte, buffer_size),
		reader: r}
}

func (r *CharReader) HasMore() bool {
	if r.pos < r.size {
		return true
	} else {
		r.fillBuffer(nil)
		return r.pos < r.size
	}
}

func (r *CharReader) NextString(size int) string {
	result := make([]rune, 0)
	for size > 0 {
		result = append(result, r.Next())
		size--
	}
	return string(result)
}

func (r *CharReader) Next() rune {
	bytes := make([]byte, 0)
	for {
		bytes = append(bytes, r.NextByte())
		if utf8.Valid(bytes) {
			r, _ := utf8.DecodeRune(bytes)
			return r
		}
	}
}

func (r *CharReader) NextByte() byte {
	if r.pos == r.size {
		r.fillBuffer("EOF")
	}
	ch := r.buffer[r.pos]
	r.pos++
	return ch
}

func (r *CharReader) Peek() rune {
	idx := r.pos
	bytes := make([]byte, 0)
	for {
		if idx == r.size {
			r.fillBuffer("EOF")
		}
		bytes = append(bytes, r.buffer[idx])
		if utf8.Valid(bytes) {
			r, _ := utf8.DecodeRune(bytes)
			return r
		} else {
			idx++
		}
	}
}

func (r *CharReader) fillBuffer(eofErrorMessage interface{}) {
	if n, err := r.reader.Read(r.buffer); err == nil {
		r.pos = 0
		r.size = n
		r.readed += n
	} else {
		if eofErrorMessage != nil {
			panic(eofErrorMessage)
		}
	}
}
