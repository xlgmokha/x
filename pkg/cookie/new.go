package cookie

import (
	"net/http"

	"github.com/xlgmokha/x/pkg/x"
)

func New(name string, options ...x.Option[*http.Cookie]) *http.Cookie {
	return x.New[*http.Cookie](x.Prepend[x.Option[*http.Cookie]](options, WithName(name))...)
}
