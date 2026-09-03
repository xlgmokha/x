package peg

import (
	"regexp"
	"strings"
	"unicode"
)

type ASTNode any
type Token map[string]ASTNode
type Parser func(c *Context) (ASTNode, bool)

func Space() Parser {
	return func(c *Context) (ASTNode, bool) {
		for c.position < len(c.stream) && unicode.IsSpace(rune(c.stream[c.position])) {
			c.position++
		}
		return nil, true
	}
}

func Str(s string) Parser {
	return func(c *Context) (ASTNode, bool) {
		start := c.position
		if c.position+len(s) <= len(c.stream) && strings.EqualFold(c.stream[c.position:c.position+len(s)], s) {
			c.position += len(s)
			return s, true
		}
		c.position = start
		return nil, false
	}
}

func Match(pattern string) Parser {
	re := regexp.MustCompile("^" + pattern)
	return func(c *Context) (ASTNode, bool) {
		start := c.position
		loc := re.FindStringIndex(c.stream[c.position:])
		if loc != nil {
			res := c.stream[c.position : c.position+loc[1]]
			c.position += loc[1]
			return res, true
		}
		c.position = start
		return nil, false
	}
}

func Sequence(parslets ...Parser) Parser {
	return func(c *Context) (ASTNode, bool) {
		start := c.position
		var results []ASTNode
		for _, p := range parslets {
			val, ok := p(c)
			if !ok {
				c.position = start
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
		start := c.position
		for _, p := range parslets {
			if val, ok := p(c); ok {
				return val, true
			}
			c.position = start
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
