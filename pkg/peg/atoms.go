package peg

import (
	"unicode"
)

type Context struct {
	src string
	pos int
}

type Parser func(c *Context) (any, bool)

func Space() Parser {
	return func(c *Context) (any, bool) {
		for c.pos < len(c.src) && unicode.IsSpace(rune(c.src[c.pos])) {
			c.pos++
		}
		return nil, true
	}
}
