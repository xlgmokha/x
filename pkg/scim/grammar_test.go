package scim

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cmp(t *testing.T, e Expression) Comparison {
	t.Helper()
	c, ok := e.(Comparison)
	require.Truef(t, ok, "expected Comparison, got %T", e)
	return c
}

func logical(t *testing.T, e Expression) Logical {
	t.Helper()
	l, ok := e.(Logical)
	require.Truef(t, ok, "expected Logical, got %T", e)
	return l
}

func notExpr(t *testing.T, e Expression) Not {
	t.Helper()
	n, ok := e.(Not)
	require.Truef(t, ok, "expected Not, got %T", e)
	return n
}

func valuePath(t *testing.T, e Expression) ValuePath {
	t.Helper()
	v, ok := e.(ValuePath)
	require.Truef(t, ok, "expected ValuePath, got %T", e)
	return v
}

func TestGrammarAttributePath(t *testing.T) {
	tt := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", `userName eq "x"`, "userName"},
		{"dotted sub-attribute", `name.familyName eq "x"`, "name.familyName"},
		{"schema-URI-prefixed", `urn:ietf:params:scim:schemas:core:2.0:User:userName eq "x"`, "urn:ietf:params:scim:schemas:core:2.0:User:userName"},
		{"schema-URI-prefixed with sub-attribute", `urn:ietf:params:scim:schemas:core:2.0:User:name.familyName eq "x"`, "urn:ietf:params:scim:schemas:core:2.0:User:name.familyName"},
	}
	g := &Grammar{}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := g.Parse(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.want, cmp(t, expr).Attribute.String())
		})
	}
}

func TestGrammarLeafComparison(t *testing.T) {
	g := &Grammar{}
	expr, err := g.Parse(`userName eq "bjensen"`)

	require.NoError(t, err)
	c := cmp(t, expr)
	assert.Equal(t, "userName", c.Attribute.String())
	assert.Equal(t, Equal, c.Operator)
	assert.Equal(t, "bjensen", c.Value)
}

func TestGrammarComparisonOperatorsAndLiterals(t *testing.T) {
	tt := []struct {
		input     string
		attribute string
		operator  CompareOperator
		value     any
	}{
		{`userName eq "bjensen"`, "userName", Equal, "bjensen"},
		{`userName ne "bjensen"`, "userName", NotEqual, "bjensen"},
		{`userName co "jen"`, "userName", Contains, "jen"},
		{`userName sw "bjen"`, "userName", StartsWith, "bjen"},
		{`userName ew "sen"`, "userName", EndsWith, "sen"},
		{`age gt 21`, "age", GreaterThan, 21.0},
		{`age ge 21`, "age", GreaterOrEqual, 21.0},
		{`age lt 21`, "age", LessThan, 21.0},
		{`age le 21`, "age", LessOrEqual, 21.0},
		{`salary ge 1.5e4`, "salary", GreaterOrEqual, 1.5e4},
		{`salary ge -3.5`, "salary", GreaterOrEqual, -3.5},
		{`active eq true`, "active", Equal, true},
		{`active eq false`, "active", Equal, false},
		{`nickName eq null`, "nickName", Equal, nil},
	}
	g := &Grammar{}
	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			expr, err := g.Parse(tc.input)

			require.NoError(t, err)
			c := cmp(t, expr)
			assert.Equal(t, tc.attribute, c.Attribute.String())
			assert.Equal(t, tc.operator, c.Operator)
			assert.Equal(t, tc.value, c.Value)
		})
	}
}

func TestGrammarPresence(t *testing.T) {
	g := &Grammar{}
	expr, err := g.Parse("userName pr")

	require.NoError(t, err)
	c := cmp(t, expr)
	assert.Equal(t, "userName", c.Attribute.String())
	assert.Equal(t, Present, c.Operator)
	assert.Nil(t, c.Value)
}

