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

func (s *Stack) Push(obj *StackValue) error {
	if s.pos == s.Size {
		return NewStackOverflowError("Maximum depth reached when parse JSON string")
	}
	s.array[s.pos] = obj
	s.pos++
	return nil
}

func (s *Stack) Pop() (*StackValue, error) {
	if s.IsEmpty() {
		return nil, NewEmptyStackError("The stack is empty")
	}
	s.pos--
	return s.array[s.pos], nil
}

func (s *Stack) PopKind(kind int) (*StackValue, error) {
	if s.IsEmpty() {
		return nil, NewEmptyStackError("The stack is empty")
	}
	s.pos--
	obj := s.array[s.pos]
	if obj.Kind == kind {
		return obj, nil
	} else {
		return nil, NewJsonParserError("unmatched object or empty", 1)
	}
}

func (s *Stack) GetTopValueType() int {
	obj := s.array[s.pos-1]
	return obj.Kind
}

func (s *Stack) Peek(kind int) (*StackValue, error) {
	if s.IsEmpty() {
		return nil, NewEmptyStackError("The stack is empty")
	}
	obj := s.array[s.pos-1]
	if obj.Kind == kind {
		return obj, nil
	}
	return nil, NewJsonParserError("Unmatched object or array.", 0)
}
