package peg

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
		merged := Token{}
		sawToken := false
		for _, r := range results {
			token, ok := r.(Token)
			if !ok {
				continue
			}
			sawToken = true
			for k, v := range token {
				merged[k] = v
			}
		}
		if !sawToken {
			return results, true
		}
		return merged, true
	}
}
