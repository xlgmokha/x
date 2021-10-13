package spec

import "testing"

// establish context (e.g. "when creating a new sparkle")
// because of (e.g. "with an empty body")
// it behaves like (e.g. "it returns an error message")
type Context struct {
	because func() interface{}
	its     []ItFunc
	Result  interface{}
}
type ContextFunc func(*Context)
type BecauseFunc func() interface{}
type ItFunc func(*testing.T)

func Establish(t *testing.T, x ContextFunc) {
	context := Context{
		its: []ItFunc{},
	}
	x(&context)
	context.Result = context.because()
	for _, it := range context.its {
		it(t)
	}
}

func (c *Context) Because(b func() interface{}) {
	c.because = b
}

func (c *Context) It(it ItFunc) {
	c.its = append(c.its, it)
}
