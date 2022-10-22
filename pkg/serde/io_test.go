package serde

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrom(t *testing.T) {
	type Example struct {
		Key   string `json:"key" jsonapi:"primary,examples"`
		Value string `json:"value" jsonapi:"attr,value"`
	}

	t.Run("parses a single item from JSON data", func(t *testing.T) {
		body := strings.NewReader(`{"key":"my-key","value":"my-value"}`)

		result, err := From[Example](body, JSON)

		require.NoError(t, err)
		assert.Equal(t, "my-key", result.Key)
		assert.Equal(t, "my-value", result.Value)
	})

	t.Run("parses a single *item from JSON data", func(t *testing.T) {
		body := strings.NewReader(`{"key":"my-key","value":"my-value"}`)

		result, err := From[*Example](body, JSON)

		require.NoError(t, err)
		assert.Equal(t, "my-key", result.Key)
		assert.Equal(t, "my-value", result.Value)
	})

	t.Run("parses a slices of items from JSON data", func(t *testing.T) {
		body := strings.NewReader(`[{"key":"my-key","value":"my-value"}]`)

		results, err := From[[]Example](body, JSON)

		require.NoError(t, err)
		require.Equal(t, 1, len(results))
		assert.Equal(t, "my-key", results[0].Key)
		assert.Equal(t, "my-value", results[0].Value)
	})

	t.Run("parses a slices of *items from JSON data", func(t *testing.T) {
		body := strings.NewReader(`[{"key":"my-key","value":"my-value"}]`)

		results, err := From[[]*Example](body, JSON)

		require.NoError(t, err)
		require.Equal(t, 1, len(results))
		assert.Equal(t, "my-key", results[0].Key)
		assert.Equal(t, "my-value", results[0].Value)
	})

	t.Run("parses a single item from JSON API data", func(t *testing.T) {
		body := strings.NewReader(`{"data":{"type":"examples","id":"my-key","attributes":{"value":"my-value"}}}`)

		result, err := From[Example](body, JSONAPI)

		require.NoError(t, err)
		assert.Equal(t, "my-key", result.Key)
		assert.Equal(t, "my-value", result.Value)
	})

	t.Run("parses a single *item from JSON API data", func(t *testing.T) {
		body := strings.NewReader(`{"data":{"type":"examples","id":"my-key","attributes":{"value":"my-value"}}}`)

		result, err := From[*Example](body, JSONAPI)

		require.NoError(t, err)
		assert.Equal(t, "my-key", result.Key)
		assert.Equal(t, "my-value", result.Value)
	})

	t.Run("parses a slice of items from JSON API data", func(t *testing.T) {
		t.Skip()
		body := strings.NewReader(`{"data":[{"type":"examples","id":"my-key","attributes":{"value":"my-value"}}]}`)

		results, err := From[[]Example](body, JSONAPI)

		require.NoError(t, err)
		require.Equal(t, 1, len(results))
		assert.Equal(t, "my-key", results[0].Key)
		assert.Equal(t, "my-value", results[0].Value)
	})

	t.Run("parses a slice of *items from JSON API data", func(t *testing.T) {
		body := strings.NewReader(`{"data":[{"type":"examples","id":"my-key","attributes":{"value":"my-value"}}]}`)

		results, err := From[[]*Example](body, JSONAPI)

		require.NoError(t, err)
		require.Equal(t, 1, len(results))
		assert.Equal(t, "my-key", results[0].Key)
		assert.Equal(t, "my-value", results[0].Value)
	})
}

func TestTo(t *testing.T) {
	type Example struct {
		Key   string `json:"key" jsonapi:"primary,examples"`
		Value string `json:"value" jsonapi:"attr,value"`
	}

	t.Run("serializes an item to JSON", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, To[Example](w, Example{
			Key:   "my-key",
			Value: "my-value",
		}, JSON))
		expected := `{
  "key": "my-key",
  "value": "my-value"
}
`
		assert.Equal(t, expected, w.String())
	})

	t.Run("serializes an *item to JSON", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, To[*Example](w, &Example{
			Key:   "my-key",
			Value: "my-value",
		}, JSON))
		expected := `{
  "key": "my-key",
  "value": "my-value"
}
`
		assert.Equal(t, expected, w.String())
	})

	t.Run("serializes items to JSON", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, To(w, []Example{{
			Key:   "my-key",
			Value: "my-value",
		}}, JSON))
		expected := `[
  {
    "key": "my-key",
    "value": "my-value"
  }
]
`
		assert.Equal(t, expected, w.String())
	})

	t.Run("serializes *items to JSON", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, To(w, []*Example{{
			Key:   "my-key",
			Value: "my-value",
		}}, JSON))
		expected := `[
  {
    "key": "my-key",
    "value": "my-value"
  }
]
`
		assert.Equal(t, expected, w.String())
	})

	t.Run("serializes an item to JSON API", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, To[Example](w, Example{
			Key:   "my-key",
			Value: "my-value",
		}, JSONAPI))
		expected := `{"data":{"type":"examples","id":"my-key","attributes":{"value":"my-value"}}}
`
		assert.Equal(t, expected, w.String())
	})

	t.Run("serializes an *item to JSON API", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, To[*Example](w, &Example{
			Key:   "my-key",
			Value: "my-value",
		}, JSONAPI))
		expected := `{"data":{"type":"examples","id":"my-key","attributes":{"value":"my-value"}}}
`
		assert.Equal(t, expected, w.String())
	})

	t.Run("serializes a slice of items to JSON API", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, To[[]Example](w, []Example{{
			Key:   "my-key",
			Value: "my-value",
		}}, JSONAPI))
		expected := `{"data":[{"type":"examples","id":"my-key","attributes":{"value":"my-value"}}]}
`
		assert.Equal(t, expected, w.String())
	})

	t.Run("serializes a slice of *items to JSON API", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, To[[]*Example](w, []*Example{{
			Key:   "my-key",
			Value: "my-value",
		}}, JSONAPI))
		expected := `{"data":[{"type":"examples","id":"my-key","attributes":{"value":"my-value"}}]}
`
		assert.Equal(t, expected, w.String())
	})
}
