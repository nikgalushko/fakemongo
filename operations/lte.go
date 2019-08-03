package operations

type Lte struct {
	DefaultOperator
}

func (l Lte) Do() interface{} {
	actual, ok := l.Record.GetByField(l.Field)
	if !ok {
		return false
	}
	cmp, ok := compare(actual, l.Expected)
	return ok && (cmp == lessThan || cmp == equal)
}

func (l Lte) Name() string {
	return "$lte"
}
