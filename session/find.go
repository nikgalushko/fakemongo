package session

import (
	"errors"
	"fakemongo/collection"
	"fakemongo/operations"
	"github.com/globalsign/mgo/bson"
)

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
func (f Finder) Select(selector bson.M) Query {
	f.selector = selector
	return f
}

func (f Finder) One(result interface{}) error {
	var (
		r   collection.Record
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

		if !matched {
			continue
		}

		r = r.WithFields(f.selector)
		data, err := bson.Marshal(r)
		if err != nil {
			return err
		}
		return bson.Unmarshal(data, result)
	}

	return errors.New("not found")
}

func (f Finder) All(result interface{}) error {
	panic("All is unimplemented")
}
