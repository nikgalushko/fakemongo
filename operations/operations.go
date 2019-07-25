package operations

import "fakemongo/collection"

type Expression struct {
	Operator OperatorExpression
}

type OperatorExpression struct {
	Cmd                    string
	Value                  interface{}
	Field                  string
	SubOperatorExpressions []OperatorExpression
}

func (o OperatorExpression) Match(record collection.Record) bool {
	cmd := NewOperator(o.Cmd)

	cmd.SetExpectedValue(o.Value)
	cmd.SetRetrievalField(o.Field)
	cmd.SetArgs(o.SubOperatorExpressions)
	cmd.SetRecord(record)

	return cmd.Do().(bool)
}

func (o OperatorExpression) Update(record collection.Record) collection.Record {
	cmd := NewOperator(o.Cmd)

	cmd.SetExpectedValue(o.Value)
	cmd.SetRetrievalField(o.Field)
	cmd.SetArgs(o.SubOperatorExpressions)
	cmd.SetRecord(record)

	return cmd.Do().(collection.Record)
}

type Operator interface {
	Do() interface{}
	SetExpectedValue(interface{})
	SetRetrievalField(string)
	SetArgs([]OperatorExpression)
	SetRecord(collection.Record)
	Name() string
}

type DefaultOperator struct {
	Expected               interface{}
	Field                  string
	SubOperatorExpressions []OperatorExpression
	Record                 collection.Record
}

func (d *DefaultOperator) SetExpectedValue(v interface{}) {
	d.Expected = v
}

func (d *DefaultOperator) SetRetrievalField(f string) {
	d.Field = f
}

func (d *DefaultOperator) SetArgs(ops []OperatorExpression) {
	d.SubOperatorExpressions = ops
}

func (d *DefaultOperator) SetRecord(r collection.Record) {
	d.Record = r
}

func NewOperator(cmd string) Operator {
	switch cmd {
	case "$eq":
		return &Eq{}
	case "$and":
		return &And{}
	case "$exists":
		return &Exists{}
	case "$elemMatch":
		return &ElemMatch{}
	case "$gt":
		return &Gt{}
	case "$gte":
		return &Gte{}
	case "$lt":
		return &Lt{}
	case "$lte":
		return &Lte{}
	case "$set":
		return &Set{}
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
		return &Eq{}
	}
}
