package operations

type Gt struct {
	DefaultOperator
}

func (g Gt) Do() bool {
	actual := g.Record[g.Field]
	cmp, ok := compare(actual, g.Expected)
	return ok && cmp == greaterThan
}

func (g Gt) Name() string {
	return "$gt"
}
