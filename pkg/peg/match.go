package peg

import (
	"regexp"
)

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
