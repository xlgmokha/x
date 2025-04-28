package log

import (
	"net/http"

	"github.com/xlgmokha/x/pkg/mapper"
)

func init() {
	mapper.Register[*http.Request, Fields](func(r *http.Request) Fields {
		return Fields{
			"host":        r.URL.Host,
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_host": r.RemoteAddr,
		}
	})
}
