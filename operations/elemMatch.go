package operations

import (
	"fakemongo/collection"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
)

type ElemMatch struct {
	DefaultOperator
}

func (e ElemMatch) Name() string {
	return "$elemMatch"
}

func (e ElemMatch) Do() interface{} {
	value := e.Record[e.Field]
	if value == nil {
		return false
	}
	rt := reflect.TypeOf(value).Kind()
	if rt != reflect.Slice && rt != reflect.Array {
		panic(fmt.Sprintf("$elemMatch doesn't support %v", rt))
	}

	// todo value can be not just []bson.M
	arr := reflect.ValueOf(value)
	for i := 0; i < arr.Len(); i++ {
		v := arr.Index(i).Interface()
		matched := true
		for _, e2 := range e.SubOperatorExpressions {
			var r collection.Record
			switch v.(type) {
			case bson.M:
				r = collection.Record(v.(bson.M))
			default:
				r = collection.Record(bson.M{"": v})
			}
			if !e2.Match(r) {
				matched = false
				break
			}
		}

		if matched {
			return true
		}
	}

	return false
}
