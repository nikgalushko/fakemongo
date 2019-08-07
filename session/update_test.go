package session

import (
	"fakemongo/collection"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDataUpd = []collection.Record{
	{"f": 10, "arr2": "is not array"},
	{"c": true, "obj": bson.M{"f": 15.6}},
	{"f": 10, "arr": []interface{}{"t", "e", "s", "t"}, "arr2": []interface{}{5}},
	{"e": 4},
	{"f": 12, "obj": []interface{}{1, 2, 3}, "e": 5},
	{"test": "#dot_notation", "obj": bson.M{"arr": []interface{}{12}}},
}

func TestUpdater_Update(t *testing.T) {
	query := bson.M{"e": 4}
	update := bson.M{"$set": bson.M{"e2": 5, "g3": "test"}}
	u := Updater{c: collection.NewCursor(&testDataUpd)}

	err := u.Update(query, update)
	assert.NoError(t, err)
	assert.Equal(t, collection.Record{"e": 4, "e2": 5, "g3": "test"}, testDataUpd[3])
}

func TestUpdater_Update_Push(t *testing.T) {
	query := bson.M{"f": 10, "arr": bson.M{"$exists": true}}
	update := bson.M{"$push": bson.M{"arr": "2", "arr2": "not_a_number"}}
	u := Updater{c: collection.NewCursor(&testDataUpd)}

	err := u.Update(query, update)
	assert.NoError(t, err)
	assert.Equal(t, collection.Record{"f": 10, "arr": []interface{}{"t", "e", "s", "t", "2"}, "arr2": []interface{}{5, "not_a_number"}}, testDataUpd[2])
}

func TestUpdater_Update_Push_Each(t *testing.T) {
	query := bson.M{"f": 12}
	update := bson.M{"$push": bson.M{"obj": bson.M{"$each": []int{90, 91, 92}}}}
	u := Updater{c: collection.NewCursor(&testDataUpd)}

	err := u.Update(query, update)
	assert.NoError(t, err)
	assert.Equal(t, collection.Record{"f": 12, "obj": []interface{}{1, 2, 3, 90, 91, 92}, "e": 5}, testDataUpd[4])
}

func TestUpdater_Update_Push_DotNotation(t *testing.T) {
	query := bson.M{"test": "#dot_notation"}
	update := bson.M{"$push": bson.M{"obj.arr": "2"}}
	u := Updater{c: collection.NewCursor(&testDataUpd)}

	err := u.Update(query, update)
	assert.NoError(t, err)
	assert.Equal(t, collection.Record{"test": "#dot_notation", "obj": bson.M{"arr": []interface{}{12, "2"}}}, testDataUpd[5])
}
