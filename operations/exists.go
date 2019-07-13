package operations

type Exists struct {
	DefaultOperator
}

func (e Exists) Do() bool {
	_, ok := e.Record[e.Field]
	return Eq{}.objectsAreEqual(e.Expected, ok)
}

func (e Exists) Name() string {
	return "$exists"
}
