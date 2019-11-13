package session

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jetuuuu/fakemongo/collection"
	"github.com/jetuuuu/fakemongo/operations"
	"reflect"
)

var notFoundError = errors.New("not found")

type Finder struct {
	query       bson.M // todo query must be an interface{}
	foundAt     int
	selector    bson.M
	expressions []operations.Expression
	c           *collection.Cursor
}

// todo query must be an interface{}
func NewFinder(query bson.M, c *collection.Cursor) Finder {
	f := Finder{query: query, c: c}
	f.expressions = SelectorParser{}.ParseQuery(query)

	return f
}

// todo selector must be an interface{}
func (f Finder) Select(selector interface{}) Query {
	f.selector = selector.(bson.M)
	return f
}

func (f Finder) One(result interface{}) error {
	arr, err := f.find(true)
	if err != nil {
		return err
	}

	if result == nil {
		return nil
	}

	data, err := bson.Marshal(arr[0])
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, result)
}

func (f Finder) All(result interface{}) error {
	arr, err := f.find(false)
	if err != nil {
		return err
	}

	valuePtr := reflect.ValueOf(result)
	value := valuePtr.Elem()

	for _, r := range arr {
		data, err := bson.Marshal(r)
		if err != nil {
			return err
		}
		v := reflect.New(reflect.TypeOf(result).Elem().Elem())
		i := v.Interface()
		err = bson.Unmarshal(data, i)
		if err != nil {
			return err
		}

		value.Set(reflect.Append(value, v.Elem()))
	}

	return nil
}

func (f Finder) find(one bool) ([]collection.Record, error) {
	var (
		r   collection.Record
		ret []collection.Record
		err error
	)
	for ; err == nil; r, err = f.c.Next() {
		matched := true
		for _, e := range f.expressions {
			if !e.Match(r) {
				matched = false
				break
			}
		}

		if !matched || r == nil {
			continue
		}

		ret = append(ret, r.WithFields(f.selector))
		if one {
			break
		}
	}

	if len(ret) == 0 {
		return nil, notFoundError
	}

	return ret, nil
}

func (f Finder) Apply(change mgo.Change, result interface{}) (*mgo.ChangeInfo, error) {
	return nil, errors.New("unimplemented")
}

func (f Finder) Count() (int, error) {
	arr, err := f.find(false)
	return len(arr), err
}

func (f Finder) Sort(fields ...string) Query {
	panic("unimplemented")
}

func (f Finder) Limit(n int) Query {
	panic("unimplemented")
}

func (f Finder) Collation(collation *mgo.Collation) Query {
	panic("unimplemented")
}
