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

type ASTNode any
type Token map[string]ASTNode
type Parser func(c *Context) (ASTNode, bool)

func Space() Parser {
	return func(c *Context) (ASTNode, bool) {
		for c.pos < len(c.stream) && unicode.IsSpace(rune(c.stream[c.pos])) {
			c.pos++
		}
		return nil, true
	}
}

func Str(s string) Parser {
	return func(c *Context) (ASTNode, bool) {
		start := c.pos
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
	return func(c *Context) (ASTNode, bool) {
		start := c.pos
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
	return func(c *Context) (ASTNode, bool) {
		start := c.pos
		var results []ASTNode
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

func Choice(parslets ...Parser) Parser {
	return func(c *Context) (ASTNode, bool) {
		start := c.pos
		for _, p := range parslets {
			if val, ok := p(c); ok {
				return val, true
			}
			c.pos = start
		}
		return nil, false
	}
}

func Tag(key string, p Parser) Parser {
	return func(c *Context) (ASTNode, bool) {
		val, ok := p(c)
		if !ok {
			return nil, false
		}
		return Token{key: val}, true
	}
}
