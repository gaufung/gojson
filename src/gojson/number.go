package gojson

type Number struct {
	value interface{}
}

func NewNumber(v interface{}) *Number {
	return &Number{value: v}
}

func (n *Number) Float64Value() (float64, bool) {
	switch n.value.(type) {
	case float64:
		return n.value.(float64), true
	default:
		return 0.0, false
	}
}

func (n *Number) Int64Value() (int64, bool) {
	switch n.value.(type) {
	case int64:
		return n.value.(int64), true
	default:
		return 0, false
	}
}
