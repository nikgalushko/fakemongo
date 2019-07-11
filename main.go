package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fakemongo/compare"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"strings"
)

var unimplemented = errors.New("unimplemented")

type Client struct {
	Id                         bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Domain                     string        `bson:"domain,omitempty" json:"domain"`
	Value                      float64       `bson:"val"`
	Status                     string        `bson:"status,omitempty" json:"status"`
	CityID                     string        `bson:"cityId,omitempty" json:"city_id,omitempty"`
	Tariff                     string        `bson:"tariff,omitempty" json:"tariff"`
	Enabled                    bool          `bson:"enabled,omitempty" json:"-"`
	IsTest                     bool          `bson:"is_test,omitempty" json:"is_test"`
	IsDemo                     bool          `bson:"is_demo" json:"is_demo"`
	ContractOATS               bool          `bson:"contract_signature,omitempty" json:"-"`
	ContractTelephony          bool          `bson:"contract_telephony,omitempty" json:"-"`
	ContractTelephonyWorldWide bool          `bson:"contract_telephony_world_wide"`
	Obj                        Object        `bson:"object"`
	Telnums                    []Telnum      `bson:"telnums"`
}

type Telnum struct {
	Telnum      string `bson:"telnum"`
	City        string `bson:"city"`
	Domainscity string `bson:"domainscity"`
}

type Object struct {
	Field1      string                 `bson:"field1"`
	Flag        bool                   `bson:"flag"`
	InnerObject map[string]interface{} `bson:"inner_obj"`
}

func main() {
	c := NewCollection("domains", testData)
	client := Client{}
	err := c.Find(bson.M{"domain": "vvuri-030719"}).Select(bson.M{"cityId": 1}).One(&client)
	if err != nil {
		panic(err)
	}

	fmt.Println("One: ", client)

	client = Client{}
	err = c.Find(bson.M{"domain": "vvuri-030719", "cityId": bson.M{"$nin": []string{"888", "999"}}}).One(&client)
	if err != nil {
		panic(err)
	}

	fmt.Println("One with $nin: ", client)

	var clients []Client
	err = c.Find(bson.M{"cityId": bson.M{"$nin": []string{"888", "999"}}}).All(&clients)
	if err != nil {
		panic(err)
	}

	fmt.Println(clients)

	var clients2 []Client
	err = c.Find(bson.M{"is_demo": bson.M{"$ne": false}}).All(&clients2)
	if err != nil {
		panic(err)
	}

	fmt.Println("One with $ne: ", clients2)

	err = c.Update(bson.M{"domain": "nikita2"}, bson.M{"$set": bson.M{"tariff": "ITooLabsIndividualRate"}})
	if err != nil {
		panic(err)
	}

	client = Client{}
	err = c.Find(bson.M{"tariff": "ITooLabsIndividualRate"}).One(&client)
	if err != nil {
		panic(err)
	}

	fmt.Println("tariff ITooLabsIndividualRate: ", client)

	client = Client{}
	err = c.Find(bson.M{"object.inner_obj.field1": "inner_field_value1"}).One(&client)
	if err != nil {
		panic(err)
	}

	fmt.Println("client Object: ", client.Obj)

	client = Client{}
	err = c.Find(bson.M{"val": bson.M{"$gt": float64(3)}}).One(&client)
	if err != nil {
		panic(err)
	}

	fmt.Println("client \"$gt\": 3: ", client)
	client = Client{}
	err = c.Find(bson.M{"val": bson.M{"$exists": true}}).One(&client)
	if err != nil {
		panic(err)
	}

	fmt.Println("client \"$exists\": true: ", client)

	client = Client{}
	err = c.Find(bson.M{"telnums": bson.M{"$size": 1}}).One(&client)
	if err != nil {
		panic(err)
	}

	fmt.Println("telnums \"$size\": 1: ", client)

	var clients3 []Client
	err = c.Find(bson.M{"val": bson.M{"$not": bson.M{"$gte": float64(6)}}}).All(&clients3)
	if err != nil {
		panic(err)
	}
	fmt.Println("client \"$not\": {$gt: 3}: ", clients3)

	var clients4 []Client
	err = c.Find(bson.M{"$and": []bson.M{
		{"val": bson.M{"$gt": float64(4)}},
		{"is_demo": true},
	}}).All(&clients4)
	if err != nil {
		panic(err)
	}
	fmt.Println("client \"$and\": ", clients4)
}

type FakeMongo struct {
	collections map[string]Collection
}

type Collection struct {
	name string
	data []Record
}

// TODO: jsonData maybe just single object
func NewCollection(name, jsonData string) Collection {
	var records []map[string]interface{}
	_ = json.Unmarshal([]byte(jsonData), &records)

	c := Collection{
		name: name,
	}

	for _, r := range records {
		c.data = append(c.data, Record(r))
	}

	return c
}

func (c Collection) Find(query bson.M) Query {
	return &FindOp{query: query, c: &c}
}

