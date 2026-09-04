package scim

import (
	"errors"

	"github.com/xlgmokha/x/pkg/peg"
)

var errUnexpectedNode = errors.New("scim: unexpected parse node")

// Raw parser token keys, written by the grammar (peg.Tag / peg.Token in
// grammar.go and combinators.go) and read only here in newExpression.
const (
	keyAttribute    = "attribute"
	keyOperator     = "operator"
	keyValue        = "value"
	keyLeft         = "left"
	keyRight        = "right"
	keyNot          = "not"
	keyPath         = "path"
	keyValueFilter  = "value_filter"
	keySubAttribute = "sub_attribute"
)

// newExpression is the single reader of the raw parser token shape: it decodes
// the grammar's peg.Token tree into the typed Expression AST.
func newExpression(raw peg.ASTNode) (Expression, error) {
	tok, ok := raw.(peg.Token)
	if !ok {
		return nil, errUnexpectedNode
	}
	if inner, ok := tok[keyNot]; ok {
		operand, err := newExpression(inner)
		if err != nil {
			return nil, err
		}
		return Not{operand}, nil
	}
	if path, ok := tok[keyPath].(AttrPath); ok {
		filter, err := newExpression(tok[keyValueFilter])
		if err != nil {
			return nil, err
		}
		sub, _ := tok[keySubAttribute].(string)
		return ValuePath{path, filter, sub}, nil
	}
	op, _ := tok[keyOperator].(string)
	switch lop := LogicalOperator(op); lop {
	case And, Or:
		left, err := newExpression(tok[keyLeft])
		if err != nil {
			return nil, err
		}
		right, err := newExpression(tok[keyRight])
		if err != nil {
			return nil, err
		}
		return Logical{lop, left, right}, nil
	}
	cop := CompareOperator(op)
	if !cop.valid() {
		return nil, errUnexpectedNode
	}
	attr, _ := tok[keyAttribute].(AttrPath)
	return Comparison{attr, cop, tok[keyValue]}, nil
}