func TestGrammarLogicalExpression(t *testing.T) {
	g := &Grammar{}

	t.Run("and", func(t *testing.T) {
		expr, err := g.Parse(`userName eq "bjensen" and active eq true`)

		require.NoError(t, err)
		l := logical(t, expr)
		assert.Equal(t, And, l.Operator)
		assert.Equal(t, "userName", cmp(t, l.Left).Attribute.String())
		assert.Equal(t, "active", cmp(t, l.Right).Attribute.String())
	})

	t.Run("or", func(t *testing.T) {
		expr, err := g.Parse(`userName eq "bjensen" or userName eq "jsmith"`)

		require.NoError(t, err)
		l := logical(t, expr)
		assert.Equal(t, Or, l.Operator)
		assert.Equal(t, "bjensen", cmp(t, l.Left).Value)
		assert.Equal(t, "jsmith", cmp(t, l.Right).Value)
	})

	t.Run("chained and", func(t *testing.T) {
		expr, err := g.Parse(`a eq "1" and b eq "2" and c eq "3"`)

		require.NoError(t, err)
		l := logical(t, expr)
		assert.Equal(t, And, l.Operator)
		assert.Equal(t, "a", cmp(t, l.Left).Attribute.String())
		r := logical(t, l.Right)
		assert.Equal(t, And, r.Operator)
		assert.Equal(t, "b", cmp(t, r.Left).Attribute.String())
		assert.Equal(t, "c", cmp(t, r.Right).Attribute.String())
	})

	t.Run("and binds tighter than or, regardless of textual order", func(t *testing.T) {
		expr, err := g.Parse(`a eq "1" and b eq "2" or c eq "3"`)

		require.NoError(t, err)
		l := logical(t, expr)
		require.Equal(t, Or, l.Operator, "expected (a and b) or c, i.e. root operator is or")
		left := logical(t, l.Left)
		assert.Equal(t, And, left.Operator)
		assert.Equal(t, "a", cmp(t, left.Left).Attribute.String())
		assert.Equal(t, "b", cmp(t, left.Right).Attribute.String())
		assert.Equal(t, "c", cmp(t, l.Right).Attribute.String())
	})

	t.Run("or on the left, and on the right, and still binds tighter", func(t *testing.T) {
		expr, err := g.Parse(`a eq "1" or b eq "2" and c eq "3"`)

		require.NoError(t, err)
		l := logical(t, expr)
		require.Equal(t, Or, l.Operator, "expected a or (b and c)")
		assert.Equal(t, "a", cmp(t, l.Left).Attribute.String())
		r := logical(t, l.Right)
		assert.Equal(t, And, r.Operator)
		assert.Equal(t, "b", cmp(t, r.Left).Attribute.String())
		assert.Equal(t, "c", cmp(t, r.Right).Attribute.String())
	})
}

func TestGrammarParensAndNot(t *testing.T) {
	g := &Grammar{}

	t.Run("not wrapping a parenthesized leaf", func(t *testing.T) {
		expr, err := g.Parse(`not (userName eq "bjensen")`)

		require.NoError(t, err)
		c := cmp(t, notExpr(t, expr).Operand)
		assert.Equal(t, "userName", c.Attribute.String())
		assert.Equal(t, Equal, c.Operator)
	})

	t.Run("nested not stays distinct at each level", func(t *testing.T) {
		expr, err := g.Parse(`not (not (userName pr))`)

		require.NoError(t, err)
		inner := notExpr(t, notExpr(t, expr).Operand)
		c := cmp(t, inner.Operand)
		assert.Equal(t, "userName", c.Attribute.String())
		assert.Equal(t, Present, c.Operator)
	})

	t.Run("parens override precedence", func(t *testing.T) {
		expr, err := g.Parse(`(a eq "1" or b eq "2") and c eq "3"`)

		require.NoError(t, err)
		l := logical(t, expr)
		require.Equal(t, And, l.Operator)
		assert.Equal(t, Or, logical(t, l.Left).Operator)
		assert.Equal(t, "c", cmp(t, l.Right).Attribute.String())
	})

	t.Run("not without parens is rejected", func(t *testing.T) {
		_, err := g.Parse(`not userName pr`)

		assert.Error(t, err)
	})
}

func TestGrammarValuePath(t *testing.T) {
	g := &Grammar{}

	t.Run("simple", func(t *testing.T) {
		expr, err := g.Parse(`emails[type eq "work"]`)

		require.NoError(t, err)
		vp := valuePath(t, expr)
		assert.Equal(t, "emails", vp.Attribute.String())
		assert.Equal(t, "", vp.SubAttribute)
		f := cmp(t, vp.Filter)
		assert.Equal(t, "type", f.Attribute.String())
		assert.Equal(t, "work", f.Value)
	})

	t.Run("with a trailing sub-attribute", func(t *testing.T) {
		expr, err := g.Parse(`emails[type eq "work"].value`)

		require.NoError(t, err)
		vp := valuePath(t, expr)
		assert.Equal(t, "emails", vp.Attribute.String())
		assert.Equal(t, "value", vp.SubAttribute)
	})

	t.Run("value filter supports and/or", func(t *testing.T) {
		expr, err := g.Parse(`emails[type eq "work" and primary eq true]`)

		require.NoError(t, err)
		l := logical(t, valuePath(t, expr).Filter)
		assert.Equal(t, And, l.Operator)
		assert.Equal(t, "type", cmp(t, l.Left).Attribute.String())
		assert.Equal(t, "primary", cmp(t, l.Right).Attribute.String())
	})
}

func TestGrammarRejectsTrailingGarbage(t *testing.T) {
	g := &Grammar{}
	_, err := g.Parse(`userName eq "bjensen" garbage`)

	assert.Error(t, err)
}

func TestParseRejectsOversizedInput(t *testing.T) {
	g := &Grammar{MaxInputBytes: 16}

	t.Run("input over the cap is rejected before parsing", func(t *testing.T) {
		_, err := g.Parse(`userName eq "bjensen"`) // 21 bytes > 16
		require.ErrorIs(t, err, ErrInputTooLarge)
		assert.NotErrorIs(t, err, ErrInvalidFilter)
	})

	t.Run("input within the cap still parses", func(t *testing.T) {
		expr, err := g.Parse(`a eq "1"`) // 8 bytes <= 16
		require.NoError(t, err)
		assert.Equal(t, "a", cmp(t, expr).Attribute.String())
	})
}

