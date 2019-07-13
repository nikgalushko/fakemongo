package operations

import (
	"bytes"
	"reflect"
)

type Eq struct {
	DefaultOperator
}

func (e Eq) Do() bool {
	return e.objectsAreEqual(e.Expected, e.Record[e.Field])
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
