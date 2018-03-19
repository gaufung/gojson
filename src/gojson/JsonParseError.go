package gojson

type JsonParseError struct {
	errorIndex int
	message string
}

func (r JsonParseError) Error() string{
	return r.message
}

func (r JsonParseError) GetErrorIndex() int{
	return r.errorIndex
}

func NewJsonParserError(message string, index int) *JsonParseError{
	return  &JsonParseError{message:message, errorIndex:index}
}
