package operations

type Lt struct {
	DefaultOperator
}

func (l Lt) Do() bool {
	actual := l.Record[l.Field]
	cmp, ok := compare(actual, l.Expected)
	return ok && cmp == lessThan
}

func (l Lt) Name() string {
	return "$lt"
}