func (c Collection) Update(selector bson.M, update bson.M) error {
	r := make(Record)
	newRecord := make(Record)
	query := c.Find(selector).(*FindOp)
	err := query.One(&r)
	if err != nil {
		return err
	}

	for k, v := range update {
		if isCmd(k) {
			switch k {
			case "$set":
				for field, newValue := range v.(bson.M) {
					r[field] = newValue
				}
			}
		} else {
			newRecord[k] = v
		}
	}

	if len(newRecord) != 0 {
		r = newRecord
	}

	c.data[query.foundAt] = r
	return nil
}

type FindOp struct {
	query    bson.M
	foundAt  int
	selector bson.M
	c        *Collection
}

func (o *FindOp) One(result interface{}) error {
	for i, r := range o.c.data {
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
	var ret []Record
	resultv := reflect.ValueOf(result)
	slicev := resultv.Elem()
	elemt := slicev.Type().Elem()

	for _, r := range o.c.data {
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
}

func (o FindOp) Select(selector bson.M) Query {
	return &FindOp{query: o.query, selector: selector, c: o.c, foundAt: o.foundAt}
}

type Query interface {
	One(interface{}) error
	All(interface{}) error
	Select(bson.M) Query
}

type Record map[string]interface{}

func (r Record) Match(template bson.M) bool {
	var ret bool
	for k, expected := range template {
		if isCmd(k) {
			switch k {
			case "$and":
				nextTemplates := expected.([]bson.M)
				// todo nextTemplates should be sorted by priority
				for _, t := range nextTemplates {
					ret = r.Match(t)
					if !ret {
						break
					}
				}
			default:
				panic(unimplemented)
			}
		} else if isObj(k) {
			subFields := strings.Split(k, ".")
			mainKey := subFields[0]
			otherKeys := strings.Join(subFields[1:], ".")
			if _, ok := r[mainKey]; !ok {
				return false
			}
			subRecord := r[mainKey].(Record)
			ret = subRecord.Match(bson.M{otherKeys: expected})
		} else {
			switch expected.(type) {
			case bson.M:
				ret = r.FieldMatch(k, expected.(bson.M))
			default:
				if actualy, ok := r[k]; ok && reflect.DeepEqual(expected, actualy) {
					ret = true
				} else {
					ret = false
				}
			}
		}

		if !ret {
			return false
		}
	}

	return true
}

func (r Record) WithFields(fields bson.M) Record {
	if fields == nil {
		return r
	}
	ret := make(Record)
	for f, v := range fields {
		switch v.(type) {
		case int:
			if value, ok := r[f]; v.(int) == 1 && ok {
				ret[f] = value
			}
		}
	}

	return ret
}

func (r Record) FieldMatch(f string, template bson.M) bool {
	var ret bool
	value := r[f]

	for k, expected := range template {
		switch k {
		case "$nin":
			ok, found := includeElement(expected, value)
			ret = ok && !found
		case "$in":
			ok, found := includeElement(expected, value)
			ret = ok && found
		case "$ne":
			ret = !reflect.DeepEqual(expected, value)
		case "$eq":
			ret = reflect.DeepEqual(expected, value)
		case "$gt", "$gte", "$lt", "$lte":
			ret = compare.CompareTo(value, expected) == string(k[1:])
		case "$exists":
			if expected.(bool) {
				ret = value != nil
			} else {
				ret = value == nil
			}
		case "$size":
			size := expected.(int)
			arr, ok := value.([]interface{})
			ret = ok && len(arr) == size
		case "$not":
			nextTemplate := expected.(bson.M)
			ret = !r.FieldMatch(f, nextTemplate)
		default:
			panic(unimplemented)
		}

		if !ret {
			return false
		}
	}

	return true
}

func isCmd(key string) bool {
	return strings.HasPrefix(key, "$")
}

func isObj(key string) bool {
	return strings.Contains(key, ".")
}

func includeElement(list interface{}, element interface{}) (ok, found bool) {
	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if ObjectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if ObjectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false

}

func ObjectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

var testData = `
[
{
    "domain" : "nikita95",
    "cityId" : "9999",
	"val": 5,
    "tariff" : "ITooLabsPro",
    "is_demo" : false,
    "contract_signature" : true,
    "contract_telephony_world_wide" : true
},
{
    "domain" : "nikita2",
    "cityId" : "9999",
    "tariff" : "ITooLabsPro",
    "is_demo" : true,
	"val": 6,
    "contract_signature" : true,
    "contract_telephony_world_wide" : true,
    "contract_telephony" : true,
	"object": {
		"field1": "field1_value",
		"flag": true,
		"inner_obj": {
			"field1": "inner_field_value",
			"flag": false
		}
	},
    "telnums" : [ 
        {
            "telnum" : "78003576242",
            "city" : "0",
            "domainscity" : "9999"
        }
    ],
    "extensions" : [ 
        {
            "id" : "63vdtxqdxa56b25nn2392rehe2",
            "name" : "AddCallRecord",
            "cityId" : "",
            "sum" : 1
        }
    ]
},
{
    "domain" : "vvuri-030719",
    "cityId" : "9999",
	"val": 4,
    "tariff" : "string",
    "is_test" : true,
    "is_demo" : true,
    "contract_signature" : true,
    "contract_telephony_world_wide" : false,
    "contract_telephony" : true
}]
`
