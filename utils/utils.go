package utils

import (
	"reflect"
)

func ToSlice(s interface{}) []interface{} {
	v := reflect.ValueOf(s)
	slice := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		slice[i] = v.Index(i).Interface()
	}
	return slice
}
