package compare

import (
	"reflect"
	"time"
)

func CompareTo(a, b interface{}) string {
	//aType := reflect.TypeOf(a)
	//bType := reflect.TypeOf(b)
	//fmt.Println(aType.Name(), bType.Name())
	/*if !areTypeSame(aType, bType) {
		return "err"
	}*/

	switch a.(type) {
	case int:
		b2, ok := b.(int)
		if !ok {
			return "err"
		}

		a2 := a.(int)
		if a2 > b2 {
			return "gt"
		} else if a2 < b2 {
			return "lt"
		} else if a2 >= b2 {
			return "gte"
		} else if a2 <= b2 {
			return "lte"
		} else {
			return "eq"
		}
	case float64, float32:
		b2, ok := b.(float64)
		if !ok {
			return "err"
		}

		a2 := a.(float64)
		if a2 > b2 {
			return "gt"
		} else if a2 < b2 {
			return "lt"
		} else if a2 >= b2 {
			return "gte"
		} else if a2 <= b2 {
			return "lte"
		} else {
			return "eq"
		}
	case time.Time:
		b2, ok := b.(time.Time)
		if !ok {
			return "err"
		}

		a2 := a.(time.Time)
		if a2.After(b2) {
			return "gt"
		} else if a2.Before(b2) {
			return "lt"
		} else if a2.After(b2) || a2.Equal(b2) {
			return "gte"
		} else if a2.Before(b2) || a2.Equal(b2) {
			return "lte"
		} else {
			return "eq"
		}
	default:
		return "err"
	}
}

func areTypeSame(a, b reflect.Type) bool {
	// todo compare types
	return a.Name() == b.Name()
}

func inSlice(value string, arr []string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}
