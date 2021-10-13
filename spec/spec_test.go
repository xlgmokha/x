package spec

import (
	"testing"
)

func TestCalculator(t *testing.T) {
	t.Run("when adding 1 + 1", func(t *testing.T) {
		Establish(func(x *Context) {
			x.Because(func() { x.Set("result", 1+1).Set("2+2", 2+2) })

			x.It(func() {
				result := x.Get("result").(int)
				if result != 2 {
					t.Errorf("Expected: 2, Got: %d", result)
				}
			})
		})
	})
}
