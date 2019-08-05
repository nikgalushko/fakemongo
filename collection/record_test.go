package collection

import (
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecord_GetByField(t *testing.T) {
	data := bson.M{
		"a": 5,
		"b": bson.M{
			"c": bson.M{
				"d": bson.M{
					"e": []int{1, 2, 3},
				},
			},
		},
	}

	v1, ok1 := Record(data).GetByField("b.c.d")
	v2, ok2 := Record(data).GetByField("b.c.d.e")
	v3, ok3 := Record(data).GetByField("a")
	v4, ok4 := Record(data).GetByField("b.c.d.e.f")
	a := assert.New(t)

	a.True(ok1)
	a.Equal(bson.M{"e": []int{1, 2, 3}}, v1)

	a.True(ok2)
	a.Equal([]int{1, 2, 3}, v2)

	a.True(ok3)
	a.Equal(5, v3)

	a.False(ok4)
	a.Nil(v4)
}

func TestRecord_GetByField_Array(t *testing.T) {
	data := bson.M{
		"arr": []bson.M{
			{"b": 1, "b2": 2},
			{"b": 5, "b2": 4},
			{"b2": 10, "b": bson.M{"b3_c": []bson.M{
				{"b3_c_d": 5},
				{"b3_c_d": 18},
				{"b3_c_1": 1},
			}}},
		},
		"b": bson.M{
			"c": bson.M{
				"d": bson.M{
					"e": []int{1, 2, 3},
				},
			},
		},
	}

	v1, ok1 := Record(data).GetByField("arr.b.b3_c.b3_c_d")
	a := assert.New(t)

	a.True(ok1)
	a.Equal([]interface{}{nil, nil, 5, 18, nil}, v1)
}
