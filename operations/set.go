package operations

type Set struct {
	DefaultOperator
	insertable bool
}

func (s Set) Do() interface{} {
	s.Record[s.Field] = s.Expected // todo dot notation
	return s.Record
}

func (s Set) Name() string {
	return "$set"
}

func (s Set) onlyForInsert() bool {
	return s.insertable
}
