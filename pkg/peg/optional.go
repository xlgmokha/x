package peg

func Optional(p Parser) Parser {
	return func(c *Context) (ASTNode, bool) {
		start := c.position
		if val, ok := p(c); ok {
			return val, true
		}
		c.position = start
		return nil, true
	}
}
