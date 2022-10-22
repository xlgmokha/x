package serde

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type example struct {
	ID   string `jsonapi:"primary,examples"`
	Name string `jsonapi:"attr,name"`
}

func TestToJSONAPI(t *testing.T) {
	t.Run("serializes a custom type", func(t *testing.T) {
		io := bytes.NewBuffer(nil)

		require.NoError(t, ToJSONAPI(io, &example{Name: "Example"}))

		assert.Equal(t, `{"data":{"type":"examples","attributes":{"name":"Example"}}}`+"\n", io.String())
	})

	t.Run("serializes a jsonapi.ErrorsPayload", func(t *testing.T) {
		io := bytes.NewBuffer(nil)

		require.NoError(t, ToJSONAPI(io, &jsonapi.ErrorsPayload{
			Errors: []*jsonapi.ErrorObject{
				{
					ID:     "id",
					Title:  "Name is required",
					Status: "400",
				},
			},
		}))

		assert.Equal(t, `{"errors":[{"id":"id","title":"Name is required","status":"400"}]}`+"\n", io.String())
	})
}

func TestFromJSONAPI(t *testing.T) {
	t.Run("from a single item", func(t *testing.T) {
		io := strings.NewReader(`{"data":{"type":"examples","id":"42","attributes":{"name":"Example"}}}`)

		item, err := FromJSONAPI[example](io)

		require.NoError(t, err)
		assert.Equal(t, "42", item.ID)
		assert.Equal(t, "Example", item.Name)
	})
}
