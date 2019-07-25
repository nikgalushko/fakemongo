package collection

import "errors"

var EOF = errors.New("EOF")

type Cursor struct {
	data            *[]Record
	currentPosition uint
	prevPosition    uint
	lastError       error
}

func NewCursor(data *[]Record) *Cursor {
	if data == nil {
		data = &[]Record{}
	}
	return &Cursor{
		data:            data,
		currentPosition: 0,
		prevPosition:    0,
	}
}

func (c *Cursor) Seek(position uint) error {
	if position >= uint(len(*c.data)) {
		err := errors.New("new position is greater than data length")
		c.lastError = err
		return err
	}

	c.currentPosition = position
	c.prevPosition = position

	c.lastError = nil

	return nil
}

func (c *Cursor) Next() (Record, error) {
	if !c.HasNext() {
		c.lastError = EOF
		return nil, EOF
	}

	data := *c.data
	r := data[c.currentPosition]
	c.prevPosition = c.currentPosition
	c.currentPosition += 1
	c.lastError = nil

	return r, nil
}

func (c *Cursor) Current() (Record, error) {
	if c.lastError != nil {
		return nil, c.lastError
	}

	data := *c.data
	return data[c.prevPosition], nil
}

func (c *Cursor) SetCurrent(r Record) {
	data := *c.data
	data[c.prevPosition] = r

	*c.data = data
}

func (c *Cursor) Insert(r Record) {
	data := *c.data
	data = append(data, r)

	*c.data = data
}

func (c Cursor) HasNext() bool {
	return c.currentPosition < uint(len(*c.data))
}
