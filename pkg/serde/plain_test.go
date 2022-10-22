package serde

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nostringer struct {
	Value string
}

type stringer struct {
	Value string
}

func (s stringer) String() string {
	return fmt.Sprintf("The %s", s.Value)
}

func TestToPlain(t *testing.T) {

	t.Run("stringafies an item that doesn't implement fmt.Stringer", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, ToPlain[nostringer](w, nostringer{Value: "example"}))
		assert.Equal(t, `{example}`, w.String())
	})

	t.Run("stringafies an item that implements fmt.Stringer", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		require.NoError(t, ToPlain[stringer](w, stringer{Value: "example"}))
		assert.Equal(t, `The example`, w.String())
	})
}
