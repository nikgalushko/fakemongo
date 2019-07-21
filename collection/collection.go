package collection

type Collection struct {
	Name string
	data []Record
	c    *Cursor
}

func NewCollection(name string, data []Record) Collection {
	return Collection{
		Name: name,
		data: data,
		c:    NewCursor(data),
	}
}

func (c Collection) Cursor() *Cursor {
	if c.c.lastError != nil {
		_ = c.c.Seek(0)
	}

	return c.c
}

/*func (r Record) Match(template bson.M) bool {
	var ret bool
	for k, operatorExpression := range template {
		key := Key(k)
		if key.IsCmd() {
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
			case "$or":
				nextTemplates := expected.([]bson.M)
				// todo nextTemplates should be sorted by priority
				for _, t := range nextTemplates {
					ret = r.Match(t)
					if ret {
						break
					}
				}
			default:
				panic(unimplemented)
			}
		} else if key.IsDotNotation() {
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
		case "$all":
			panic(unimplemented)
		default:
			panic(unimplemented)
		}

		if !ret {
			return false
		}
	}

	return true
}*/
