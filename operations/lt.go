package operations

type Lt struct {
	DefaultOperator
}

func (l Lt) Do() interface{} {
	actual, ok := l.Record.GetByField(l.Field)
	if !ok {
		return false
	}

	cmp, ok := compare(actual, l.Expected)
	return ok && cmp == lessThan
}

func (l Lt) Name() string {
	return "$lt"
}
