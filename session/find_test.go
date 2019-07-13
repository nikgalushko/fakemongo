package session

import (
	"fakemongo/collection"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testData = []collection.Record{
	{"f": 10, "arr2": "is not array"},
	{"c": true, "obj": bson.M{"f": 15.6}},
	{"f": 10, "arr": []interface{}{"t", "e", "s", "t"}},
	{"e": 4},
	{"f": 12, "obj": bson.M{"e": []interface{}{1, 2, 3}, "f": 18.2}, "e": 5},
}

func TestFinder_One_SimpleFields(t *testing.T) {
	f := NewFinder(bson.M{"e": 5, "f": 12}, testData)
	m := make(bson.M)
	err := f.One(&m)
	assert.NoError(t, err)
	assert.Equal(t, testData[4], collection.Record(m))
}

func TestFinder_One_OperatorAnd(t *testing.T) {
	f := NewFinder(bson.M{"$and": []bson.M{{"f": bson.M{"$eq": 10}}, {"arr": bson.M{"$exists": true}}}}, testData)
	m := make(bson.M)
	err := f.One(&m)
	assert.NoError(t, err)
	assert.Equal(t, testData[2], collection.Record(m))
}

func TestFinder_Select(t *testing.T) {
	f := NewFinder(bson.M{"e": 5}, testData).Select(bson.M{"obj": 1})
	m := make(bson.M)
	err := f.One(&m)
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"obj": testData[4]["obj"]}, (m))
}
