package cookie

import (
	"net/http"

	"github.com/xlgmokha/x/pkg/x"
)

func New(name string, options ...x.Option[*http.Cookie]) *http.Cookie {
	options = x.Prepend[x.Option[*http.Cookie]](
		options,
		With(func(c *http.Cookie) {
			c.Name = name
		}),
	)
	return x.New[*http.Cookie](options...)
}
