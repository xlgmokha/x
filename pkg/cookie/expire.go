package cookie

import (
	"net/http"

	"github.com/xlgmokha/x/pkg/x"
)

func Expire(w http.ResponseWriter, name string, options ...x.Option[*http.Cookie]) {
	http.SetCookie(w, Reset(name, options...))
}
