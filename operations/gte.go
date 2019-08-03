package operations

type Gte struct {
	DefaultOperator
}

func (g Gte) Do() interface{} {
	actual, ok := g.Record.GetByField(g.Field)
	if !ok {
		return false
	}
	cmp, ok := compare(actual, g.Expected)
	return ok && (cmp == greaterThan || cmp == equal)
}

func (g Gte) Name() string {
	return "$gte"
}
