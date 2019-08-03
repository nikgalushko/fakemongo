package operations

type Exists struct {
	DefaultOperator
}

func (e Exists) Do() interface{} {
	_, ok := e.Record.GetByField(e.Field)
	return Eq{}.objectsAreEqual(e.Expected, ok)
}

func (e Exists) Name() string {
	return "$exists"
}
