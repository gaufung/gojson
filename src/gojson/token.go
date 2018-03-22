package gojson

type Token int

const (
	END_DOCUMENT = iota // 0

	START_OBJECT // 1

	END_OBJECT // 2

	START_ARRAY // 3

	END_ARRAY // 4

	COLON_SEPERATOR // 5

	COMA_SEPERATOR // 6

	STRING // 7

	BOOLEAN // 8

	NUMBER // 9

	NULL // 10
)
