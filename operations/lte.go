package operations

type Lte struct {
	DefaultOperator
}

func (l Lte) Do() interface{} {
	actual, ok := l.Record.GetByField(l.Field)
	if !ok {
		return false
	}

	return applyTo(func(expected, actual interface{}) bool {
		cmp, ok := compare(actual, expected)
		return ok && (cmp == lessThan || cmp == equal)
	}, l.Expected, actual)
}

func (l Lte) Name() string {
	return "$lte"
}
