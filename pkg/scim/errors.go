package scim

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidFilter = errors.New("scim: invalid filter")
	ErrInputTooLarge = errors.New("scim: filter exceeds MaxInputBytes")
)

// ParseError reports a filter that does not conform to the RFC 7644 grammar,
// with the byte offset where parsing stopped. It unwraps to ErrInvalidFilter.
type ParseError struct {
	Input    string
	Position int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("scim: invalid filter at position %d: %q", e.Position, e.Input)
}

func (e *ParseError) Unwrap() error { return ErrInvalidFilter }
