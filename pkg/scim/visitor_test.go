package scim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringifyVisitorRoundTrip(t *testing.T) {
	tt := []string{
		`userName eq "bjensen"`,
		`userName pr`,
		`age gt 21`,
		`active eq true`,
		`nickName eq null`,
		`userName eq "bjensen" and active eq true`,
		`a eq "1" and b eq "2" or c eq "3"`,
		`not (userName eq "bjensen")`,
		`not (not (userName pr))`,
		`(a eq "1" or b eq "2") and c eq "3"`,
		`emails[type eq "work"]`,
		`emails[type eq "work"].value`,
		`emails[type eq "work" and primary eq true]`,
	}

	g := &Grammar{}
	for _, input := range tt {
		t.Run(input, func(t *testing.T) {
			original, err := g.Parse(input)
			require.NoError(t, err)

			text, err := Visit[string](StringifyVisitor{}, original)
			require.NoError(t, err)

			reparsed, err := g.Parse(text)
			require.NoError(t, err)

			assert.Equal(t, original, reparsed)
		})
	}
}

type countingVisitor struct{}

func (countingVisitor) VisitAnd(left, right int) (int, error)           { return left + right, nil }
func (countingVisitor) VisitOr(left, right int) (int, error)            { return left + right, nil }
func (countingVisitor) VisitNot(operand int) (int, error)               { return operand + 100, nil }
func (countingVisitor) VisitEquals(string, any) (int, error)            { return 1, nil }
func (countingVisitor) VisitNotEquals(string, any) (int, error)         { return 1, nil }
func (countingVisitor) VisitContains(string, any) (int, error)          { return 1, nil }
func (countingVisitor) VisitStartsWith(string, any) (int, error)        { return 1, nil }
func (countingVisitor) VisitEndsWith(string, any) (int, error)          { return 1, nil }
func (countingVisitor) VisitGreaterThan(string, any) (int, error)       { return 1, nil }
func (countingVisitor) VisitGreaterThanEquals(string, any) (int, error) { return 1, nil }
func (countingVisitor) VisitLessThan(string, any) (int, error)          { return 1, nil }
func (countingVisitor) VisitLessThanEquals(string, any) (int, error)    { return 1, nil }
func (countingVisitor) VisitPresence(string) (int, error)               { return 1, nil }
func (countingVisitor) VisitValuePath(_ string, _ string, valueFilter func() (int, error)) (int, error) {
	return valueFilter()
}

func TestVisitDispatch(t *testing.T) {
	g := &Grammar{}

	t.Run("counts leaves across and/or", func(t *testing.T) {
		node, err := g.Parse(`a eq "1" and b eq "2" and c eq "3"`)
		require.NoError(t, err)

		count, err := Visit[int](countingVisitor{}, node)

		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("not adds 100 on top of the underlying result", func(t *testing.T) {
		node, err := g.Parse(`not (userName eq "bjensen")`)
		require.NoError(t, err)

		count, err := Visit[int](countingVisitor{}, node)

		require.NoError(t, err)
		assert.Equal(t, 101, count)
	})

	t.Run("nil node is an error, not a panic", func(t *testing.T) {
		_, err := Visit[int](countingVisitor{}, nil)

		assert.Error(t, err)
	})
}

type scopeVisitor struct{ scope []string }

func (v *scopeVisitor) qualify(attribute string) string {
	if len(v.scope) == 0 {
		return attribute
	}
	return v.scope[len(v.scope)-1] + "." + attribute
}

func (v *scopeVisitor) VisitAnd(l, r string) (string, error) { return l + " AND " + r, nil }
func (v *scopeVisitor) VisitOr(l, r string) (string, error)  { return l + " OR " + r, nil }
func (v *scopeVisitor) VisitNot(o string) (string, error)    { return "NOT " + o, nil }
func (v *scopeVisitor) VisitEquals(a string, _ any) (string, error) {
	return v.qualify(a) + " = ?", nil
}
func (v *scopeVisitor) VisitNotEquals(a string, _ any) (string, error)   { return v.qualify(a), nil }
func (v *scopeVisitor) VisitContains(a string, _ any) (string, error)    { return v.qualify(a), nil }
func (v *scopeVisitor) VisitStartsWith(a string, _ any) (string, error)  { return v.qualify(a), nil }
func (v *scopeVisitor) VisitEndsWith(a string, _ any) (string, error)    { return v.qualify(a), nil }
func (v *scopeVisitor) VisitGreaterThan(a string, _ any) (string, error) { return v.qualify(a), nil }
func (v *scopeVisitor) VisitGreaterThanEquals(a string, _ any) (string, error) {
	return v.qualify(a), nil
}
func (v *scopeVisitor) VisitLessThan(a string, _ any) (string, error)       { return v.qualify(a), nil }
func (v *scopeVisitor) VisitLessThanEquals(a string, _ any) (string, error) { return v.qualify(a), nil }
func (v *scopeVisitor) VisitPresence(a string) (string, error)              { return v.qualify(a), nil }
func (v *scopeVisitor) VisitValuePath(path string, _ string, valueFilter func() (string, error)) (string, error) {
	v.scope = append(v.scope, path)
	inner, err := valueFilter()
	v.scope = v.scope[:len(v.scope)-1]
	if err != nil {
		return "", err
	}
	return "EXISTS(" + inner + ")", nil
}

func TestVisitValuePathThreadsScope(t *testing.T) {
	g := &Grammar{}
	node, err := g.Parse(`emails[type eq "work" and primary eq true]`)
	require.NoError(t, err)

	got, err := Visit[string](&scopeVisitor{}, node)

	require.NoError(t, err)
	assert.Equal(t, "EXISTS(emails.type = ? AND emails.primary = ?)", got)
}
