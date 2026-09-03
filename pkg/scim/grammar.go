// Package scim implements a parser for the SCIM RFC-7644 filter grammar.
package scim

import (
	"regexp"
	"strconv"

	"github.com/xlgmokha/x/pkg/peg"
)

type Grammar struct{}

func (g *Grammar) Parse(text string) (*Node, bool) {
	ctx := peg.NewContext(text)
	raw, ok := g.Filter()(ctx)
	if !ok {
		return nil, false
	}
	peg.Space()(ctx)
	if ctx.Position() != len(text) {
		return nil, false
	}
	node := newNode(raw)
	return node, node != nil
}

// FILTER = attrExp / logExp / valuePath / *1"not" "(" FILTER ")"
func (g *Grammar) Filter() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		return peg.Choice(
			g.or(),
			g.and(),
		)(c)
	}
}

func (g *Grammar) FilterAtom() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		if res, ok := g.notParenGroup()(c); ok {
			return res, true
		}
		return peg.Choice(
			g.AttributeExpression(),
			g.ValuePath(),
		)(c)
	}
}

// attrExp = (attrPath SP "pr") / (attrPath SP compareOp SP compValue)
func (g *Grammar) AttributeExpression() peg.Parser {
	return peg.Choice(
		peg.Sequence(
			peg.Tag("attribute", g.AttributePath()),
			peg.Space(),
			peg.Tag("operator", g.Presence()),
		),
		peg.Sequence(
			peg.Tag("attribute", g.AttributePath()),
			peg.Space(),
			peg.Tag("operator", g.ComparisonOperator()),
			peg.Space(),
			peg.Tag("value", g.ComparisonValue()),
		),
	)
}

// valuePath = attrPath "[" valFilter "]"
func (g *Grammar) ValuePath() peg.Parser {
	return peg.Sequence(
		peg.Tag("path", g.AttributePath()),
		peg.Str("["),
		peg.Tag("value_filter", g.ValueFilter()),
		peg.Str("]"),
		peg.Optional(peg.Tag("sub_attribute", g.SubAttribute())),
	)
}

// valFilter = attrExp / logExp / *1"not" "(" valFilter ")"
func (g *Grammar) ValueFilter() peg.Parser {
	return g.Filter()
}

// attrPath = [URI ":"] ATTRNAME *1subAttr
func (g *Grammar) AttributePath() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		start := c.Position()
		prefix := ""
		if uri, ok := g.schemaURI()(c); ok {
			prefix = uri.(string) + ":"
		}
		name, ok := g.AttributeName()(c)
		if !ok {
			c.Seek(start)
			return nil, false
		}
		path := prefix + name.(string)
		if sub, ok := g.SubAttribute()(c); ok {
			path += "." + sub.(string)
		}
		return path, true
	}
}

