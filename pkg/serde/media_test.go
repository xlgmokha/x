package serde

import (
	"fmt"
	"testing"

	"github.com/google/jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestMediaTypeFor(t *testing.T) {
	tt := []struct {
		input    string
		expected MediaType
	}{
		{input: jsonapi.MediaType, expected: JSONAPI},
		{input: "application/json", expected: JSON},
		{input: "text/plain", expected: Text},
		{input: "text/html", expected: Default},
		{input: "*/*", expected: Default},
		{input: "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8", expected: Default},
		{input: "application/vnd.api+json, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8", expected: JSONAPI},
		{input: "application/vnd.api+json;q=0.1, application/json;q=0.2, text/plain;q=0.3, unknown/thing;q=0.4", expected: Text},
		{input: "text/html; charset=UTF-8", expected: Default},
		{input: "application/json; charset=UTF-8", expected: JSON},
		{input: "application/vnd.api+json; charset=UTF-8", expected: JSONAPI},
		{input: "text/plain; charset=UTF-8", expected: Text},
		{input: "application/json, text/plain, */*", expected: JSON},
	}
	for _, row := range tt {
		t.Run(fmt.Sprintf("%v returns %v", row.input, row.expected), func(t *testing.T) {
			assert.Equal(t, row.expected, MediaTypeFor(row.input))
		})
	}
}
