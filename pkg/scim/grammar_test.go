package scim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xlgmokha/x/pkg/peg"
)

func TestGrammarAttributePath(t *testing.T) {
	tt := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "userName eq", "userName"},
		{"dotted sub-attribute", "name.familyName eq", "name.familyName"},
		{"schema-URI-prefixed", `urn:ietf:params:scim:schemas:core:2.0:User:userName eq`, "urn:ietf:params:scim:schemas:core:2.0:User:userName"},
		{"schema-URI-prefixed with sub-attribute", `urn:ietf:params:scim:schemas:core:2.0:User:name.familyName eq`, "urn:ietf:params:scim:schemas:core:2.0:User:name.familyName"},
	}
	g := &Grammar{}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, ok := g.AttributePath()(peg.NewContext(tc.input))
			require.True(t, ok)
			assert.Equal(t, tc.want, result)
		})
	}
}

func TestGrammarLeafComparison(t *testing.T) {
	g := &Grammar{}
	node, ok := g.Parse(`userName eq "bjensen"`)

	require.True(t, ok)
	require.NotNil(t, node)
	assert.Equal(t, "userName", node.Attribute())
	assert.Equal(t, "eq", node.Operator())
	assert.Equal(t, "bjensen", node.Value())
}

func TestGrammarComparisonOperatorsAndLiterals(t *testing.T) {
	tt := []struct {
		input     string
		attribute string
		operator  string
		value     any
	}{
		{`userName eq "bjensen"`, "userName", "eq", "bjensen"},
		{`userName ne "bjensen"`, "userName", "ne", "bjensen"},
		{`userName co "jen"`, "userName", "co", "jen"},
		{`userName sw "bjen"`, "userName", "sw", "bjen"},
		{`userName ew "sen"`, "userName", "ew", "sen"},
		{`age gt 21`, "age", "gt", 21.0},
		{`age ge 21`, "age", "ge", 21.0},
		{`age lt 21`, "age", "lt", 21.0},
		{`age le 21`, "age", "le", 21.0},
		{`salary ge 1.5e4`, "salary", "ge", 1.5e4},
		{`salary ge -3.5`, "salary", "ge", -3.5},
		{`active eq true`, "active", "eq", true},
		{`active eq false`, "active", "eq", false},
		{`nickName eq null`, "nickName", "eq", nil},
	}
	g := &Grammar{}
	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			node, ok := g.Parse(tc.input)

			require.True(t, ok)
			require.NotNil(t, node)
			assert.Equal(t, tc.attribute, node.Attribute())
			assert.Equal(t, tc.operator, node.Operator())
			assert.Equal(t, tc.value, node.Value())
		})
	}
}

func TestGrammarPresence(t *testing.T) {
	g := &Grammar{}
	node, ok := g.Parse("userName pr")

	require.True(t, ok)
	require.NotNil(t, node)
	assert.Equal(t, "userName", node.Attribute())
	assert.Equal(t, "pr", node.Operator())
	assert.Nil(t, node.Value())
}

func TestGrammarLogicalExpression(t *testing.T) {
	g := &Grammar{}

	t.Run("and", func(t *testing.T) {
		node, ok := g.Parse(`userName eq "bjensen" and active eq true`)

		require.True(t, ok)
		assert.Equal(t, "and", node.Operator())
		assert.Equal(t, "userName", node.Left().Attribute())
		assert.Equal(t, "active", node.Right().Attribute())
	})

	t.Run("or", func(t *testing.T) {
		node, ok := g.Parse(`userName eq "bjensen" or userName eq "jsmith"`)

		require.True(t, ok)
		assert.Equal(t, "or", node.Operator())
		assert.Equal(t, "bjensen", node.Left().Value())
		assert.Equal(t, "jsmith", node.Right().Value())
	})

	t.Run("chained and", func(t *testing.T) {
		node, ok := g.Parse(`a eq "1" and b eq "2" and c eq "3"`)

		require.True(t, ok)
		assert.Equal(t, "and", node.Operator())
		assert.Equal(t, "a", node.Left().Attribute())
		assert.Equal(t, "and", node.Right().Operator())
		assert.Equal(t, "b", node.Right().Left().Attribute())
		assert.Equal(t, "c", node.Right().Right().Attribute())
	})

	t.Run("and binds tighter than or, regardless of textual order", func(t *testing.T) {
		node, ok := g.Parse(`a eq "1" and b eq "2" or c eq "3"`)

		require.True(t, ok)
		require.Equal(t, "or", node.Operator(), "expected (a and b) or c, i.e. root operator is or")
		assert.Equal(t, "and", node.Left().Operator())
		assert.Equal(t, "a", node.Left().Left().Attribute())
		assert.Equal(t, "b", node.Left().Right().Attribute())
		assert.Equal(t, "c", node.Right().Attribute())
	})

	t.Run("or on the left, and on the right, and still binds tighter", func(t *testing.T) {
		node, ok := g.Parse(`a eq "1" or b eq "2" and c eq "3"`)

		require.True(t, ok)
		require.Equal(t, "or", node.Operator(), "expected a or (b and c)")
		assert.Equal(t, "a", node.Left().Attribute())
		assert.Equal(t, "and", node.Right().Operator())
		assert.Equal(t, "b", node.Right().Left().Attribute())
		assert.Equal(t, "c", node.Right().Right().Attribute())
	})
}

