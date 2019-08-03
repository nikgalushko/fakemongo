package collection

import (
	"github.com/globalsign/mgo/bson"
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
			rec, ok = v.(bson.M)
			if !ok {
				return nil, false
			}
		} else {
			return v, ok
		}
	}

	return nil, false
}