// ATTRNAME = ALPHA *(nameChar)
// nameChar = "-" / "_" / DIGIT / ALPHA
func (g *Grammar) AttributeName() peg.Parser {
	return peg.Match(regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_\-]*`))
}

// subAttr = "." ATTRNAME
func (g *Grammar) SubAttribute() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		if _, ok := peg.Str(".")(c); !ok {
			return nil, false
		}
		return g.AttributeName()(c)
	}
}

// compareOp = "eq" / "ne" / "co" / "sw" / "ew" / "gt" / "lt" / "ge" / "le"
func (g *Grammar) ComparisonOperator() peg.Parser {
	return peg.Choice(
		g.Equal(), g.NotEqual(), g.Contains(), g.StartsWith(), g.EndsWith(),
		g.GreaterThan(), g.LessThan(), g.GreaterThanEquals(), g.LessThanEquals(),
	)
}

func (g *Grammar) Equal() peg.Parser             { return peg.Str("eq") }
func (g *Grammar) NotEqual() peg.Parser          { return peg.Str("ne") }
func (g *Grammar) Contains() peg.Parser          { return peg.Str("co") }
func (g *Grammar) StartsWith() peg.Parser        { return peg.Str("sw") }
func (g *Grammar) EndsWith() peg.Parser          { return peg.Str("ew") }
func (g *Grammar) GreaterThan() peg.Parser       { return peg.Str("gt") }
func (g *Grammar) LessThan() peg.Parser          { return peg.Str("lt") }
func (g *Grammar) GreaterThanEquals() peg.Parser { return peg.Str("ge") }
func (g *Grammar) LessThanEquals() peg.Parser    { return peg.Str("le") }

// "pr" ; presence operator, attrExp's other alternative
func (g *Grammar) Presence() peg.Parser { return peg.Str("pr") }

// "and" ; logExp operator
func (g *Grammar) AndOp() peg.Parser { return peg.Str("and") }

// "or" ; logExp operator
func (g *Grammar) OrOp() peg.Parser { return peg.Str("or") }

// "not" ; FILTER/valFilter negation prefix
func (g *Grammar) NotOp() peg.Parser { return peg.Str("not") }

// compValue = false / null / true / number / string
func (g *Grammar) ComparisonValue() peg.Parser {
	return peg.Choice(
		g.Falsey(),
		g.Null(),
		g.Truthy(),
		g.Number(),
		g.StringLiteral(),
	)
}

func (g *Grammar) Falsey() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		if _, ok := peg.Str("false")(c); ok {
			return false, true
		}
		return nil, false
	}
}

func (g *Grammar) Truthy() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		if _, ok := peg.Str("true")(c); ok {
			return true, true
		}
		return nil, false
	}
}

func (g *Grammar) Null() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		if _, ok := peg.Str("null")(c); ok {
			return nil, true
		}
		return nil, false
	}
}

func (g *Grammar) Number() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		val, ok := peg.Match(regexp.MustCompile(`-?[0-9]+(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?`))(c)
		if !ok {
			return nil, false
		}
		n, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return nil, false
		}
		return n, true
	}
}

func (g *Grammar) StringLiteral() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		val, ok := peg.Match(regexp.MustCompile(`"(?:[^"\\]|\\.)*"`))(c)
		if !ok {
			return nil, false
		}
		s := val.(string)
		return s[1 : len(s)-1], true
	}
}

func (g *Grammar) or() peg.Parser {
	return peg.Sequence(
		peg.Tag("left", g.and()),
		peg.Space(),
		peg.Tag("operator", g.OrOp()),
		peg.Space(),
		peg.Tag("right", g.Filter()),
	)
}

func (g *Grammar) and() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		return peg.Choice(
			peg.Sequence(
				peg.Tag("left", g.FilterAtom()),
				peg.Space(),
				peg.Tag("operator", g.AndOp()),
				peg.Space(),
				peg.Tag("right", g.and()),
			),
			g.FilterAtom(),
		)(c)
	}
}

func (g *Grammar) notParenGroup() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		start := c.Position()
		_, negated := peg.Sequence(g.NotOp(), peg.Space())(c)
		if _, ok := peg.Str("(")(c); !ok {
			c.Seek(start)
			return nil, false
		}
		inner, ok := g.Filter()(c)
		if !ok {
			c.Seek(start)
			return nil, false
		}
		if _, ok := peg.Str(")")(c); !ok {
			c.Seek(start)
			return nil, false
		}
		if !negated {
			return inner, true
		}
		return peg.Token{"not": inner}, true
	}
}

func (g *Grammar) schemaURI() peg.Parser {
	return func(c *peg.Context) (peg.ASTNode, bool) {
		val, ok := peg.Match(regexp.MustCompile(`[A-Za-z][A-Za-z0-9.\-]*(?::[A-Za-z0-9.\-]+)*:`))(c)
		if !ok {
			return nil, false
		}
		s := val.(string)
		return s[:len(s)-1], true
	}
}
