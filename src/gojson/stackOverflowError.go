package gojson

type StackOverflowError struct {
	message string
}

func (s *StackOverflowError) Error() string {
	return s.message
}

func NewStackOverflowError(message string) *StackOverflowError {
	return &StackOverflowError{message: message}
}
