package operations

type And struct{}

func (a And) Match(expected, actual interface{}) (bool, error) {
	panic("and is unimplemented")
	return false, nil
}

func (a And) Name() string {
	return "$and"
}
