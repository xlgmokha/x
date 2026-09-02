package peg

import (
	"regexp"
	"strings"
	"unicode"
)

type Context struct {
	stream string
	pos    int
}

type Parser func(c *Context) (any, bool)

func Space() Parser {
	return func(c *Context) (any, bool) {
		for c.pos < len(c.stream) && unicode.IsSpace(rune(c.stream[c.pos])) {
			c.pos++
		}
		return nil, true
	}
}

func Str(s string) Parser {
	return func(c *Context) (any, bool) {
		start := c.pos
		Space()(c)
		if c.pos+len(s) <= len(c.stream) && strings.EqualFold(c.stream[c.pos:c.pos+len(s)], s) {
			c.pos += len(s)
			return s, true
		}
		c.pos = start
		return nil, false
	}
}

func Match(pattern string) Parser {
	re := regexp.MustCompile("^" + pattern)
	return func(c *Context) (any, bool) {
		start := c.pos
		Space()(c)
		loc := re.FindStringIndex(c.stream[c.pos:])
		if loc != nil {
			res := c.stream[c.pos : c.pos+loc[1]]
			c.pos += loc[1]
			return res, true
		}
		c.pos = start
		return nil, false
	}
}

func Sequence(parslets ...Parser) Parser {
	return func(c *Context) (any, bool) {
		start := c.pos
		var results []any
		for _, p := range parslets {
			val, ok := p(c)
			if !ok {
				c.pos = start
				return nil, false
			}
			if val != nil {
				results = append(results, val)
			}
		}
		return results, true
	}
}
