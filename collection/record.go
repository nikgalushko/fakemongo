package collection

import "github.com/globalsign/mgo/bson"

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
