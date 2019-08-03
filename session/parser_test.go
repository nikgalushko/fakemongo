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

	assert.Len(t, expressions, 2)
	assert.Len(t, expressions[second].(operations.OperatorExpression).SubOperatorExpressions, 2)

	firstExpression := expressions[first].(operations.OperatorExpression)
	assert.Equal(t, "a", firstExpression.Field)
	assert.Equal(t, "$eq", firstExpression.Cmd)
	assert.Equal(t, 10, firstExpression.Value)

	secondExpression := expressions[second].(operations.OperatorExpression)
	assert.Equal(t, "$and", secondExpression.Cmd)
	assert.Nil(t, secondExpression.Value)

	subOperators := secondExpression.SubOperatorExpressions
	firstSubOperation := subOperators[first].(operations.OperatorExpression)
	assert.Equal(t, "$exists", firstSubOperation.Cmd)
	assert.Equal(t, true, firstSubOperation.Value)
	assert.Equal(t, "c", firstSubOperation.Field)

	secondSubOperation := subOperators[second].(operations.OperatorExpression)
	assert.Equal(t, "$eq", secondSubOperation.Cmd)
	assert.Equal(t, []string{"7", "8"}, secondSubOperation.Value)
	assert.Equal(t, "d", secondSubOperation.Field)

}

func TestSelectorParser_ParseQuery_Exists(t *testing.T) {
	p := SelectorParser{}

	expressions := p.ParseQuery(bson.M{"c": bson.M{"$exists": true}})

	assert.Len(t, expressions, 1)
	assert.Nil(t, expressions[0].(operations.OperatorExpression).SubOperatorExpressions)

	firstExpression := expressions[0].(operations.OperatorExpression)
	assert.Equal(t, "c", firstExpression.Field)
	assert.Equal(t, "$exists", firstExpression.Cmd)
	assert.Equal(t, true, firstExpression.Value)
}

func TestSelectorParser_ParseQuery_ElemMatch(t *testing.T) {
	p := SelectorParser{}

	expressions := p.ParseQuery(bson.M{"arr": bson.M{"$elemMatch": bson.M{"a": 2, "d": bson.M{"$gte": 14}}}})

	assert.Len(t, expressions, 1)

	op := expressions[0].(operations.OperatorExpression)
	assert.Equal(t, "arr", op.Field)
	assert.Equal(t, "$elemMatch", op.Cmd)
	assert.Len(t, op.SubOperatorExpressions, 2)

	first := 0
	second := 1

	for i, o := range op.SubOperatorExpressions {
		if o.(operations.OperatorExpression).Field == "a" {
			first = i
		} else {
			second = i
		}
	}

	assert.Equal(t, "a", op.SubOperatorExpressions[first].(operations.OperatorExpression).Field)
	assert.Equal(t, "$eq", op.SubOperatorExpressions[first].(operations.OperatorExpression).Cmd)
	assert.Equal(t, 2, op.SubOperatorExpressions[first].(operations.OperatorExpression).Value)

	assert.Equal(t, "d", op.SubOperatorExpressions[second].(operations.OperatorExpression).Field)
	assert.Equal(t, "$gte", op.SubOperatorExpressions[second].(operations.OperatorExpression).Cmd)
	assert.Equal(t, 14, op.SubOperatorExpressions[second].(operations.OperatorExpression).Value)
}
