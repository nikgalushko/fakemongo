package operations

type And struct {
	DefaultOperator
}

func (a And) Do() bool {
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
