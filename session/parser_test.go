package session

import (
	"fakemongo/operations"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectorParser_ParseQuery(t *testing.T) {
	p := SelectorParser{}

	and := []bson.M{
		{"c": bson.M{"$exists": true}},
		{"d": []string{"7", "8"}},
	}
	expressions := p.ParseQuery(bson.M{"a": 10, "$and": and})

	first := 0
	second := 1
	for i, e := range expressions {
		if e.(operations.OperatorExpression).Field == "a" {
			first = i
		} else {
			second = i
		}
	}

	a := assert.New(t)
	a.Len(expressions, 2)
	a.Len(expressions[second].(operations.OperatorExpression).SubOperatorExpressions, 2)

	firstExpression := expressions[first].(operations.OperatorExpression)
	a.Equal("a", firstExpression.Field)
	a.Equal("$eq", firstExpression.Cmd)
	a.Equal(10, firstExpression.Value)

	secondExpression := expressions[second].(operations.OperatorExpression)
	a.Equal("$and", secondExpression.Cmd)
	a.Nil(secondExpression.Value)

	subOperators := secondExpression.SubOperatorExpressions
	firstSubOperation := subOperators[first].(operations.OperatorExpression)
	a.Equal("$exists", firstSubOperation.Cmd)
	a.Equal(true, firstSubOperation.Value)
	a.Equal("c", firstSubOperation.Field)

	secondSubOperation := subOperators[second].(operations.OperatorExpression)
	a.Equal("$eq", secondSubOperation.Cmd)
	a.Equal([]string{"7", "8"}, secondSubOperation.Value)
	a.Equal("d", secondSubOperation.Field)

}

func TestSelectorParser_ParseQuery_Exists(t *testing.T) {
	p := SelectorParser{}

	expressions := p.ParseQuery(bson.M{"c": bson.M{"$exists": true}})
	a := assert.New(t)

	a.Len(expressions, 1)
	a.Nil(expressions[0].(operations.OperatorExpression).SubOperatorExpressions)

	firstExpression := expressions[0].(operations.OperatorExpression)
	a.Equal("c", firstExpression.Field)
	a.Equal("$exists", firstExpression.Cmd)
	a.Equal(true, firstExpression.Value)
}

func TestSelectorParser_ParseQuery_ElemMatch(t *testing.T) {
	p := SelectorParser{}

	expressions := p.ParseQuery(bson.M{"arr": bson.M{"$elemMatch": bson.M{"a": 2, "d": bson.M{"$gte": 14}}}})
	a := assert.New(t)
	a.Len(expressions, 1)

	op := expressions[0].(operations.OperatorExpression)
	a.Equal("arr", op.Field)
	a.Equal("$elemMatch", op.Cmd)
	a.Len(op.SubOperatorExpressions, 2)

	first := 0
	second := 1

	for i, o := range op.SubOperatorExpressions {
		if o.(operations.OperatorExpression).Field == "a" {
			first = i
		} else {
			second = i
		}
	}

	a.Equal("a", op.SubOperatorExpressions[first].(operations.OperatorExpression).Field)
	a.Equal("$eq", op.SubOperatorExpressions[first].(operations.OperatorExpression).Cmd)
	a.Equal(2, op.SubOperatorExpressions[first].(operations.OperatorExpression).Value)

	a.Equal("d", op.SubOperatorExpressions[second].(operations.OperatorExpression).Field)
	a.Equal("$gte", op.SubOperatorExpressions[second].(operations.OperatorExpression).Cmd)
	a.Equal(14, op.SubOperatorExpressions[second].(operations.OperatorExpression).Value)
}

func TestSelectorParser_ParseQuery_DotNotation(t *testing.T) {
	p := SelectorParser{}
	expressions := p.ParseQuery(bson.M{"obj1.sub_obj1.sub_obj2.field": 10})

	a := assert.New(t)
	e := expressions[0].(operations.OperatorExpression)

	a.Len(expressions, 1)
	a.Equal("obj1.sub_obj1.sub_obj2.field", e.Field)
	a.Equal(10, e.Value)
	a.Equal("$eq", e.Cmd)
}

func TestSelectorParser_ParseQuery_MultiCondition(t *testing.T) {
	p := SelectorParser{}

	expressions := p.ParseQuery(bson.M{"c": bson.M{"$gt": 10, "$lt": 20}})
	a := assert.New(t)

	a.Len(expressions, 2)

	gt := expressions[0].(operations.OperatorExpression)
	a.Equal("c", gt.Field)
	a.Equal("$gt", gt.Cmd)
	a.Equal(10, gt.Value)

	lt := expressions[1].(operations.OperatorExpression)
	a.Equal("c", lt.Field)
	a.Equal("$lt", lt.Cmd)
	a.Equal(20, lt.Value)
}

func TestUpdateParameterParser_ParseUpdate(t *testing.T) {
	p := UpdateParameterParser{}

	expressions := p.ParseUpdate(bson.M{"$push": bson.M{"scores": bson.M{"$each": []int{90, 92, 85}}}})
	a := assert.New(t)

	a.Len(expressions, 1)

	push := expressions[0]
	a.Equal("$push", push.Cmd)
	a.Equal("scores", push.Field)
	a.Nil(push.Value)
	a.Len(push.SubOperatorExpressions, 1)

	each := push.SubOperatorExpressions[0].(operations.OperatorExpression)
	a.Equal("$each", each.Cmd)
	a.Equal("scores", each.Field)
	a.Equal([]int{90, 92, 85}, each.Value)
}

func TestUpdateParameterParser_ParseUpdate_MultiCondition(t *testing.T) {
	p := UpdateParameterParser{}

	expressions := p.ParseUpdate(bson.M{"$push": bson.M{"scores": 12, "cars": "car#2"}})
	a := assert.New(t)

	a.Len(expressions, 2)

	scores := expressions[0]
	a.Equal("scores", scores.Field)
	a.Equal("$push", scores.Cmd)
	a.Equal(12, scores.Value)

	cars := expressions[1]
	a.Equal("cars", cars.Field)
	a.Equal("$push", cars.Cmd)
	a.Equal("car#2", cars.Value)
}
