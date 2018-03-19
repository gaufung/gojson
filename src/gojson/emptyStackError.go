package gojson

type EmptyStackError struct {
	message string
}

func (e *EmptyStackError) Error() string {
	return e.message
}

func NewEmptyStackError(mess string) *EmptyStackError {
	return &EmptyStackError{message: mess}
}
