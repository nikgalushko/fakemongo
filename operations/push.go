package operations

import "fakemongo/utils"

type Push struct {
	DefaultOperator
}

func (p Push) Do() interface{} {
	// todo dot notation
	r := p.Record
	if isSlice(p.Record[p.Field]) {

		arr := utils.ToSlice(r[p.Field])
		if len(p.SubOperatorExpressions) > 0 {
			for _, op := range p.SubOperatorExpressions {
				r = op.Update(r)
			}
		} else {
			r[p.Field] = append(arr, p.Expected)
		}
	}

	return r
}

func (p Push) Name() string {
	return "$push"
}
