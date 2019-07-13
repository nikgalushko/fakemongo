package session

import (
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
		if e.Operator.Field == "a" {
			first = i
		} else {
			second = i
		}
	}

	assert.Len(t, expressions, 2)
	assert.Len(t, expressions[second].Operator.SubOperatorExpressions, 2)

	assert.Equal(t, "a", expressions[first].Operator.Field)
	assert.Equal(t, "$eq", expressions[first].Operator.Cmd)
	assert.Equal(t, 10, expressions[first].Operator.Value)

	assert.Equal(t, "$and", expressions[second].Operator.Cmd)
	assert.Nil(t, expressions[second].Operator.Value)

	subOperators := expressions[second].Operator.SubOperatorExpressions
	assert.Equal(t, "$exists", subOperators[first].Cmd)
	assert.Equal(t, true, subOperators[first].Value)
	assert.Equal(t, "c", subOperators[first].Field)

	assert.Equal(t, "$eq", subOperators[second].Cmd)
	assert.Equal(t, []string{"7", "8"}, subOperators[second].Value)
	assert.Equal(t, "d", subOperators[second].Field)

}

func TestSelectorParser_ParseQuery_Exists(t *testing.T) {
	p := SelectorParser{}

	expressions := p.ParseQuery(bson.M{"c": bson.M{"$exists": true}})

	assert.Len(t, expressions, 1)
	assert.Nil(t, expressions[0].Operator.SubOperatorExpressions)

	assert.Equal(t, "c", expressions[0].Operator.Field)
	assert.Equal(t, "$exists", expressions[0].Operator.Cmd)
	assert.Equal(t, true, expressions[0].Operator.Value)

}
