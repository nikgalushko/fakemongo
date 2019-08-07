package session

import (
	"github.com/globalsign/mgo/bson"
	"github.com/jetuuuu/fakemongo/collection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession_Remove(t *testing.T) {
	data := []collection.Record{
		{"f": 10, "arr2": "is not array"},
		{"c": true, "obj": bson.M{"f": 15.6}},
	}

	s := NewSession([]collection.Collection{collection.NewCollection("test", data)})
	err := s.Remove("test", bson.M{"obj.f": bson.M{"$gt": 10}})

	assert.NoError(t, err)
	assert.Len(t, *s.data["test"].Data, 1)
}