func TestParseDefaultIsUnbounded(t *testing.T) {
	g := New() // MaxInputBytes defaults to 0
	const depth = 4000
	input := strings.Repeat("(", depth) + `a eq "1"` + strings.Repeat(")", depth)

	expr, err := g.Parse(input)

	require.NoError(t, err)
	assert.Equal(t, "a", cmp(t, expr).Attribute.String())
}

func TestParseErrorReportsPosition(t *testing.T) {
	g := &Grammar{}
	input := `userName eq "bjensen" garbage`
	_, err := g.Parse(input)

	require.ErrorIs(t, err, ErrInvalidFilter)
	var perr *ParseError
	require.ErrorAs(t, err, &perr)
	assert.Equal(t, input, perr.Input)
	assert.Equal(t, len("userName eq \"bjensen\" "), perr.Position)
}

func TestGrammarTrailingWhitespaceIsTolerated(t *testing.T) {
	g := &Grammar{}
	expr, err := g.Parse(`userName eq "bjensen"   `)

	require.NoError(t, err)
	assert.Equal(t, "userName", cmp(t, expr).Attribute.String())
}

func TestGrammarDeepNestingIsLinear(t *testing.T) {
	g := &Grammar{}
	const depth = 2000
	input := strings.Repeat("(", depth) + `userName eq "bjensen"` + strings.Repeat(")", depth)

	expr, err := g.Parse(input)

	require.NoError(t, err)
	c := cmp(t, expr)
	assert.Equal(t, "userName", c.Attribute.String())
	assert.Equal(t, Equal, c.Operator)
}

func TestGrammarOperatorsAreCaseInsensitive(t *testing.T) {
	// RFC 7644 3.4.2.2: attribute operators and logical keywords are case insensitive.
	g := &Grammar{}

	t.Run("comparison operator", func(t *testing.T) {
		expr, err := g.Parse(`userName EQ "bjensen"`)
		require.NoError(t, err)
		assert.Equal(t, Equal, cmp(t, expr).Operator)
	})

	t.Run("logical keyword", func(t *testing.T) {
		expr, err := g.Parse(`userName eq "a" AND active pr`)
		require.NoError(t, err)
		assert.Equal(t, And, logical(t, expr).Operator)
	})

	t.Run("presence", func(t *testing.T) {
		expr, err := g.Parse(`userName PR`)
		require.NoError(t, err)
		assert.Equal(t, Present, cmp(t, expr).Operator)
	})
}

func TestGrammarJSONLiteralsAreCaseSensitive(t *testing.T) {
	// compValue uses JSON rules (RFC 7159); true/false/null are lowercase only.
	g := &Grammar{}

	for _, good := range []string{`active eq true`, `active eq false`, `nickName eq null`} {
		t.Run("accepts "+good, func(t *testing.T) {
			_, err := g.Parse(good)
			assert.NoError(t, err)
		})
	}
	for _, bad := range []string{`active eq TRUE`, `active eq False`, `nickName eq Null`} {
		t.Run("rejects "+bad, func(t *testing.T) {
			_, err := g.Parse(bad)
			assert.Error(t, err)
		})
	}
}

func TestGrammarStringEscapesAreDecoded(t *testing.T) {
	// compValue strings follow JSON rules, so escapes must be decoded.
	tt := []struct {
		input string
		want  string
	}{
		{`userName eq "a\"b"`, `a"b`},
		{`userName eq "a\\b"`, `a\b`},
		{`userName eq "tab\tend"`, "tab\tend"},
		{`userName eq "A"`, "A"},
	}
	g := &Grammar{}
	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			expr, err := g.Parse(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.want, cmp(t, expr).Value)
		})
	}
}

func TestGrammarRejectsMalformed(t *testing.T) {
	tt := []struct {
		name  string
		input string
	}{
		{"dangling and", `userName eq "a" and`},
		{"dangling or", `userName eq "a" or`},
		{"invalid string escape", `userName eq "\x"`},
		{"number out of range", `age gt 1e400`},
		{"unclosed paren", `not (userName eq "a"`},
		{"unopened paren", `userName eq "a")`},
	}
	g := &Grammar{}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := g.Parse(tc.input)
			assert.ErrorIs(t, err, ErrInvalidFilter)
		})
	}
}

func TestGrammarValueFilterRejectsNestedValuePath(t *testing.T) {
	// valFilter = attrExp / logExp / *1"not" "(" valFilter ")" -- no nested valuePath.
	g := &Grammar{}

	t.Run("nested value path is rejected", func(t *testing.T) {
		_, err := g.Parse(`emails[members[type eq "work"]]`)
		assert.Error(t, err)
	})

	t.Run("not-group inside brackets is accepted", func(t *testing.T) {
		expr, err := g.Parse(`emails[not (type eq "work")]`)
		require.NoError(t, err)
		_, ok := valuePath(t, expr).Filter.(Not)
		assert.True(t, ok)
	})
}
