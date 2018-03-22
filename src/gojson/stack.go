package gojson

type Stack struct {
	Size  int
	pos   int
	array []*StackValue
}

func NewStack() *Stack {
	size := 100
	return &Stack{Size: size, pos: 0, array: make([]*StackValue, size)}
}

func (s *Stack) IsEmpty() bool {
	return s.pos == 0
}

func (s *Stack) Push(obj *StackValue) bool {
	if s.pos == s.Size {
		return false
	}
	s.array[s.pos] = obj
	s.pos++
	return true
}

func (s *Stack) Pop() (*StackValue, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	s.pos--
	return s.array[s.pos], true
}

func (s *Stack) PopKind(kind int) (*StackValue, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	s.pos--
	obj := s.array[s.pos]
	if obj.Kind == kind {
		return obj, true
	} else {
		return nil, false
	}
}

func (s *Stack) GetTopValueType() int {
	obj := s.array[s.pos-1]
	return obj.Kind
}

func (s *Stack) Peek(kind int) (*StackValue, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	obj := s.array[s.pos-1]
	if obj.Kind == kind {
		return obj, true
	}
	return nil, false
}
