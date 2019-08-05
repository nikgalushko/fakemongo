package operations

type Gte struct {
	DefaultOperator
}

func (g Gte) Do() interface{} {
	actual, ok := g.Record.GetByField(g.Field)
	if !ok {
		return false
	}

	return applyTo(func(expected, actual interface{}) bool {
		cmp, ok := compare(actual, expected)
		return ok && (cmp == greaterThan || cmp == equal)
	}, g.Expected, actual)
}

func (g Gte) Name() string {
	return "$gte"
}
