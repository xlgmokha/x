package peg

func Tag(key string, p Parser) Parser {
	return func(c *Context) (ASTNode, bool) {
		val, ok := p(c)
		if !ok {
			return nil, false
		}
		return Token{key: val}, true
	}
}
