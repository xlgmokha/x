package peg

import (
	"strings"
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

func Str(s string) Parser {
	return func(c *Context) (any, bool) {
		start := c.pos
		Space()(c)
		if c.pos+len(s) <= len(c.src) && strings.EqualFold(c.src[c.pos:c.pos+len(s)], s) {
			c.pos += len(s)
			return s, true
		}
		c.pos = start
		return nil, false
	}
}
