package scim

import "fmt"

type Visitor[T any] interface {
	VisitAnd(left, right T) (T, error)
	VisitOr(left, right T) (T, error)
	VisitNot(operand T) (T, error)
	VisitEquals(attribute AttrPath, value any) (T, error)
	VisitNotEquals(attribute AttrPath, value any) (T, error)
	VisitContains(attribute AttrPath, value any) (T, error)
	VisitStartsWith(attribute AttrPath, value any) (T, error)
	VisitEndsWith(attribute AttrPath, value any) (T, error)
	VisitGreaterThan(attribute AttrPath, value any) (T, error)
	VisitGreaterThanEquals(attribute AttrPath, value any) (T, error)
	VisitLessThan(attribute AttrPath, value any) (T, error)
	VisitLessThanEquals(attribute AttrPath, value any) (T, error)
	VisitPresence(attribute AttrPath) (T, error)
	VisitValuePath(path AttrPath, subAttribute string, valueFilter func() (T, error)) (T, error)
}

func Visit[T any](v Visitor[T], e Expression) (T, error) {
	var zero T
	switch e := e.(type) {
	case Comparison:
		return visitComparison(v, e)
	case Logical:
		left, err := Visit(v, e.Left)
		if err != nil {
			return zero, err
		}
		right, err := Visit(v, e.Right)
		if err != nil {
			return zero, err
		}
		if e.Operator == Or {
			return v.VisitOr(left, right)
		}
		return v.VisitAnd(left, right)
	case Not:
		operand, err := Visit(v, e.Operand)
		if err != nil {
			return zero, err
		}
		return v.VisitNot(operand)
	case ValuePath:
		return v.VisitValuePath(e.Attribute, e.SubAttribute, func() (T, error) {
			return Visit(v, e.Filter)
		})
	default:
		return zero, fmt.Errorf("scim: cannot visit %T", e)
	}
}

func visitComparison[T any](v Visitor[T], c Comparison) (T, error) {
	switch c.Operator {
	case Equal:
		return v.VisitEquals(c.Attribute, c.Value)
	case NotEqual:
		return v.VisitNotEquals(c.Attribute, c.Value)
	case Contains:
		return v.VisitContains(c.Attribute, c.Value)
	case StartsWith:
		return v.VisitStartsWith(c.Attribute, c.Value)
	case EndsWith:
		return v.VisitEndsWith(c.Attribute, c.Value)
	case GreaterThan:
		return v.VisitGreaterThan(c.Attribute, c.Value)
	case GreaterOrEqual:
		return v.VisitGreaterThanEquals(c.Attribute, c.Value)
	case LessThan:
		return v.VisitLessThan(c.Attribute, c.Value)
	case LessOrEqual:
		return v.VisitLessThanEquals(c.Attribute, c.Value)
	case Present:
		return v.VisitPresence(c.Attribute)
	default:
		var zero T
		return zero, fmt.Errorf("scim: unknown operator %q", c.Operator)
	}
}
