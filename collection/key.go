package collection

import "strings"

type KeyType int

const (
	Cmd KeyType = iota
	DotNotation
	Literal
	Illegal
)

type Key string

func (k Key) IsCmd() bool {
	return strings.HasPrefix(string(k), "$")
}

func (k Key) IsDotNotation() bool {
	return strings.Contains(string(k), ".")
}

func (k Key) Type() KeyType {
	if k.IsCmd() {
		return Cmd
	} else if k.IsDotNotation() {
		return DotNotation
	} else {
		return Literal
	}
}
