package operations

import (
	"bytes"
	"reflect"
)

type Eq struct {
	DefaultOperator
}

func (e Eq) Do() interface{} {
	actual, ok := e.Record.GetByField(e.Field)
	if !ok {
		return false
	}

	return applyTo(e.objectsAreEqual, e.Expected, actual)
}

func (e Eq) objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

func (e Eq) Name() string {
	return "$eq"
}

func isSlice(i interface{}) bool {
	if i == nil {
		return false
	}

	t := reflect.TypeOf(i).Kind()
	return t == reflect.Slice || t == reflect.Array
}

func applyTo(f func(interface{}, interface{}) bool, args ...interface{}) bool {
	expected := args[0]
	actual := args[1]

	var value []interface{}
	if isSlice(actual) {
		value = actual.([]interface{})
	} else {
		value = append(value, actual)
	}

	for _, v := range value {
		if f(expected, v) {
			return true
		}
	}

	return false
}
