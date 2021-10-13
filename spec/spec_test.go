package spec

import (
	"testing"
)

func TestCalculator(t *testing.T) {
	t.Run("when adding 1 + 1", func(t *testing.T) {
		Establish(t, func(x *Context) {
			x.Because(func() interface{} { return 1 + 1 })

			x.It(func(t *testing.T) {
				if x.Result != 2 {
					t.Errorf("Expected: 2, Got: %d", x.Result)
				}
			})
		})
	})
}
