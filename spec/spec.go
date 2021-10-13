package spec

// establish context (e.g. "when creating a new sparkle")
// because of (e.g. "with an empty body")
// it behaves like (e.g. "it returns an error message")
type Context struct {
	because BecauseFunc
	its     []ItFunc
	state   map[string]interface{}
}
type ContextFunc func(*Context)
type BecauseFunc func()
type ItFunc func()

func Establish(x ContextFunc) {
	context := Context{
		its:     []ItFunc{},
		because: func() {},
		state:   make(map[string]interface{}),
	}
	x(&context)
	context.because()
	for _, it := range context.its {
		it()
	}
}

func (c *Context) Because(b BecauseFunc) {
	c.because = b
}

func (c *Context) It(it ItFunc) {
	c.its = append(c.its, it)
}

func (c *Context) Set(key string, value interface{}) *Context {
	c.state[key] = value
	return c
}

func (c *Context) Get(key string) interface{} {
	return c.state[key]
}
