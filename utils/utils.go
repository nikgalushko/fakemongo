package utils

import (
	"github.com/globalsign/mgo/bson"
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

func ToBsonM(i interface{}) (bson.M, error) {
	data, err := bson.Marshal(i)
	if err != nil {
		return nil, err
	}

	var ret bson.M
	err = bson.Unmarshal(data, &ret)

	return ret, err
}
