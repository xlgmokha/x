package peg

import "strings"

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
