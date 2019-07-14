package operations

type cmp int

const (
	equal cmp = iota
	lessThanOrEqual
	lessThan
	greaterThanOrEqual
	greaterThan
	illegal
)

func compare(a, b interface{}) (cmp, bool) {
	switch a.(type) {
	case int, int8, int16, int32, float32, float64:
		actual, ok := toFloat64(a)
		if !ok {
			return illegal, false
		}

		expected, ok := toFloat64(b)
		if !ok {
			return illegal, false
		}

		if actual < expected {
			return lessThan, true
		} else if actual > expected {
			return greaterThan, true
		} else {
			return equal, true
		}
	}

	return illegal, false
}

func toFloat64(i interface{}) (float64, bool) {
	switch i.(type) {
	case int:
		return float64(i.(int)), true
	case int8:
		return float64(i.(int8)), true
	case int16:
		return float64(i.(int16)), true
	case int32:
		return float64(i.(int32)), true
	case int64:
		return 0, false
	case float32:
		return float64(i.(float32)), true
	case float64:
		return i.(float64), true
	default:
		return 0, false
	}
}
