package collection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testData = &[]Record{
	{"a": 1, "b": 2},
	{"c": 3, "d": 4, "e": "test string"},
	{"f": 0.2, "h": 0, "a": 5},
}

func TestCursor_Seek(t *testing.T) {
	c := NewCursor(testData)
	err := c.Seek(2)

	assert.NoError(t, err)
	assert.NoError(t, c.lastError)
	assert.Equal(t, c.currentPosition, uint(2))
}

func TestCursor_Seek_Error(t *testing.T) {
	c := NewCursor(testData)
	err := c.Seek(20)

	assert.Error(t, err)
	assert.Error(t, c.lastError)
	assert.Equal(t, c.currentPosition, uint(0))
}

func TestCursor_Current(t *testing.T) {
	c := NewCursor(testData)

	_, err1 := c.Next()
	_, err2 := c.Next()
	current, err3 := c.Current()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	assert.Equal(t, (*testData)[1], current)
}

func TestCursor_HasNext(t *testing.T) {
	c := NewCursor(testData)
	b1 := c.HasNext()
	_, err1 := c.Next()

	b2 := c.HasNext()
	_, err2 := c.Next()

	b3 := c.HasNext()
	_, err3 := c.Next()

	b4 := c.HasNext()
	_, err4 := c.Next()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	assert.True(t, b1)
	assert.True(t, b2)
	assert.True(t, b3)

	assert.Error(t, err4)
	assert.False(t, b4)
}
