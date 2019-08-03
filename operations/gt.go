package operations

type Gt struct {
	DefaultOperator
}

func (g Gt) Do() interface{} {
	actual, ok := g.Record.GetByField(g.Field)
	if !ok {
		return false
	}
	cmp, ok := compare(actual, g.Expected)
	return ok && cmp == greaterThan
}

func (g Gt) Name() string {
	return "$gt"
}
