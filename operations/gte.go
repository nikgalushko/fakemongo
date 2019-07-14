package operations

type Gte struct {
	DefaultOperator
}

func (g Gte) Do() bool {
	actual := g.Record[g.Field]
	cmp, ok := compare(actual, g.Expected)
	return ok && (cmp == greaterThan || cmp == equal)
}

func (g Gte) Name() string {
	return "$gte"
}
