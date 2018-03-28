package gojson

type stack struct {
	Size  int
	pos   int
	array []*stackValue
}

func newStack() *stack {
	size := 100
	return &stack{Size: size, pos: 0, array: make([]*stackValue, size)}
}

func (s *stack) isEmpty() bool {
	return s.pos == 0
}

func (s *stack) push(obj *stackValue) bool {
	if s.pos == s.Size {
		return false
	}
	s.array[s.pos] = obj
	s.pos++
	return true
}

func (s *stack) pop() (*stackValue, bool) {
	if s.isEmpty() {
		return nil, false
	}
	s.pos--
	return s.array[s.pos], true
}

func (s *stack) popKind(kind int) (*stackValue, bool) {
	if s.isEmpty() {
		return nil, false
	}
	s.pos--
	obj := s.array[s.pos]
	if obj.kind == kind {
		return obj, true
	} else {
		return nil, false
	}
}

func (s *stack) getTopValueType() int {
	obj := s.array[s.pos-1]
	return obj.kind
}

func (s *stack) peek(kind int) (*stackValue, bool) {
	if s.isEmpty() {
		return nil, false
	}
	obj := s.array[s.pos-1]
	if obj.kind == kind {
		return obj, true
	}
	return nil, false
}
