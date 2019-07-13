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

/*
type OperatorExpression struct {
	Cmd                    operations.Operator
	Value                  interface{}
	Field                  string
	SubOperatorExpressions []OperatorExpression
}

*/

/*func (r Record) MatchAll(expressions []Expression) bool {
	for _, e := range expressions {
		op := e.Operator
		ret, _ := op.Cmd.Match(op.Value, r[op.Field])
		if !ret {
			return false
		}
	}

	return true
}*/
