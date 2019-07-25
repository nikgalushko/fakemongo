package collection

type Collection struct {
	Name string
	Data *[]Record
	c    *Cursor
}

func NewCollection(name string, data []Record) Collection {
	coll := Collection{
		Name: name,
		Data: &data,
	}

	coll.c = NewCursor(coll.Data)
	return coll
}

func (c Collection) Cursor() *Cursor {
	if c.c.lastError != nil {
		_ = c.c.Seek(0)
	}

	return c.c
}
