package scim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttrPathSplitsParts(t *testing.T) {
	tt := []struct {
		input string
		want  AttrPath
	}{
		{`userName eq "x"`, AttrPath{Name: "userName"}},
		{`name.familyName eq "x"`, AttrPath{Name: "name", Sub: "familyName"}},
		{
			`urn:ietf:params:scim:schemas:core:2.0:User:userName eq "x"`,
			AttrPath{URI: "urn:ietf:params:scim:schemas:core:2.0:User", Name: "userName"},
		},
		{
			`urn:ietf:params:scim:schemas:core:2.0:User:name.familyName eq "x"`,
			AttrPath{URI: "urn:ietf:params:scim:schemas:core:2.0:User", Name: "name", Sub: "familyName"},
		},
	}
	g := &Grammar{}
	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			expr, err := g.Parse(tc.input)
			require.NoError(t, err)

			got := expr.(Comparison).Attribute

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAttrPathOfValuePath(t *testing.T) {
	g := &Grammar{}
	expr, err := g.Parse(`emails[type eq "work"]`)
	require.NoError(t, err)

	assert.Equal(t, AttrPath{Name: "emails"}, expr.(ValuePath).Attribute)
}
