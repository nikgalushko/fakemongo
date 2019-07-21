package session

import (
	"fakemongo/collection"
	"fakemongo/operations"
	"github.com/globalsign/mgo/bson"
)

type SelectorParser struct{}

func (p SelectorParser) ParseQuery(query bson.M) []operations.Expression {
	var result []operations.Expression
	for k, expression := range query {
		key := collection.Key(k)
		switch key.Type() {
		case collection.Cmd:
			var e bson.M
			switch expression.(type) {
			case bson.M:
				e = expression.(bson.M)
			case []bson.M:
				e = bson.M{k: expression}
			default:
				e = bson.M{k: expression}
			}
			result = append(result, operations.Expression{Operator: p.ParseOperatorExpression(e)})
		case collection.DotNotation:
		case collection.Literal:
			e := operations.Expression{}
			e.Operator = p.ParseLiteralSubQuery(expression)
			e.Operator.Field = k
			result = append(result, e)
		default:
			panic(key)
		}
	}

	return result
}

func (p SelectorParser) ParseLiteralSubQuery(query interface{}) operations.OperatorExpression {
	switch query.(type) {
	case bson.M:
		for k, value := range query.(bson.M) {
			if collection.Key(k).IsCmd() {
				return p.ParseOperatorExpression(bson.M{k: value})
			}
		}
	}

	return operations.OperatorExpression{
		Value: query,
		Cmd:   "$eq",
	}
}

func (p SelectorParser) ParseOperatorExpression(query bson.M) operations.OperatorExpression {
	for cmd, value := range query {
		e := operations.OperatorExpression{Cmd: cmd}
		switch cmd {
		case "$eq", "$exists", "$gt", "$gte", "$lt", "$lte":
			e.Value = value
		case "$and":
			// todo to slice
			slice := value.([]bson.M)
			for _, s := range slice {
				e.SubOperatorExpressions = append(e.SubOperatorExpressions, p.ParseQuery(s)[0].Operator)
			}
			return e
		case "$elemMatch":
			v := value.(bson.M)
			parsed := p.ParseQuery(v)
			for _, expression := range parsed {
				e.SubOperatorExpressions = append(e.SubOperatorExpressions, expression.Operator)
			}
		default:
			panic(cmd)
		}

		return e
	}

	return operations.OperatorExpression{}
}
