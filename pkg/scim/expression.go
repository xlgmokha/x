package scim

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Expression is a node in a parsed SCIM filter (RFC 7644 3.4.2.2). The concrete
// types are Comparison, Logical, Not, and ValuePath. String renders canonical,
// re-parseable filter text.
type Expression interface {
	fmt.Stringer
	expr()
}

type CompareOperator string

const (
	Present        CompareOperator = "pr"
	Equal          CompareOperator = "eq"
	NotEqual       CompareOperator = "ne"
	Contains       CompareOperator = "co"
	StartsWith     CompareOperator = "sw"
	EndsWith       CompareOperator = "ew"
	GreaterThan    CompareOperator = "gt"
	LessThan       CompareOperator = "lt"
	GreaterOrEqual CompareOperator = "ge"
	LessOrEqual    CompareOperator = "le"
)

func (o CompareOperator) valid() bool {
	switch o {
	case Present, Equal, NotEqual, Contains, StartsWith, EndsWith,
		GreaterThan, LessThan, GreaterOrEqual, LessOrEqual:
		return true
	}
	return false
}

type LogicalOperator string

const (
	And LogicalOperator = "and"
	Or  LogicalOperator = "or"
)

// Comparison is attrExp: an attribute compared against a value, or a presence
// test. Value is nil when Operator is Present.
type Comparison struct {
	Attribute AttrPath
	Operator  CompareOperator
	Value     any
}

// Logical is logExp: two sub-filters joined by "and" / "or".
type Logical struct {
	Operator    LogicalOperator
	Left, Right Expression
}

// Not is a negated, parenthesized sub-filter.
type Not struct {
	Operand Expression
}

// ValuePath is attrPath "[" valFilter "]" [subAttr]. SubAttribute is "" when
// absent; attribute names inside Filter are relative to Attribute.
type ValuePath struct {
	Attribute    AttrPath
	Filter       Expression
	SubAttribute string
}

func (Comparison) expr() {}
func (Logical) expr()    {}
func (Not) expr()        {}
func (ValuePath) expr()  {}

func (c Comparison) String() string {
	if c.Operator == Present {
		return fmt.Sprintf("%s pr", c.Attribute)
	}
	return fmt.Sprintf("%s %s %s", c.Attribute, c.Operator, stringifyValue(c.Value))
}

func (l Logical) String() string {
	return fmt.Sprintf("(%s %s %s)", l.Left, l.Operator, l.Right)
}

func (n Not) String() string {
	return fmt.Sprintf("not (%s)", n.Operand)
}

func (v ValuePath) String() string {
	if v.SubAttribute == "" {
		return fmt.Sprintf("%s[%s]", v.Attribute, v.Filter)
	}
	return fmt.Sprintf("%s[%s].%s", v.Attribute, v.Filter, v.SubAttribute)
}

func stringifyValue(value any) string {
	switch v := value.(type) {
	case string:
		b, _ := json.Marshal(v)
		return string(b)
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}
