package operations

type Lt struct {
	DefaultOperator
}

func (l Lt) Do() interface{} {
	actual, ok := l.Record.GetByField(l.Field)
	if !ok {
		return false
	}

	return applyTo(func(expected, actual interface{}) bool {
		cmp, ok := compare(actual, expected)
		return ok && cmp == lessThan
	}, l.Expected, actual)

}

func (l Lt) Name() string {
	return "$lt"
}
