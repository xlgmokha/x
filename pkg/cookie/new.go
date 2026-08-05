package cookie

import (
	"net/http"

	"github.com/xlgmokha/x/pkg/x"
)

func New(name string, options ...x.Option[*http.Cookie]) *http.Cookie {
	return x.New(x.Prepend(options, WithName(name))...)
}
