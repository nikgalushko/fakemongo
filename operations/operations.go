package operations

type Operator interface {
	Match(expected, actual interface{}) (bool, error)
	Name() string
}

func New(cmd string) Operator {
	switch cmd {
	case "$eq":
		return Eq{}
	case "$and":
		return And{}
	case "$exists":
		return Exists{}
	/*case "$or":
		return Or{}
	case "$nin":
		return Nin{}
	case "$in":
		return In{}
	case "$ne":
		return Ne{}
	case "$exists":
		return Exists{}
	case "$not":
		return Not{}
	case "$size":
		return Size{}*/
	default:
		return Eq{}
	}
}
