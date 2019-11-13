package operations

import "github.com/jetuuuu/fakemongo/collection"

type Expression interface {
	Match(collection.Record) bool
	Update(collection.Record) collection.Record
}

type OperatorExpression struct {
	Cmd                    string
	Value                  interface{}
	Field                  string
	SubOperatorExpressions []Expression
}

func (o OperatorExpression) Match(record collection.Record) bool {
	cmd := NewOperator(o.Cmd, o.Value, o.Field, o.SubOperatorExpressions, record)

	return cmd.Do().(bool)
}

func (o OperatorExpression) Update(record collection.Record) collection.Record {
	cmd := NewOperator(o.Cmd, o.Value, o.Field, o.SubOperatorExpressions, record)

	return cmd.Do().(collection.Record)
}

type KindOperation uint8

const (
	Unknow KindOperation = iota
	Update
	Insert
)

func (o OperatorExpression) Upsert(record collection.Record, kind KindOperation) collection.Record {
	cmd := NewOperator(o.Cmd, o.Value, o.Field, o.SubOperatorExpressions, record)
	if kind == Update && cmd.onlyForInsert() {
		return record
	}
	return cmd.Do().(collection.Record)
}

type Operator interface {
	Do() interface{}
	Name() string
	onlyForInsert() bool
}

type DefaultOperator struct {
	Expected               interface{}
	Field                  string
	SubOperatorExpressions []Expression
	Record                 collection.Record
}

func (o DefaultOperator) onlyForInsert() bool {
	return false
}

func NewOperator(cmd string, v interface{}, f string, args []Expression, r collection.Record) Operator {
	d := DefaultOperator{Expected: v, Field: f, SubOperatorExpressions: args, Record: r}
	switch cmd {
	case "$eq":
		return &Eq{d}
	case "$and":
		return &And{d}
	case "$exists":
		return &Exists{d}
	case "$elemMatch":
		return &ElemMatch{d}
	case "$gt":
		return &Gt{d}
	case "$gte":
		return &Gte{d}
	case "$lt":
		return &Lt{d}
	case "$lte":
		return &Lte{d}
	case "$set":
		return &Set{d, false}
	case "$setOnInsert":
		return &Set{d, true}
	case "$in":
		return &In{d}
	case "$push":
		return &Push{d}
	case "$each":
		return &Each{d}
	/*case "$or":
		return Or{}
	case "$nin":
		return Nin{}
	case "$ne":
		return Ne{}
	case "$exists":
		return Exists{}
	case "$not":
		return Not{}
	case "$size":
		return Size{}*/
	default:
		return &Eq{d}
	}
}
