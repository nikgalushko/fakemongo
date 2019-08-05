package session

import (
	"fakemongo/collection"
	"fakemongo/utils"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testData = []collection.Record{
	{"f": 10, "arr2": "is not array"},
	{"c": true, "obj": bson.M{"f": 15.6}},
	{"f": 10, "arr": []interface{}{"t", "e", "s", "t"}},
	{"e": 4},
	{"obj": bson.M{"e": []interface{}{1, 2, 3}, "f": 18.2}, "e": 5, "f": 12},
	{"e": 5, "arr": []interface{}{bson.M{"price": 12}, bson.M{"price": 14}}, "shop": "#1"},
	{"e": 7, "arr": []interface{}{bson.M{"price": 45, "g": 14}, bson.M{"price": 36, "g": 16}}, "shop": "#2"},
	{"h": 19324, "arr": []interface{}{9, 4, 5, 1}, "car": "#3"},
}

func TestFinder_One_SimpleFields(t *testing.T) {
	f := NewFinder(bson.M{"e": 5, "f": 12}, cursor())
	m := make(bson.M)
	err := f.One(&m)
	assert.NoError(t, err)
	assert.Equal(t, testData[4], collection.Record(m))
}

func TestFinder_One_OperatorAnd(t *testing.T) {
	f := NewFinder(bson.M{"$and": []bson.M{{"f": bson.M{"$eq": 10}}, {"arr": bson.M{"$exists": true}}}}, cursor())
	m := make(bson.M)
	err := f.One(&m)
	assert.NoError(t, err)
	assert.Equal(t, testData[2], collection.Record(m))
}

func TestFinder_Select(t *testing.T) {
	f := NewFinder(bson.M{"e": 5}, cursor()).Select(bson.M{"obj": 1})
	m := make(bson.M)
	err := f.One(&m)
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"obj": testData[4]["obj"]}, m)
}

func TestFinder_One_ElemMatch_FlatArray(t *testing.T) {
	f := NewFinder(bson.M{"arr": bson.M{"$elemMatch": bson.M{"$eq": "s"}}}, cursor())
	m := make(bson.M)
	err := f.One(&m)
	assert.NoError(t, err)
	assert.Equal(t, testData[2], collection.Record(m))
}

func TestFinder_One_ElemMatch_ArrayOfObjects(t *testing.T) {
	f := NewFinder(bson.M{"arr": bson.M{"$elemMatch": bson.M{"g": bson.M{"$exists": true}, "price": bson.M{"$gte": 36}}}}, cursor())
	m := make(bson.M)
	err := f.One(&m)
	assert.NoError(t, err)
	assert.Equal(t, len(testData[6]), len(m))
	assert.Equal(t, testData[6]["e"], m["e"])
	assert.Equal(t, testData[6]["shop"], m["shop"])
	assert.Equal(t, utils.ToSlice(testData[6]["arr"]), utils.ToSlice(m["arr"]))
}

func TestFinder_One_DotNotation(t *testing.T) {
	f := NewFinder(bson.M{"obj.f": bson.M{"$lt": 20, "$gt": 18}}, cursor())
	m := make(bson.M)
	err := f.One(&m)

	a := assert.New(t)
	a.NoError(err)
	a.Equal(testData[4], collection.Record(m))
}

func TestFinder_One_DotNotation_Array_Eq(t *testing.T) {
	f := NewFinder(bson.M{"arr.price": 12}, cursor())
	m := make(bson.M)
	err := f.One(&m)

	a := assert.New(t)
	a.NoError(err)
	a.Equal(testData[5], collection.Record(m))
}

func TestFinder_One_DotNotation_Array_Gt(t *testing.T) {
	f := NewFinder(bson.M{"arr.price": bson.M{"$gt": 40}}, cursor())
	m := make(bson.M)
	err := f.One(&m)

	a := assert.New(t)
	a.NoError(err)
	a.Equal(testData[6], collection.Record(m))
}

func TestFinder_One_In(t *testing.T) {
	f := NewFinder(bson.M{"car": bson.M{"$in": []string{"#1", "2", "#3"}}}, cursor())
	m := make(bson.M)
	err := f.One(&m)

	a := assert.New(t)
	a.NoError(err)
	a.Equal(testData[7], collection.Record(m))
}

func TestFinder_One_DotNotation_In(t *testing.T) {
	f := NewFinder(bson.M{"arr.price": bson.M{"$in": []int{40, 41, 44, 36}}}, cursor())
	m := make(bson.M)
	err := f.One(&m)

	a := assert.New(t)
	a.NoError(err)
	a.Equal(testData[6], collection.Record(m))
}

func cursor() *collection.Cursor {
	return collection.NewCursor(&testData)
}
