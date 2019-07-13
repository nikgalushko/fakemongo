package operations

type Exists struct{}

func (e Exists) Match(expected, actual interface{}) (bool, error) {
	panic("$exists is unimplemented")
	return false, nil
}

func (e Exists) Name() string {
	return "$exists"
}
