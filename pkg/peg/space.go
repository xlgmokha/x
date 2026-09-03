package peg

import "unicode"

func Space() Parser {
	return func(c *Context) (ASTNode, bool) {
		for c.position < len(c.stream) && unicode.IsSpace(rune(c.stream[c.position])) {
			c.position++
		}
		return nil, true
	}
}
