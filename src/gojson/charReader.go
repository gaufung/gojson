package gojson

import (
	"io"
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

func (r *CharReader) Next(size int) string {
	result := make([]byte, 0)
	for size > 0 {
		result = append(result, r.NextChar())
		size--
	}
	return string(result)
}

func (r *CharReader) NextChar() byte {
	if r.pos == r.size {
		r.fillBuffer("EOF")
	}
	ch := r.buffer[r.pos]
	r.pos++
	return ch
}

func (r *CharReader) Peek() byte {
	if r.pos == r.size {
		r.fillBuffer("EOF")
	}
	return r.buffer[r.pos]
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
