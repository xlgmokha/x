package serde

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromHTTP(t *testing.T) {
	type input struct {
		ID   string `json:"id" jsonapi:"primary,inputs"`
		Name string `json:"name" jsonapi:"attr,name"`
	}

	t.Run("parses a JSON request body by default", func(t *testing.T) {
		body := strings.NewReader(`{"id":"1","name": "Faux"}`)
		r := httptest.NewRequest("GET", "/", body)

		item, err := FromHTTP[input](r)

		require.NoError(t, err)
		assert.NotZero(t, item)
		assert.Equal(t, "1", item.ID)
		assert.Equal(t, "Faux", item.Name)
	})

	t.Run("parses the request body as JSON", func(t *testing.T) {
		body := strings.NewReader(`{"id":"1","name": "Faux"}`)
		r := httptest.NewRequest("GET", "/", body)
		r.Header.Set("Content-Type", "application/json")

		item, err := FromHTTP[input](r)

		require.NoError(t, err)
		assert.NotZero(t, item)
		assert.Equal(t, "1", item.ID)
		assert.Equal(t, "Faux", item.Name)
	})

	t.Run("parses a JSON API request body", func(t *testing.T) {
		body := strings.NewReader(`{"data":{"type":"inputs","id":"1","attributes":{"name": "Faux"}}}`)
		r := httptest.NewRequest("GET", "/", body)
		r.Header.Set("Content-Type", jsonapi.MediaType)

		item, err := FromHTTP[input](r)

		require.NoError(t, err)
		assert.NotZero(t, item)
		assert.Equal(t, "1", item.ID)
		assert.Equal(t, "Faux", item.Name)
	})
}

func TestToHTTP(t *testing.T) {
	type output struct {
		ID   string `json:"id" jsonapi:"primary,inputs"`
		Name string `json:"name" jsonapi:"attr,name"`
	}

	t.Run("serializes the data as JSON by default", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		require.NoError(t, ToHTTP[*output](w, r, &output{
			ID:   "2",
			Name: "Dave East",
		}))

		result, err := From[output](w.Body, JSON)
		require.NoError(t, err)
		assert.Equal(t, "2", result.ID)
		assert.Equal(t, "Dave East", result.Name)
	})

	t.Run("serializes the data as JSON when specified in the request Accept header", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		require.NoError(t, ToHTTP[*output](w, r, &output{
			ID:   "2",
			Name: "Dave East",
		}))

		result, err := From[output](w.Body, JSON)
		require.NoError(t, err)
		assert.Equal(t, "2", result.ID)
		assert.Equal(t, "Dave East", result.Name)
	})

	t.Run("serializes the data as JSON API when specified in the request Accept header", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept", jsonapi.MediaType)
		w := httptest.NewRecorder()

		require.NoError(t, ToHTTP[*output](w, r, &output{
			ID:   "2",
			Name: "Dave East",
		}))

		result, err := From[output](w.Body, JSONAPI)
		require.NoError(t, err)
		assert.Equal(t, "2", result.ID)
		assert.Equal(t, "Dave East", result.Name)
	})
}
