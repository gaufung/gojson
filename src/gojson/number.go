package gojson

import "errors"

type Number struct {
	value interface{}
}

func NewNumber(v interface{}) *Number {
	return &Number{value: v}
}

func (n *Number) float64Value() (float64, error) {
	switch n.value.(type) {
	case float64:
		return n.value.(float64), nil
	default:
		return 0.0, errors.New("Not a float")
	}
}

func (n *Number) int64Value() (int64, error) {
	switch n.value.(type) {
	case int64:
		return n.value.(int64), nil
	default:
		return 0, errors.New("Not a int")
	}
}
