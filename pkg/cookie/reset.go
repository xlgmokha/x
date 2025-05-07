package cookie

import (
	"net/http"
	"time"

	"github.com/xlgmokha/x/pkg/x"
)

func Reset(name string, options ...x.Option[*http.Cookie]) *http.Cookie {
	return New(name, append(
		options,
		WithValue(""),
		WithExpiration(time.Unix(0, 0)),
	)...)
}
