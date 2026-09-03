package scim

import (
	"fmt"
	"strings"
)

type Visitor[T any] interface {
	VisitAnd(left, right T) (T, error)
	VisitOr(left, right T) (T, error)
	VisitNot(operand T) (T, error)
	VisitEquals(attribute string, value any) (T, error)
	VisitNotEquals(attribute string, value any) (T, error)
	VisitContains(attribute string, value any) (T, error)
	VisitStartsWith(attribute string, value any) (T, error)
	VisitEndsWith(attribute string, value any) (T, error)
	VisitGreaterThan(attribute string, value any) (T, error)
	VisitGreaterThanEquals(attribute string, value any) (T, error)
	VisitLessThan(attribute string, value any) (T, error)
	VisitLessThanEquals(attribute string, value any) (T, error)
	VisitPresence(attribute string) (T, error)
	VisitValuePath(path string, subAttribute string, valueFilter func() (T, error)) (T, error)
}

func Visit[T any](v Visitor[T], n *Node) (T, error) {
	var zero T
	if n == nil {
		return zero, fmt.Errorf("scim: cannot visit a nil node")
	}
	if n.Not() {
		operand, err := Visit(v, n.Operand())
		if err != nil {
			return zero, err
		}
		return v.VisitNot(operand)
	}
	return dispatch(v, n)
}

func dispatch[T any](v Visitor[T], n *Node) (T, error) {
	var zero T

	if n.HasPath() {
		return v.VisitValuePath(n.Path(), n.SubAttribute(), func() (T, error) {
			return Visit(v, n.ValueFilter())
		})
	}

	switch strings.ToLower(n.Operator()) {
	case "and":
		return visitBinary(v, n, v.VisitAnd)
	case "or":
		return visitBinary(v, n, v.VisitOr)
	case "eq":
		return v.VisitEquals(n.Attribute(), n.Value())
	case "ne":
		return v.VisitNotEquals(n.Attribute(), n.Value())
	case "co":
		return v.VisitContains(n.Attribute(), n.Value())
	case "sw":
		return v.VisitStartsWith(n.Attribute(), n.Value())
	case "ew":
		return v.VisitEndsWith(n.Attribute(), n.Value())
	case "gt":
		return v.VisitGreaterThan(n.Attribute(), n.Value())
	case "ge":
		return v.VisitGreaterThanEquals(n.Attribute(), n.Value())
	case "lt":
		return v.VisitLessThan(n.Attribute(), n.Value())
	case "le":
		return v.VisitLessThanEquals(n.Attribute(), n.Value())
	case "pr":
		return v.VisitPresence(n.Attribute())
	default:
		return zero, fmt.Errorf("scim: unrecognized node shape: %+v", n.raw)
	}
}

func visitBinary[T any](v Visitor[T], n *Node, combine func(left, right T) (T, error)) (T, error) {
	var zero T
	left, err := Visit(v, n.Left())
	if err != nil {
		return zero, err
	}
	right, err := Visit(v, n.Right())
	if err != nil {
		return zero, err
	}
	return combine(left, right)
}
