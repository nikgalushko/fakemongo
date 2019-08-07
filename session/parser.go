package session

import (
	"github.com/globalsign/mgo/bson"
	"github.com/jetuuuu/fakemongo/collection"
	"github.com/jetuuuu/fakemongo/operations"
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
			default:
				e = bson.M{k: expression}
			}
			result = append(result, p.ParseOperatorExpression(e))
		case collection.DotNotation, collection.Literal:
			for _, op := range p.ParseLiteralSubQuery(expression) {
				op.Field = k
				result = append(result, op)
			}
		default:
			panic(key)
		}
	}

	return result
}

func (p SelectorParser) ParseLiteralSubQuery(query interface{}) []operations.OperatorExpression {
	switch query.(type) {
	case bson.M:
		var ret []operations.OperatorExpression
		for k, value := range query.(bson.M) {
			if collection.Key(k).IsCmd() {
				ret = append(ret, p.ParseOperatorExpression(bson.M{k: value}))
			}
		}
		return ret
	}

	return []operations.OperatorExpression{
		{
			Value: query,
			Cmd:   "$eq",
		},
	}

}

func (p SelectorParser) ParseOperatorExpression(query bson.M) operations.OperatorExpression {
	for cmd, value := range query {
		e := operations.OperatorExpression{Cmd: cmd}
		switch cmd {
		case "$eq", "$exists", "$gt", "$gte", "$lt", "$lte", "$in":
			e.Value = value
		case "$and":
			// todo to slice
			slice := value.([]bson.M)
			for _, s := range slice {
				e.SubOperatorExpressions = append(e.SubOperatorExpressions, p.ParseQuery(s)[0])
			}
			return e
		case "$elemMatch":
			v := value.(bson.M)
			parsed := p.ParseQuery(v)
			for _, expression := range parsed {
				e.SubOperatorExpressions = append(e.SubOperatorExpressions, expression)
			}
		default:
			panic(cmd)
		}

		return e
	}

	return operations.OperatorExpression{}
}

type UpdateParameterParser struct{}

type UpdateOperator struct {
	Cmd        string
	Operations []operations.OperatorExpression
}

func (p UpdateParameterParser) ParseUpdate(update bson.M) []operations.OperatorExpression {
	var result []operations.OperatorExpression
	for k, v := range update {
		key := collection.Key(k)
		if key.IsCmd() {
			switch k {
			case "$set", "$push":
				fields := v.(bson.M)
				var ops []operations.OperatorExpression
				for f, v := range fields {
					op := operations.OperatorExpression{Cmd: k}
					op.Field = f

					if sub := p.ParseSubCmd(v); k == "$push" && len(sub) > 0 {
						for i := 0; i < len(sub); i++ {
							sub[i].Field = f
							op.SubOperatorExpressions = append(op.SubOperatorExpressions, sub[i])
						}
					} else {
						op.Value = v
					}

					ops = append(ops, op)
				}
				result = append(result, ops...)
			default:
				panic("unimplemented")
			}
		} else {
			panic("unimplemented")
		}
	}

	return result
}

func (p UpdateParameterParser) ParseSubCmd(v interface{}) []operations.OperatorExpression {
	var result []operations.OperatorExpression
	update, ok := v.(bson.M)
	if !ok {
		return nil
	}
	for k, v := range update {
		switch k {
		case "$each":
			result = append(result, operations.OperatorExpression{
				Cmd:   "$each",
				Value: v,
			})
		}
	}

	return result
}
