package scim

import "github.com/xlgmokha/x/pkg/peg"

// binExpr parses `left [SP op SP right]`, tagging a Logical token when the
// operator is present and returning left unchanged otherwise.
func binExpr(left, op peg.Parser, name string, right peg.Parser) peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, error) {
		l, err := left(c)
		if err != nil {
			return nil, err
		}
		start := c.Position()
		peg.Space()(c)
		if _, err := op(c); err != nil {
			c.Seek(start)
			return l, nil
		}
		peg.Space()(c)
		r, err := right(c)
		if err != nil {
			c.Seek(start)
			return l, nil
		}
		return peg.Token{keyLeft: l, keyOperator: name, keyRight: r}, nil
	}
}

// constant matches a literal and yields a fixed value.
func constant(lit string, val peg.ASTNode) peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, error) {
		if _, err := peg.Str(lit)(c); err != nil {
			return nil, err
		}
		return val, nil
	}
}

// convert runs p, then maps its matched string through fn.
func convert(p peg.Parser, fn func(string) (peg.ASTNode, error)) peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, error) {
		start := c.Position()
		val, err := p(c)
		if err != nil {
			return nil, err
		}
		out, err := fn(val.(string))
		if err != nil {
			c.Seek(start)
			return nil, err
		}
		return out, nil
	}
}
