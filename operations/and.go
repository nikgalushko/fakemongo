package operations

type And struct {
	DefaultOperator
}

func (a And) Do() interface{} {
	for _, e := range a.SubOperatorExpressions {
		if !e.Match(a.Record) {
			return false
		}
	}

	return true
}

func (a And) Name() string {
	return "$and"
}
