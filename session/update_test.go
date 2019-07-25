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
	{"f": 10, "arr": []interface{}{"t", "e", "s", "t"}},
	{"e": 4},
	{"f": 12, "obj": bson.M{"e": []interface{}{1, 2, 3}, "f": 18.2}, "e": 5},
	{"e": 5, "arr": []bson.M{{"price": 12}, {"price": 14}}, "shop": "#1"},
	{"e": 7, "arr": []bson.M{{"price": 45, "g": 14}, {"price": 36, "g": 16}}, "shop": "#2"},
}

func TestUpdater_Update(t *testing.T) {
	query := bson.M{"e": 4}
	update := bson.M{"$set": bson.M{"e2": 5, "g3": "test"}}
	u := Updater{c: collection.NewCursor(&testDataUpd)}

	err := u.Update(query, update)
	assert.NoError(t, err)
	assert.Equal(t, collection.Record{"e": 4, "e2": 5, "g3": "test"}, testDataUpd[3])
}
