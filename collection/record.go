package collection

import (
	"github.com/globalsign/mgo/bson"
	"github.com/jetuuuu/fakemongo/utils"
	"reflect"
	"strconv"
	"strings"
)

type Record bson.M

func (r Record) WithFields(fields bson.M) Record {
	if fields == nil {
		return r
	}

	ret := make(Record)
	for f, v := range fields {
		switch v.(type) {
		case int:
			if value, ok := r[f]; v.(int) == 1 && ok {
				ret[f] = value
			}
		}
	}

	return ret
}

func (r Record) GetByField(f string) (interface{}, bool) {
	fields := strings.Split(f, ".")
	rec := bson.M(r)

	for i := 0; i < len(fields); i++ {
		v, ok := rec[fields[i]]
		if !ok {
			return nil, false
		}

		if i < (len(fields) - 1) {
			switch v.(type) {
			case bson.M:
				rec = v.(bson.M)
			case []bson.M, []interface{}:
				var ret []interface{}
				arr := utils.ToSlice(v)
				if index, err := strconv.Atoi(fields[i+1]); err == nil {
					arr = []interface{}{arr[index]} // todo check array size
					i += 1
				}
				for _, el := range arr {
					var v interface{}
					if obj, ok := el.(bson.M); ok {
						v, _ = Record(obj).GetByField(strings.Join(fields[i+1:], "."))
					}
					if s, ok := toSlice(v); ok {
						ret = append(ret, s...)
					} else {
						ret = append(ret, v)
					}
				}

				return ret, true
			default:
				return nil, false
			}
		} else {
			return v, ok
		}
	}

	return nil, false
}

func toSlice(i interface{}) ([]interface{}, bool) {
	if i == nil {
		return nil, false
	}

	t := reflect.TypeOf(i).Kind()
	if t == reflect.Slice || t == reflect.Array {
		return i.([]interface{}), true
	}

	return nil, false
}
