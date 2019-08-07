package operations

import "github.com/jetuuuu/fakemongo/utils"

type Each struct {
	DefaultOperator
}

func (e Each) Do() interface{} { // todo if field is not array
	arr := utils.ToSlice(e.Record[e.Field]) // todo dot notation
	newValues := utils.ToSlice(e.Expected)
	arr = append(arr, newValues...)

	e.Record[e.Field] = arr

	return e.Record
}

func (e Each) Name() string {
	return "$each"
}
