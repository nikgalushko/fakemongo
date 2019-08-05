package operations

type Gt struct {
	DefaultOperator
}

func (g Gt) Do() interface{} {
	actual, ok := g.Record.GetByField(g.Field)
	if !ok {
		return false
	}

	return applyTo(func(expected, actual interface{}) bool {
		cmp, ok := compare(actual, expected)
		return ok && cmp == greaterThan
	}, g.Expected, actual)
}

func (g Gt) Name() string {
	return "$gt"
}
