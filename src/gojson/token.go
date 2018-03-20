package gojson

type Token int

const (
	END_DOCUMENT = iota

	START_OBJECT

	END_OBJECT

	START_ARRAY

	END_ARRAY

	COLON_SEPERATOR

	COMA_SEPERATOR

	STRING

	BOOLEAN

	NUMBER

	NULL
)