func TestGrammarParensAndNot(t *testing.T) {
	g := &Grammar{}

	t.Run("not wrapping a parenthesized leaf", func(t *testing.T) {
		node, ok := g.Parse(`not (userName eq "bjensen")`)

		require.True(t, ok)
		assert.True(t, node.Not())
		assert.Equal(t, "userName", node.Operand().Attribute())
		assert.Equal(t, "eq", node.Operand().Operator())
	})

	t.Run("nested not stays distinct at each level", func(t *testing.T) {
		node, ok := g.Parse(`not (not (userName pr))`)

		require.True(t, ok)
		require.True(t, node.Not())
		require.True(t, node.Operand().Not())
		assert.Equal(t, "userName", node.Operand().Operand().Attribute())
		assert.Equal(t, "pr", node.Operand().Operand().Operator())
	})

	t.Run("parens override precedence", func(t *testing.T) {
		node, ok := g.Parse(`(a eq "1" or b eq "2") and c eq "3"`)

		require.True(t, ok)
		require.Equal(t, "and", node.Operator())
		assert.Equal(t, "or", node.Left().Operator())
		assert.Equal(t, "c", node.Right().Attribute())
	})

	t.Run("not without parens is rejected", func(t *testing.T) {
		_, ok := g.Parse(`not userName pr`)

		assert.False(t, ok)
	})
}

func TestGrammarValuePath(t *testing.T) {
	g := &Grammar{}

	t.Run("simple", func(t *testing.T) {
		node, ok := g.Parse(`emails[type eq "work"]`)

		require.True(t, ok)
		require.True(t, node.HasPath())
		assert.Equal(t, "emails", node.Path())
		assert.Equal(t, "", node.SubAttribute())
		assert.Equal(t, "type", node.ValueFilter().Attribute())
		assert.Equal(t, "work", node.ValueFilter().Value())
	})

	t.Run("with a trailing sub-attribute", func(t *testing.T) {
		node, ok := g.Parse(`emails[type eq "work"].value`)

		require.True(t, ok)
		assert.Equal(t, "emails", node.Path())
		assert.Equal(t, "value", node.SubAttribute())
	})

	t.Run("value filter supports and/or", func(t *testing.T) {
		node, ok := g.Parse(`emails[type eq "work" and primary eq true]`)

		require.True(t, ok)
		vf := node.ValueFilter()
		assert.Equal(t, "and", vf.Operator())
		assert.Equal(t, "type", vf.Left().Attribute())
		assert.Equal(t, "primary", vf.Right().Attribute())
	})
}

func TestGrammarRejectsTrailingGarbage(t *testing.T) {
	g := &Grammar{}
	_, ok := g.Parse(`userName eq "bjensen" garbage`)

	assert.False(t, ok)
}

func TestGrammarTrailingWhitespaceIsTolerated(t *testing.T) {
	g := &Grammar{}
	node, ok := g.Parse(`userName eq "bjensen"   `)

	require.True(t, ok)
	assert.Equal(t, "userName", node.Attribute())
}
