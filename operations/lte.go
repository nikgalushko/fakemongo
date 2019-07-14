package operations

type Lte struct {
	DefaultOperator
}

func (l Lte) Do() bool {
	actual := l.Record[l.Field]
	cmp, ok := compare(actual, l.Expected)
	return ok && (cmp == lessThan || cmp == equal)
}

func (l Lte) Name() string {
	return "$lte"
}
