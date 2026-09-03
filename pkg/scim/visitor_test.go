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
			original, ok := g.Parse(input)
			require.True(t, ok)

			text, err := Visit[string](StringifyVisitor{}, original)
			require.NoError(t, err)

			reparsed, ok := g.Parse(text)
			require.True(t, ok)

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
func (countingVisitor) VisitValuePath(string, int, string) (int, error) { return 1, nil }

func TestVisitDispatch(t *testing.T) {
	g := &Grammar{}

	t.Run("counts leaves across and/or", func(t *testing.T) {
		node, ok := g.Parse(`a eq "1" and b eq "2" and c eq "3"`)
		require.True(t, ok)

		count, err := Visit[int](countingVisitor{}, node)

		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("not adds 100 on top of the underlying result", func(t *testing.T) {
		node, ok := g.Parse(`not (userName eq "bjensen")`)
		require.True(t, ok)

		count, err := Visit[int](countingVisitor{}, node)

		require.NoError(t, err)
		assert.Equal(t, 101, count)
	})

	t.Run("nil node is an error, not a panic", func(t *testing.T) {
		_, err := Visit[int](countingVisitor{}, nil)

		assert.Error(t, err)
	})
}
