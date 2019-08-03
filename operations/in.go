package operations

import (
	"reflect"
	"strings"
)

type In struct {
	DefaultOperator
}

func (i In) Name() string {
	return "$in"
}

func (i In) Do() interface{} {
	actual, ok := i.Record.GetByField(i.Field)
	if !ok {
		return false
	}

	ok, found := i.includeElement(i.Expected, actual)
	return ok && found
}

// todo rewrite
func (i In) includeElement(list interface{}, element interface{}) (ok, found bool) {
	eq := Eq{}
	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if eq.objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if eq.objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false

}
