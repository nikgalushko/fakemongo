package session

import (
	"fakemongo/collection"
	"github.com/globalsign/mgo/bson"
)

type Updater struct {
	f Finder
	c *collection.Cursor
}

func (u Updater) Update(selector, update bson.M) error {
	u.f = NewFinder(selector, u.c)

	r := make(collection.Record)
	if err := u.f.One(&r); err != nil {
		return err
	}

	operations := UpdateParameterParser{}.ParseUpdate(update)
	for _, op := range operations {
		r = op.Update(r)
	}

	u.c.SetCurrent(r)

	return nil
}
