package session

import (
	"github.com/globalsign/mgo/bson"
	"github.com/jetuuuu/fakemongo/collection"
	"github.com/stretchr/testify/assert"
	"testing"
)

type R struct {
	N string `bson:"name"`
}

func TestSession_Insert(t *testing.T) {
	collections := []collection.Collection{
		collection.NewCollection("c1", []collection.Record{
			{"e2": 1},
		}),
	}
	s := NewSession(collections)
	err := s.Insert("c1", R{N: "test"}, bson.M{"_bson_": 15})
	c := s.data["c1"].Cursor()

	assert.NoError(t, err)

	assert.Len(t, *s.data["c1"].Data, 3)

	firstRecord, err := c.Next()
	assert.NoError(t, err)
	assert.Equal(t, collection.Record{"e2": 1}, firstRecord)

	secondRecord, err := c.Next()
	assert.NoError(t, err)
	assert.Equal(t, collection.Record{"name": "test"}, secondRecord)

	thirdRecord, err := c.Next()
	assert.NoError(t, err)
	assert.Equal(t, collection.Record{"_bson_": 15}, thirdRecord)
}
