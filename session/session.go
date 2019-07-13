package session

import (
	"fakemongo/collection"
	"github.com/globalsign/mgo/bson"
)

type Query interface {
	One(interface{}) error
	All(interface{}) error
	Select(bson.M) Query
}

type FindOp struct {
	query    bson.M
	foundAt  int
	selector bson.M
	parser   SelectorParser
	c        *collection.Collection
}

/*func (o *FindOp) One(result interface{}) error {
	for i, r := range o.c.Data {
		if r.Match(o.query) {
			o.foundAt = i
			r = r.WithFields(o.selector)
			data, _ := bson.Marshal(r)
			return bson.Unmarshal(data, result)
		}
	}
	return nil
}

func (o FindOp) All(result interface{}) error {
	var ret []collection.Record
	resultv := reflect.ValueOf(result)
	slicev := resultv.Elem()
	elemt := slicev.Type().Elem()

	for _, r := range o.c.Data {
		if r.Match(o.query) {
			r = r.WithFields(o.selector)
			ret = append(ret, r)
		}
	}

	for _, r := range ret {
		elemp := reflect.New(elemt)
		data, _ := bson.Marshal(r)
		err := bson.Unmarshal(data, elemp.Interface())
		if err != nil {
			panic(err)
		}

		slicev.Set(reflect.Append(slicev, elemp.Elem()))
	}

	resultv.Elem().Set(slicev.Slice(0, len(ret)))

	return nil
}*/

/*func (o FindOp) Select(selector bson.M) Query {
	return &FindOp{query: o.query, selector: selector, c: o.c, foundAt: o.foundAt}
}*/
