package session

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jetuuuu/fakemongo/collection"
	"github.com/jetuuuu/fakemongo/operations"
	"github.com/jetuuuu/fakemongo/utils"
)

type Session struct {
	data map[string]collection.Collection
}

type Query interface {
	One(interface{}) error
	All(interface{}) error
	Select(interface{}) Query
	Apply(mgo.Change, interface{}) (*mgo.ChangeInfo, error)
	Count() (int, error)
	Sort(...string) Query
	Limit(int) Query
	Collation(*mgo.Collation) Query
}

func NewSession(collections []collection.Collection) Session {
	data := make(map[string]collection.Collection, len(collections))
	for _, c := range collections {
		data[c.Name] = c
	}
	return Session{data: data}
}

// todo query must be an interface{}
func (s Session) Find(collectionName string, query bson.M) Query {
	return NewFinder(query, s.data[collectionName].Cursor())
}

// todo query/set must be an interface{}
func (s Session) Update(collectionName string, query, update bson.M) error {
	u := Updater{c: s.data[collectionName].Cursor()}
	return u.Update(query, update)
}

// todo query/set must be an interface{}
func (s Session) Upsert(collectionName string, query, update bson.M) error {
	c := s.data[collectionName].Cursor()
	f := NewFinder(query, c)

	r := make(collection.Record)
	err := f.One(&r)
	needToInsert := false
	if errors.Is(err, notFoundError) {
		needToInsert = true
	} else if err != nil {
		return err
	}

	parsedOperations := UpdateParameterParser{}.ParseUpdate(update)
	kind := operations.Update
	if needToInsert {
		kind = operations.Insert
	}
	for _, op := range parsedOperations {
		r = op.Upsert(r, kind)
	}

	if needToInsert {
		return s.Insert(collectionName, r)
	} else {
		c.SetCurrent(r)
	}

	return nil
}

func (s Session) Insert(collectionName string, docs ...interface{}) error {
	for _, d := range docs {
		var item bson.M
		switch d.(type) {
		case bson.M:
			item = d.(bson.M)
		default:
			converted, err := utils.ToBsonM(d)
			if err != nil {
				return err
			}

			item = converted
		}

		c := s.data[collectionName].Cursor()
		c.Insert(collection.Record(item))
	}

	return nil
}

func (s Session) Remove(collectionName string, selector interface{}) error {
	query, err := utils.ToBsonM(selector)
	if err != nil {
		return err
	}

	c := s.data[collectionName].Cursor()
	f := NewFinder(query, c)

	if err := f.One(nil); err != nil {
		return err
	}

	f.c.RemoveCurrent()
	return nil
}
