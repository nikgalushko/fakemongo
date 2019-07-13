package session

import (
	"fakemongo/collection"
	"fakemongo/operations"
	"github.com/globalsign/mgo/bson"
)

type SelectorParser struct{}

func (p SelectorParser) ParseQuery(query bson.M) []collection.Expression {
	var result []collection.Expression
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
				panic(expression)
			}
			result = append(result, collection.Expression{Operator: p.ParseOperatorExpression(e)})
		case collection.DotNotation:
		case collection.Literal:
			e := collection.Expression{}
			e.Operator = p.ParseLiteralSubQuery(expression)
			e.Operator.Field = k
			result = append(result, e)
		default:
			panic(key)
		}
	}

	return result
}

func (p SelectorParser) ParseLiteralSubQuery(query interface{}) collection.OperatorExpression {
	switch query.(type) {
	case bson.M:
		for k, value := range query.(bson.M) {
			if collection.Key(k).IsCmd() {
				return p.ParseOperatorExpression(bson.M{k: value})
			}
		}
	}

	return collection.OperatorExpression{
		Value: query,
		Cmd:   operations.Eq{},
	}
}

func (p SelectorParser) ParseOperatorExpression(query bson.M) collection.OperatorExpression {
	for cmd, value := range query {
		e := collection.OperatorExpression{Cmd: operations.New(cmd)}
		switch e.Cmd.(type) {
		case operations.Eq, operations.Exists:
			e.Value = value
		case operations.And:
			// todo to slice
			slice := value.([]bson.M)
			for _, s := range slice {
				e.SubOperatorExpressions = append(e.SubOperatorExpressions, p.ParseQuery(s)[0].Operator)
			}
			return e
		default:
			panic(cmd)
		}

		return e
	}

	return collection.OperatorExpression{}
}
