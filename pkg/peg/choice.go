package peg

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
